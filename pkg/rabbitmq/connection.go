// Package rabbitmq package rabbitmq
package rabbitmq

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

// Connection -.
type Connection struct {
	Exchange   string
	URL        string
	WaitTime   time.Duration
	Attempts   int
	Connection *amqp.Connection
}

// NewConnection -.
func NewConnection(exchange string, url string, waitTime int64, attempts int) *Connection {
	conn := &Connection{
		Exchange: exchange,
		URL:      url,
		WaitTime: time.Duration(waitTime) * time.Second,
		Attempts: attempts,
	}
	return conn
}

// AttemptConnect -.
func (c *Connection) AttemptConnect() error {
	// TODO: attempts not working
	var err error
	for i := c.Attempts; i > 0; i-- {
		if err = c.connect(); err == nil {
			break
		}

		fmt.Printf("RabbitMQ is trying to connect, attempts left: %d\n", i)
		time.Sleep(c.WaitTime)
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) connect() error {
	var err error
	c.Connection, err = amqp.Dial(c.URL)
	if err != nil {
		return err
	}
	return nil
}

// BindAndConsume -.
func (c *Connection) BindAndConsume(key string) (<-chan amqp.Delivery, error) {
	channel, err := c.Connection.Channel()
	if err != nil {
		return nil, err
	}
	err = channel.ExchangeDeclare(
		c.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	queue, err := channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	err = channel.QueueBind(
		queue.Name,
		key,
		c.Exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	delivery, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return delivery, nil
}

func (c *Connection) Publish(key string, message []byte) error {
	channel, err := c.Connection.Channel()
	if err != nil {
		return err
	}
	err = channel.ExchangeDeclare(
		c.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = channel.Publish(
		c.Exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) Close() error {
	return c.Connection.Close()
}
