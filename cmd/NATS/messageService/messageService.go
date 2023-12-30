package main

import (
	"log"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/message"
	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(config.NATS_URL)
	if err != nil {
		log.Fatal(err)
	}
	n := message.NATSServer{
		NATS:    nc,
		Queries: *message.New(config.Sqlite()),
	}
	config.PrintListening(config.MessageService, config.MessageServiceAddr)
	log.Fatal(n.Run())
}
