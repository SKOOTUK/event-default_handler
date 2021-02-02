package rabbit

import (
	"github.com/SKOOTUK/event-default_handler/pkg/utils"
	"github.com/streadway/amqp"
)

//Consume consumes the messages from the queues and passes it as map of chan of amqp.Delivery
func (c *Connection) Consume() (map[string]<-chan amqp.Delivery, error) {
	m := make(map[string]<-chan amqp.Delivery)
	for _, q := range c.queues {
		deliveries, err := c.channel.Consume(q, "", false, false, false, false, nil)
		if err != nil {
			return nil, err
		}
		m[q] = deliveries
	}
	return m, nil
}

//HandleConsumedDeliveries handles the consumed deliveries from the queues. Should be called only for a consumer connection
func (c *Connection) HandleConsumedDeliveries(delivery <-chan amqp.Delivery, fn func(<-chan amqp.Delivery)) {
	for {
		go fn(delivery)
		if err := <-c.err; err != nil {
			c.Reconnect()
			_, err := c.Consume()
			if err != nil {
				utils.FailOnError(err, "Unable to consume after reconnect")
			}
		}
	}
}
