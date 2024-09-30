package models

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker interface {
	Publish(string) 
}

type BrokerStruct struct {
    ConnectionString string
    errLog *log.Logger
}

func NewBroker(connectionString string, infoLog *log.Logger, errLog *log.Logger) *BrokerStruct {
	return &BrokerStruct{
        ConnectionString: connectionString,
		errLog:  errLog,
	}
}

func (b BrokerStruct) Publish(message string) {
	conn, err := amqp.Dial(b.ConnectionString)
	if err != nil {
		b.errLog.Printf("Could not establish connection to rabbitmq")
	}
	defer conn.Close()
	ch, err := conn.Channel()
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		b.errLog.Printf("Failed to declare q")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		b.errLog.Printf("Failed to publish")
	}
	log.Printf(" [x] Sent %s\n", message)
}
