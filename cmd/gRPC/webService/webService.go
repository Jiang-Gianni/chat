package main

import (
	"log"
	"log/slog"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/message"
	"github.com/Jiang-Gianni/chat/room"
	"github.com/Jiang-Gianni/chat/user"
	"github.com/Jiang-Gianni/chat/web"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	userClient, err := user.NewGRPCClient(
		config.UserServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer userClient.Conn.Close()

	roomClient, err := room.NewGRPCClient(
		config.RoomServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer roomClient.Conn.Close()

	webServer := web.GRPCServer{
		Addr: config.WebServiceAddr,
		// TODO Logger
		Log:            slog.Default(),
		UserClient:     userClient,
		RoomClient:     roomClient,
		MessageQuerier: message.New(config.Sqlite),
		RoomQuerier:    room.New(config.Sqlite),
	}

	config.PrintListening(config.WebService, config.WebServiceAddr)
	log.Fatal(webServer.Run())
}
