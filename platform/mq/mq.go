package mq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type connection struct {
	Conn *amqp.Connection
}

var c *connection

func (cn *connection) connect(user, password, host, port string) *connection {
	var err error

	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)

	cn.Conn, err = amqp.Dial(connStr)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Connected to MQ on %s:%s", host, port)
	}

	return cn
}

func (cn *connection) getChannel() (*amqp.Channel, error) {
	channel, err := cn.Conn.Channel()
	if err != nil {
		log.Println(err)
	}

	return channel, err
}
func (cn *connection) Publish(message string) error {

	channel, _ := cn.getChannel()
	defer channel.Close()

	err := channel.ExchangeDeclare(
		"autodialer", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		log.Println(err)
	}

	err = channel.Publish(
		"autodialer",  // exchange
		"msisdn.list", // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	return err

}

func (cn *connection) Consume() {
	channel, _ := cn.getChannel()
	defer channel.Close()

	err := channel.ExchangeDeclare(
		"autodialer", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		log.Println(err)
	}

	q, err := channel.QueueDeclare(
		"msisdn", // name
		false,    // durable
		false,    // delete when usused
		true,     // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Println(err)
	}

	err = channel.QueueBind(
		q.Name,        // queue name
		"msisdn.list", // routing key
		"autodialer",  // exchange
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}

	msgs, err := channel.Consume(
		q.Name,            // queue
		"msisdn_consumer", // consumer
		true,              // auto ack
		false,             // exclusive
		false,             // no local
		false,             // no wait
		nil,               // args
	)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

type Connection interface {
	Publish(string) error
	Consume()
}

func ConnectionProvider(user, password, host, port string) Connection {
	if c == nil {
		c = new(connection).connect(user, password, host, port)
	}

	return c
}

func GetConnection() Connection {
	return c
}
