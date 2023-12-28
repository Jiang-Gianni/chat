package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/message"
	"github.com/Jiang-Gianni/chat/room"
	"github.com/Jiang-Gianni/chat/user"
	"github.com/Jiang-Gianni/chat/views"
	"github.com/Jiang-Gianni/chat/web"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	userClient *user.GRPCClient
	roomClient *room.GRPCClient

	err error
	db  *sql.DB
	// TODO
	lg *slog.Logger
)

func init() {
	userClient, err = user.NewGRPCClient(
		config.UserServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	roomClient, err = room.NewGRPCClient(
		config.RoomServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	lg = slog.Default()
	db = config.Sqlite
}

func main() {
	defer userClient.Conn.Close()
	webServer := web.Server{
		Addr:         config.WebServiceAddr,
		Log:          lg,
		ChatWs:       ChatWs,
		PostLogin:    PostLogin(),
		PostRegister: PostRegister(),
		RoomQuerier:  room.New(db),
		PostRoom:     PostRoom(),
	}
	fmt.Printf("webService litening on port %s\n", config.WebServiceAddr)
	log.Fatal(webServer.Run())
}

func ChatWs(ws *websocket.Conn) {
	for {
		_, b, err := ws.ReadMessage()
		if err != nil {
			lg.Error("web", "error", fmt.Errorf("ws.ReadMessage: %w", err))
			return
		}
		wsm := message.Message{}
		json.Unmarshal(b, &wsm)
		wsm.Username = "ITS ME"
		err = web.WriteChatMessage(wsm, ws)
		if err != nil {
			lg.Error("web", "error", fmt.Errorf("web.WriteChatMessage: %w", err))
			return
		}

		fmt.Printf("%T and %s", wsm, wsm.Message)
	}
}

func PostLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password := r.FormValue("username"), r.FormValue("password")
		req := &user.LoginRequest{
			Username: username,
			Password: password,
		}
		_, err := userClient.Login(r.Context(), req)
		if err != nil {
			views.WriteLoginRegisterError(w, "Invalid credentials")
			lg.Info("web", "error", fmt.Errorf("userClient.Login: %w", err), "username", username)
			return
		}
		w.Header().Add("HX-Redirect", config.ChatEndpoint)
	}
}

func PostRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password := r.FormValue("username"), r.FormValue("password")
		req := &user.RegisterRequest{
			Username: username,
			Password: password,
		}
		_, err := userClient.Register(r.Context(), req)
		if err != nil {
			lg.Info(fmt.Sprintf("userClient.Register: %s", err), "service", "web")
			if errors.Is(err, user.UsernameTakenError) {
				views.WriteLoginRegisterError(w, "Username already taken")
				return
			}
			views.WriteLoginRegisterError(w, "Internal Server Error")
			return
		}
		w.Header().Add("HX-Redirect", config.ChatEndpoint)
	}
}

func PostRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomName := r.FormValue("room-name")
		fmt.Println("POST ROOM CALLED")
		fmt.Println(roomName)
		req := &room.CreateRequest{
			RoomName: roomName,
		}
		resp, err := roomClient.Create(r.Context(), req)
		if err != nil {
			lg.Error(fmt.Sprintf("roomClient.Create: %s", err), "service", "web")
			views.WriteNewChatError(w, "post error")
			return
		}
		w.Header().Add("HX-Redirect", config.ChatRoomIDEndpoint(int(resp.RoomId)))
	}
}
