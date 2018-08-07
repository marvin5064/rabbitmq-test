package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/streadway/amqp"
)

const AutoDeletedQueueTTL = 60 * 60 * 12

type MessageQueueSetup struct {
	Exchanges []MqExchange `json:"exchanges"`
}

type MqExchange struct {
	Name   string       `json:"name"`
	Type   string       `json:"type"`
	Queues []queueSetup `json:"queues"`
	// able to have exchange inside exchange
	Exchanges []MqExchange `json:"exchanges"`
}

type queueSetup struct {
	Name        string   `json:"name"`
	AutoDeleted bool     `json:"auto_deleted"`
	Topics      []string `json:"topics"`
}

func main() {
	mqConn, err := setupRabbitmq("amqp", "guest", "guest", "127.0.0.1", 8080)
	if err != nil {
		fmt.Println("setupRabbitmq", err)
		return
	}

	err = setupFromFile("./config/rabbitmq.json", mqConn)
	if err != nil {
		fmt.Println("setupFromFile", err)
		return
	}
	return
}

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

func setupFromFile(filename string, conn *amqp.Connection) error {
	if filename == "" {
		return nil
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	var mqSetup MessageQueueSetup
	err = json.Unmarshal(file, &mqSetup)
	if err != nil {
		return err
	}
	return setupExchangesAndQueues(mqSetup, conn)
}

func setupExchangesAndQueues(mqSetup MessageQueueSetup, conn *amqp.Connection) error {
	for _, e := range mqSetup.Exchanges {
		err := setupExchange(e, conn)
		if err != nil {
			return err
		}
	} // exchanges
	return nil
}

func bindQueueTopics(name, exchange string, topics []string, conn *amqp.Connection) error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	for _, topic := range topics {
		err = channel.QueueBind(name, topic, exchange, false, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func setupExchange(e MqExchange, conn *amqp.Connection) error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	if e.Type != "" {
		args := make(amqp.Table)
		if e.Type == "x-delayed-message" {
			args["x-delayed-type"] = "direct"
		}
		fmt.Println("Creating Exchange", e.Name, e.Type)

		err = channel.ExchangeDeclare(e.Name, e.Type, true, false, false, false, args)
		if err != nil {
			return fmt.Errorf("exchangeDeclare (%s): %v", e.Name, err)
		}
	}
	for _, q := range e.Queues {
		fmt.Println("Creating Queue", q.Name)
		var (
			queue amqp.Queue
			err   error
		)
		if q.AutoDeleted {
			queue, err = channel.QueueDeclare(q.Name, false, true, false, false, amqp.Table{
				"x-expires": AutoDeletedQueueTTL,
			})
		} else {
			queue, err = channel.QueueDeclare(q.Name, true, false, false, false, nil)
		}
		if err != nil {
			return fmt.Errorf("queueDeclare auto deleted?(%v) (%s): %v", q.AutoDeleted, q.Name, err)
		}

		if len(q.Topics) > 0 {
			fmt.Println("Binding Topics", q.Topics)
			err = bindQueueTopics(queue.Name, e.Name, q.Topics, conn)
			if err != nil {
				return fmt.Errorf("queueBind (%v): %v", q.Topics, err)
			}
			continue
		}
		fmt.Println("Bind Queue", q.Name, "to", e.Name)
		err = channel.QueueBind(queue.Name, "", e.Name, false, nil)
		if err != nil {
			return fmt.Errorf("queueBind (%v -> %v): %v", queue.Name, e.Name, err)
		}
	} // queues

	for _, exchange := range e.Exchanges {
		err := setupExchange(exchange, conn)
		if err != nil {
			return err
		}

		err = channel.ExchangeBind(exchange.Name, "", e.Name, false, nil)
		if err != nil {
			return err
		}
		fmt.Println("Bind Exchange", exchange.Name, "to", e.Name)
	}
	return nil
}
