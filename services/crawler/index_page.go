package crawler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/john-ayodeji/Seekr/internal"
	amqp "github.com/rabbitmq/amqp091-go"
)

func IndexPage(conn *amqp.Connection) {
	msgs, err := internal.ConsumeFromQueue(conn, "index_page.jobs")
	if err != nil {
		fmt.Println(err)
		return
	}

	type ParsedHtml struct {
		ID    uuid.UUID
		Jobid string
		Url   string
	}

	for msg := range msgs {
		var m ParsedHtml

		if err := json.Unmarshal(msg.Body, &m); err != nil {
			fmt.Println(err)
			msg.Nack(false, true)
			continue
		}

		if err := internal.Cfg.Db.ReindexParsedHTMLByID(context.Background(), m.ID); err != nil {
			fmt.Println(err)
			msg.Nack(false, true)
			continue
		}

		if err := internal.Cfg.Db.MarkTokenizedAndIndexed(context.Background(), m.Url); err != nil {
			fmt.Println(err)
			msg.Nack(false, true)
			continue
		}

		msg.Ack(false)
	}
}
