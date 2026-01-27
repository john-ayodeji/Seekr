package crawler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/john-ayodeji/Seekr/internal"
	"github.com/john-ayodeji/Seekr/internal/database"
	amqp "github.com/rabbitmq/amqp091-go"
)

type payload struct {
	JobID string `json:"job_id"`
	URL   string `json:"url"`
}

func ProcessHTML(conn *amqp.Connection) {
	msgs, err := internal.ConsumeFromQueue(conn, "html_fetcher.jobs")
	if err != nil {
		fmt.Println(err)
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}

	for msg := range msgs {
		var p payload
		if err := json.Unmarshal(msg.Body, &p); err != nil {
			fmt.Println(err)
			fmt.Println("invalid payload:", err)
			msg.Ack(false)
			continue
		}

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
				// unknown → retry but track via DLQ
				msg.Nack(false, false)
				continue
			}
		}

		data, err := internal.Cfg.Db.AddHtml(context.Background(), database.AddHtmlParams{
			ID:    uuid.New(),
			Jobid: p.JobID,
			Url:   p.URL,
			Html:  html.HTML,
		})
		if err != nil {
			fmt.Printf("fialed to save to db: %v", err)
			msg.Nack(false, true)
			continue
		}

		_ = internal.Cfg.Db.MarkFetched(context.Background(), p.URL)

		_ = msg.Ack(false)

		if err := internal.PublishToQueue(ch, internal.RabbitCfg.Exchange, "page.fetch.success", data); err != nil {
			fmt.Println(err)
			continue
		}
	}
}
