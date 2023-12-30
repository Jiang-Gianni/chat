package main

import (
	"log"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/user"
	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(config.NATS_URL)
	if err != nil {
		log.Fatal(err)
	}
	n := user.NATSServer{
		NATS:    nc,
		Queries: *user.New(config.Sqlite()),
	}
	config.PrintListening(config.UserService, config.UserServiceAddr)
	log.Fatal(n.Run())
}
