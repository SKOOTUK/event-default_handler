package nats

import (
	"fmt"
	"log"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
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
	conn       stan.Conn
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

// Connect connect to STAN and listen to notifyClose
func (c *Connection) Connect() error {
	var err error
	uri := os.Getenv("NATS_ADDRESS")
	nc, err := nats.Connect(uri)
	if err != nil {
		sentry.CaptureException(err)
		return fmt.Errorf("error in creating NATS connection with %s : %s", uri, err.Error())
	}
	c.conn, err = stan.Connect(
		os.Getenv("NATS_CLUSTER_NAME"),
		os.Getenv("NATS_CLIENT_PREFIX")+uuid.New().String(),
		stan.NatsConn(nc),
		stan.Pings(10, 5),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			sentry.CaptureException(err)
			log.Fatalf("Connection lost, reason: %v", reason)
			for reason != nil {
				c.Connect()
			}
		}))
	if err != nil {
		sentry.CaptureException(err)
		return fmt.Errorf("error in creating STAN connection with %s : %s", uri, err.Error())
	}
	return nil
}

//Reconnect reconnects the connection
func (c *Connection) Reconnect() error {
	return c.Connect()
}

// Close the connection
func (c *Connection) Close() error {
	return c.conn.Close()
}
