package nats

import (
	"fmt"

	"github.com/nats-io/stan.go"
)

//Publish publishes a request to the NATS queue
func (c *Connection) Publish(m *stan.Msg) error {
	select { //non blocking channel - if there is no error will go to default where we do nothing
	case err := <-c.err:
		if err != nil {
			c.Reconnect()
		}
	default:
	}

	if err := c.conn.Publish(c.routingKey, m.Data); err != nil {
		return fmt.Errorf("Error in Publishing: %s", err)
	}
	return nil
}
