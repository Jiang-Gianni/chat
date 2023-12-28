package main

import (
	"fmt"
	"log"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/room"
)

func main() {
	g := room.GRPCServer{
		Queries: *room.New(config.Sqlite),
	}
	fmt.Printf("roomService litening on port %s\n", config.RoomServiceAddr)
	log.Fatal(g.RunGRPC(config.RoomServiceAddr))
}
