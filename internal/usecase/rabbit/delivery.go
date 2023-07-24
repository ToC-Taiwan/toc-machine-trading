package rabbit

import "github.com/streadway/amqp"

func (c *Rabbit) establishDelivery(key string) <-chan amqp.Delivery {
	delivery, err := c.conn.BindAndConsume(key)
	if err != nil {
		c.logger.Fatal(err)
	}
	return delivery
}
