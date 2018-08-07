package wildcard

import (
	"fmt"

	"github.com/streadway/amqp"
)

func setupRabbitmq(protocol, username, password, host string, port int) (*amqp.Connection, error) {
	url := fmt.Sprintf("%s://%s:%s@%s:%d",
		protocol,
		username,
		password,
		host,
		port,
	)
	return amqp.Dial(url)
}
