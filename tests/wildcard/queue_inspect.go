package wildcard

import (
	"fmt"

	"github.com/marvin5064/rabbitmq-test/lib"

	"github.com/streadway/amqp"
)

func QueueInspectTest(conn *amqp.Connection) {
	// sent out random message to the queue
	lib.Publish(conn, "sportsbook_test", "test", []byte("payload"))
	lib.Publish(conn, "sportsbook_test", "test", []byte("payload"))
	lib.Publish(conn, "sportsbook_test", "test", []byte("payload"))

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Channel", err)
		return
	}
	for _, qName := range []string{"sportsbook_test_queue_1",
		"sportsbook_test_queue_2",
		"sportsbook_test_queue_3"} {

		q, err := ch.QueueInspect(qName)
		fmt.Println(qName, q, err)
	}

	qName := "sportsbook_test_queue"
	q, err := ch.QueueInspect(qName)
	fmt.Println(qName, q, err)

}
