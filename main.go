package main

import (
	"log"

	"github.com/SKOOTUK/event-default_handler/pkg/handler"
	"github.com/streadway/amqp"
)

func main() {
	e := handler.Init{"", "", "", func(msg amqp.Delivery) {}}
	log.Println(e)
}
