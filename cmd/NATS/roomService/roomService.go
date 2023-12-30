package main

import (
	"log"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/room"
	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(config.NATS_URL)
	if err != nil {
		log.Fatal(err)
	}
	n := room.NATSServer{
		NATS:    nc,
		Queries: *room.New(config.Sqlite),
	}
	config.PrintListening(config.RoomService, config.RoomServiceAddr)
	log.Fatal(n.Run())
}
