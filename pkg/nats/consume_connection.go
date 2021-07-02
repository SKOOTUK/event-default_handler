package nats

import (
	"log"

	"github.com/nats-io/nats.go"
)

//HandleDeliveries handles the consumed deliveries from the queues
func (c *Connection) HandleDeliveries(fn nats.MsgHandler) {
	_, err := c.conn.QueueSubscribe(
		c.routingKey, // subject
		c.exchange,   // queue group
		fn,           // message handler function
	)
	if err != nil {
		// TODO other error handling
		// c.Reconnect()
		log.Fatal(err)
	}
}
