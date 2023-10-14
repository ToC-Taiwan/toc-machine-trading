// Package rabbitmq package rabbitmq
package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	_defaultConnAttempts = 10
	_defaultConnWaitTime = 5
)

// Connection -.
type Connection struct {
	exchange string
	url      string
	waitTime time.Duration
	attempts int

	conn   *amqp.Connection
	logger Logger
}

// New -.
func New(exchange, url string, opts ...Option) (*Connection, error) {
	conn := &Connection{
		exchange: exchange,
		url:      url,
		waitTime: time.Duration(_defaultConnWaitTime) * time.Second,
		attempts: _defaultConnAttempts,
	}

	for _, opt := range opts {
		opt(conn)
	}

	if err := conn.connect(); err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

// Connect -.
func (c *Connection) connect() error {
	var err error
	for c.attempts > 0 {
		if c.conn, err = amqp.Dial(c.url); err == nil {
			break
		}

		c.attempts--
		if err != nil && c.attempts == 0 {
			return err
		}
		c.Infof("RabbitMQ is trying to connect, attempts left: %d\n", c.attempts)
		time.Sleep(c.waitTime)
	}
	return nil
}

// BindAndConsume -.
func (c *Connection) BindAndConsume(key string) (<-chan amqp.Delivery, error) {
	channel, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}
	err = channel.ExchangeDeclare(
		c.exchange,
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
		c.exchange,
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
	channel, err := c.conn.Channel()
	if err != nil {
		return err
	}
	err = channel.ExchangeDeclare(
		c.exchange,
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
	err = channel.PublishWithContext(
		context.Background(),
		c.exchange,
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

func (c *Connection) Infof(format string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Infof(format, args...)
	} else {
		fmt.Printf(format, args...)
	}
}

func (c *Connection) Errorf(format string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Errorf(format, args...)
	} else {
		fmt.Printf(format, args...)
	}
}
