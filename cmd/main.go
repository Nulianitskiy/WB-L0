package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

func main() {
	// Connect to a server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	// Simple Async Subscriber
	_, err = nc.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	if err != nil {
		log.Fatal(err)
	}
	// Close connection
	for {
	}
}
