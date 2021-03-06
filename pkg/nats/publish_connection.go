package nats

import (
	"github.com/nats-io/nats.go"
)

//Publish publishes a request to the NATS queue
func (c *Connection) Publish(m *nats.Msg) error {
	select { //non blocking channel - if there is no error will go to default where we do nothing
	case err := <-c.err:
		if err != nil {
			c.Reconnect()
		}
	default:
	}

	return c.conn.Publish(c.routingKey, m.Data)
}
