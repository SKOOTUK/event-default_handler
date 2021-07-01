// Package handler a generic event handler that can be passed a function to process events.
package stan

import (
	"log"
	"time"

	"github.com/SKOOTUK/event-default_handler/pkg/utils"
	sentry "github.com/getsentry/sentry-go"
	"github.com/nats-io/stan.go"
)

// Init pointer receiver for configuring HandleQueuedMessages
type Init struct {
	Name           string
	TopicName      string
	RoutingKey     string // STAN RoutingKey cannot contain wildcards
	MessageHandler stan.MsgHandler
}

// HandleQueuedMessages iterates and handles messages in queue based on config in init
func (e *Init) HandleQueuedMessages() {
	log.Printf("Starting up [%v] event listener", e.Name)

	// Setup Sentry
	utils.SetupSentry()
	defer sentry.Flush(2 * time.Second) // Flush buffered events before the program terminates

	conn := NewConnection(e.Name, e.TopicName, e.RoutingKey)
	if err := conn.Connect(); err != nil {
		sentry.CaptureException(err)
		log.Printf("reported to Sentry: %s", err)
		log.Fatalf("%e", err)
	}
	defer conn.Close()

	// Iterate NATS queue
	forever := make(chan bool)

	conn.HandleDeliveries(e.MessageHandler)

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")

	<-forever
}
