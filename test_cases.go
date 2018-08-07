package main

import (
	"github.com/marvin5064/rabbitmq-test/tests/wildcard"
	"github.com/streadway/amqp"
)

func RunTestCases(conn *amqp.Connection) {
	wildcard.QueueInspectTest(conn)
}
