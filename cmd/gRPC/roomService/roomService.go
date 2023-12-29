package main

import (
	"log"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/room"
)

func main() {
	g := room.GRPCServer{
		Queries: *room.New(config.Sqlite),
	}
	config.PrintListening("roomService", config.RoomServiceAddr)
	log.Fatal(g.RunGRPC(config.RoomServiceAddr))
}
