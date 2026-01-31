package crawler

import (
    "context"
    "encoding/json"
    "errors"

    "github.com/google/uuid"
    "github.com/john-ayodeji/Seekr/internal"
    "github.com/john-ayodeji/Seekr/internal/database"
    amqp "github.com/rabbitmq/amqp091-go"
    "github.com/lib/pq"
)

type payload struct {
	URL string `json:"url"`
}

func ProcessHTML(conn *amqp.Connection) {
 msgs, err := internal.ConsumeFromQueue(conn, "html_fetcher.jobs")
 if err != nil {
     // fail fast on consumer setup issues
     return
 }

 ch, err := conn.Channel()
 if err != nil {
     // cannot open publishing channel
     return
 }

	for msg := range msgs {
		var p payload
  if err := json.Unmarshal(msg.Body, &p); err != nil {
            msg.Ack(false)
            continue
        }
		jobID := uuid.New().String()

		// Ensure a job row exists for this URL as soon as we pick it from the queue using existing CreateJob.
		// Ignore duplicate errors due to unique(url) constraint.
		_ = internal.Cfg.Db.CreateJob(context.Background(), database.CreateJobParams{
			ID:    uuid.New(),
			Jobid: jobID,
			Url:   p.URL,
		})

        html, err := GetHTML(p.URL)
        if err != nil {
            switch {
            case errors.Is(err, ErrForbidden),
                errors.Is(err, ErrNotFound):
                // permanent failure
                msg.Ack(false)
                continue

            case errors.Is(err, ErrRateLimited),
                errors.Is(err, ErrServerError):
                // temporary failure → retry
                msg.Nack(false, true)
                continue

            default:
                // unknown → drop to avoid stuck messages
                msg.Ack(false)
                continue
            }
        }

        data, err := internal.Cfg.Db.AddHtml(context.Background(), database.AddHtmlParams{
            ID:    uuid.New(),
            Jobid: jobID,
            Url:   p.URL,
            Html:  html.HTML,
        })
        if err != nil {
            // If the URL already exists (unique violation), treat as success and continue
            if pqErr, ok := err.(*pq.Error); ok && string(pqErr.Code) == "23505" {
                // Construct a compatible payload for downstream using the freshly fetched HTML
                data = database.RawHtml{
                    ID:    uuid.New(),
                    Jobid: jobID,
                    Url:   p.URL,
                    Html:  html.HTML,
                }
            } else {
                // For other DB errors: if already redelivered once, ack to drop; else requeue
                if msg.Redelivered {
                    msg.Ack(false)
                } else {
                    msg.Nack(false, true)
                }
                continue
            }
        }

        _ = internal.Cfg.Db.MarkFetched(context.Background(), p.URL)

        _ = msg.Ack(false)

        if err := internal.PublishToQueue(ch, internal.RabbitCfg.Exchange, "page.fetch.success", data); err != nil {
            continue
        }
    }
}
