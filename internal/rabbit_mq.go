package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitConfig struct {
	URL        string
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Exchange   string
}

var RabbitCfg *RabbitConfig

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func ConnectRabbitMQ(connString string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(connString)
	if err != nil {
		fmt.Printf("Rabbitmq connection failed: %v\n", err)
		return nil, nil, err
	}
	fmt.Println(conn.Config.Properties)
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Rabbitmq channel creation failed: %v\n", err)
		return nil, nil, err
	}

	return conn, ch, nil
}

func (r *RabbitConfig) CreateExchange() {
	if err := r.Channel.ExchangeDeclare(r.Exchange, "direct", true, false, false, false, nil); err != nil {
		FailOnError(err, "Failed to create exchange")
	}

	fmt.Println("Exchange created successfully")
}

func (r *RabbitConfig) DeclareAndBindQueue(queueName, key, exchange string, durable bool) (amqp.Queue, error) {
	q, err := r.Channel.QueueDeclare(queueName, durable, false, false, false, nil)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("failed to create queue")
	}
	fmt.Printf("Queue created successfully\n%v - %v - %v\n", q.Name, q.Messages, q.Consumers)

	if err := r.Channel.QueueBind(queueName, key, exchange, false, nil); err != nil {
		return amqp.Queue{}, fmt.Errorf("failed to bind queue to exchange")
	}
	fmt.Println("Queue successfully bind to exchange")

	return q, nil
}

func PublishToQueue[T any](ch *amqp.Channel, exchange string, key string, body T) error {
	bytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err1 := ch.PublishWithContext(ctx, exchange, key, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		})
	if err1 != nil {
		return err
	}
	return nil
}

func ConsumeFromQueue(conn *amqp.Connection, queue string) (<-chan amqp.Delivery, error) {
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("failed to create channel")
		return nil, err
	}
	msgs, err := ch.Consume(queue, "", false, false, false, false, nil)
	FailOnError(err, "Failed to register a consumer")

	return msgs, nil
}
