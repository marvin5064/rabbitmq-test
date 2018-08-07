package lib

import (
	"time"

	"github.com/streadway/amqp"
)

func Publish(conn *amqp.Connection, exchange, topic string, payload []byte) error {
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now().UTC(),
		ContentType:  "text/plain",
		Body:         payload,
	}

	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return channel.Publish(exchange, topic, false, false, msg)
}
