package main

import (
	"fmt"
	"log"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/user"
)

func main() {
	g := user.GRPCServer{
		Queries: *user.New(config.Sqlite),
	}
	fmt.Printf("userService litening on port %s\n", config.UserServiceAddr)
	log.Fatal(g.RunGRPC(config.UserServiceAddr))
}
