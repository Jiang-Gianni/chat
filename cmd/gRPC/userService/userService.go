package main

import (
	"log"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/user"
)

func main() {
	g := user.GRPCServer{
		Queries: *user.New(config.Sqlite()),
	}
	config.PrintListening(config.UserService, config.UserServiceAddr)
	log.Fatal(g.Run(config.UserServiceAddr))
}
