package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/dfrr"
	"github.com/Jiang-Gianni/chat/room"
	"github.com/Jiang-Gianni/chat/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

const KB = 1024

var wsUpgrader = websocket.Upgrader{
	HandshakeTimeout: time.Second,
	ReadBufferSize:   KB,
	WriteBufferSize:  KB,
	CheckOrigin:      func(r *http.Request) bool { return true },
}

type Server struct {
	Addr string
	Log  *slog.Logger

	// Chat websocket handler
	ChatWs func(*websocket.Conn)

	// Post '/login' http handler
	PostLogin http.HandlerFunc

	// Post '/register' http handler
	PostRegister http.HandlerFunc

	// To get the room data
	RoomQuerier room.Querier

	// Post '/room' http handler
	PostRoom http.HandlerFunc
}

func (s *Server) Run() (err error) {
	defer dfrr.Wrap(&err, "s.Run()")
	r := chi.NewRouter()
	r.Get("/", s.index())

	r.Get(config.DiscardEndpoint, func(w http.ResponseWriter, r *http.Request) {})
	r.Get(config.ChatEndpoint, s.getChat())
	r.Get(config.ChatParamEndpoint, s.getChat())
	r.Get(config.ChatRedirectParamEndpoint, s.getChatRediect())
	r.Post(config.LoginEndpoint, s.PostLogin)
	r.Post(config.RegisterEndpoint, s.PostRegister)
	r.Post(config.RoomEndpoint, s.PostRoom)

	// TODO
	r.Get("/chat/ws", s.chatWebSocket())
	r.Get("/sse", s.sse())

	// Write timeout removed to support server side events
	srv := http.Server{
		Addr:        s.Addr,
		Handler:     r,
		ReadTimeout: 5 * time.Second,
		IdleTimeout: 5 * time.Second,
	}
	return srv.ListenAndServe()
}

func (s *Server) index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		views.WriteLoginPage(w)
	}
}

func (s *Server) getChatRediect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomIDString := chi.URLParam(r, "roomID")
		roomID, err := strconv.Atoi(roomIDString)
		if err != nil {
			http.Error(w, "Inexistent Chat Room", http.StatusNotFound)
			return
		}
		w.Header().Add("HX-Redirect", config.ChatRoomIDEndpoint(roomID))
	}
}

func (s *Server) sse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("SSE CONNECTED")
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "SSE not supported", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		flusher.Flush()
		for {
			// _, err := fmt.Fprintf(w, "data: %s\n\n", views.SSEMessage2())
			_, err := fmt.Fprintf(w, "data: %s", views.SSEMessage())
			if err != nil {
				s.Log.Error(err.Error())
				return
			}
			flusher.Flush()
			fmt.Printf("data: %s\n\n", views.SSEMessage2())
			time.Sleep(time.Second)
			fmt.Println("SLEPT")
		}
	}
}
