package crawler

import (
	"encoding/json"
	"fmt"

	"github.com/john-ayodeji/Seekr/internal"
	amqp "github.com/rabbitmq/amqp091-go"
)

type payload struct {
	URL string `json:"url"`
}

func ProcessHTML(conn *amqp.Connection) {
	msgs, err := internal.ConsumeFromQueue(conn, "html_fetcher.jobs")
	if err != nil {
		fmt.Println(err)
		return
	}

	for msg := range msgs {
		var p payload
		if err := json.Unmarshal(msg.Body, &p); err != nil {
			fmt.Println(err)
			return
		}

		html, err := GetHTML(p.URL)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
