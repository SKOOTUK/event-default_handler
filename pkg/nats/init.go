package nats

import (
	"fmt"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/nats-io/nats.go"
)

//MessageBody is the struct for the body passed in the NATS message. The type will be set on the Request header
type MessageBody struct {
	Data []byte
	Type string
}

//Message is the NATS request to publish
type Message struct {
	Queue string
	Body  MessageBody
}

//Connection is the connection created
type Connection struct {
	name       string
	conn       *nats.Conn
	exchange   string
	routingKey string
	err        chan error
}

var connectionPool = make(map[string]*Connection)

//NewConnection returns the new connection object
func NewConnection(name, exchange string, routingKey string) *Connection {
	if c, ok := connectionPool[name]; ok {
		return c
	}
	c := &Connection{
		exchange:   exchange,
		routingKey: routingKey,
		err:        make(chan error),
	}
	connectionPool[name] = c
	return c
}

//GetConnection returns the instantiated connection
func GetConnection(name string) *Connection {
	return connectionPool[name]
}

// Connect connect to nats
func (c *Connection) Connect() error {
	var err error
	uri := os.Getenv("NATS_ADDRESS")
	c.conn, err = nats.Connect(uri)
	if err != nil {
		sentry.CaptureException(err)
		return fmt.Errorf("error in creating NATS connection with %s : %s", uri, err.Error())
	}
	return nil
}

//Reconnect reconnects the connection
func (c *Connection) Reconnect() error {
	return c.Connect()
}

// Close the connection
func (c *Connection) Close() error {
	c.conn.Close()
	return nil
}
