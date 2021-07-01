package nats

import (
	"os"

	sentry "github.com/getsentry/sentry-go"
	"github.com/nats-io/stan.go"
)

//HandleDeliveries handles the consumed deliveries from the queues.
func (c *Connection) HandleDeliveries(fn stan.MsgHandler) {
	_, err := c.conn.QueueSubscribe(
		c.routingKey, // subject
		c.exchange,   // queue group
		fn,           // message handler function
		stan.DurableName(os.Getenv("NATS_DURABLE_NAME")),
		stan.StartWithLastReceived(),
	)
	if err != nil {
		sentry.CaptureException(err)
		c.Reconnect()
	}
}
