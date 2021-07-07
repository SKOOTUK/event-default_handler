package stan

import (
	"os"

	sentry "github.com/getsentry/sentry-go"
	"github.com/nats-io/stan.go"
)

//HandleDeliveries handles the consumed deliveries from the queues.
func (c *Connection) HandleDeliveries(fn stan.MsgHandler, mack bool) {
	var err error = nil
	if mack {
		_, err = c.conn.QueueSubscribe(
			c.routingKey, // subject
			c.exchange,   // queue group
			fn,           // message handler function
			stan.DurableName(os.Getenv("NATS_DURABLE_NAME")),
			stan.StartWithLastReceived(),
			stan.SetManualAckMode(), // sets manual ack mode
		)
	} else {
		_, err = c.conn.QueueSubscribe(
			c.routingKey, // subject
			c.exchange,   // queue group
			fn,           // message handler function
			stan.DurableName(os.Getenv("NATS_DURABLE_NAME")),
			stan.StartWithLastReceived(),
		)
	}

	if err != nil {
		sentry.CaptureException(err)
		c.Reconnect()
	}
}
