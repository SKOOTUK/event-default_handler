package nats

import (
	"log"
	"os"

	"github.com/nats-io/stan.go"
)

//HandleDeliveries handles the consumed deliveries from the queues. Should be called only for a consumer connection
func (c *Connection) HandleDeliveries(fn stan.MsgHandler) {
	_, err := c.conn.QueueSubscribe(
		c.routingKey, // subject
		c.exchange,   // queue group
		fn,           // message handler function
		stan.DurableName(os.Getenv("NATS_DURABLE_NAME")),
		stan.StartWithLastReceived(),
	)
	if err != nil {
		// TODO other error handling
		// c.Reconnect()
		log.Fatal(err)
	}
}
