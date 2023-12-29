package main

import (
	"log"
	"sync"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/message"
)

func main() {
	g := message.GRPCServer{
		Queries:     *message.New(config.Sqlite),
		RoomStreams: map[string]map[*message.Message_StreamMessageServer]struct{}{},
		Lock:        &sync.Mutex{},
	}
	config.PrintListening("messageService", config.MessageServiceAddr)
	log.Fatal(g.RunGRPC(config.MessageServiceAddr))
}
