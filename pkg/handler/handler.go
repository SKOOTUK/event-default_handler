// Package handler a generic event handler that can be passed a function to process events.
package handler

import (
	"log"
	"time"

	"github.com/SKOOTUK/event-default_handler/pkg/rabbit"
	"github.com/SKOOTUK/event-default_handler/pkg/utils"
	sentry "github.com/getsentry/sentry-go"
	"github.com/streadway/amqp"
)

// Init pointer receiver for configuring HandleQueuedMessages
type Init struct {
	Name, TopicName, RoutingKey string
	MessageHandler              func(amqp.Delivery)
}

// HandleQueuedMessages iterates and handles messages in queue based on config in init
func (e *Init) HandleQueuedMessages() {
	log.Printf("Starting up [%v] event listener", e.Name)

	// Setup Sentry
	utils.SetupSentry()
	defer sentry.Flush(2 * time.Second) // Flush buffered events before the program terminates

	// Iterate Rabbit queue
	forever := make(chan bool)

	conn := rabbit.NewConnection(e.Name, e.TopicName, []string{e.Name + "-1"}, e.RoutingKey)
	if err := conn.Connect(); err != nil {
		utils.FailOnError(err, "Failed to connect to rabbit")
	}
	if err := conn.BindQueue(); err != nil {
		utils.FailOnError(err, "Failed to bind to rabbit queue")
	}
	deliveries, err := conn.Consume()
	if err != nil {
		utils.FailOnError(err, "Failed to consume from rabbit queue")
	}

	log.Printf("About to attempt to process queue items")
	for _, d := range deliveries {
		go conn.HandleConsumedDeliveries(
			d,
			func(msgs <-chan amqp.Delivery) {
				for d := range msgs {
					e.MessageHandler(d)
				}
			},
		)
	}
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
