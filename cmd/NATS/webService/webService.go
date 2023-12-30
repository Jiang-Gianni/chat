package main

import (
	"log"
	"log/slog"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/message"
	"github.com/Jiang-Gianni/chat/room"
	"github.com/Jiang-Gianni/chat/web"
	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(config.NATS_URL)
	if err != nil {
		log.Fatal(err)
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	n := web.NATSServer{
		Addr: config.WebServiceAddr,
		NATS: nc,
		EC:   ec,
		// TODO Logger
		Log:            slog.Default(),
		MessageQuerier: message.New(config.Sqlite),
		RoomQuerier:    room.New(config.Sqlite),
	}
	config.PrintListening(config.WebService, config.WebServiceAddr)
	log.Fatal(n.Run())
}
