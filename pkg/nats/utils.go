package nats

import "fmt"

// PrintConnection reconnects the connection
func (c *Connection) PrintConnection() {
	fmt.Println("Name        :", c.name)
	fmt.Println("Conn        :", c.conn)
	fmt.Println("Exchange    :", c.exchange)
	fmt.Println("Routing Key :", c.routingKey)
	fmt.Println("Err         :", c.err)
}
