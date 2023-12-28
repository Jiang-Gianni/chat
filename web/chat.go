package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Jiang-Gianni/chat/dfrr"
	"github.com/Jiang-Gianni/chat/message"
	"github.com/Jiang-Gianni/chat/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type ChatMessage struct{}

func (s *Server) getChat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := s.RoomQuerier.GetRooms(r.Context())
		if err != nil {
			s.Log.Error(fmt.Sprintf("s.RoomQuerier.GetRooms: %s", err), "service", "web")
		}
		roomIDString := chi.URLParam(r, "roomID")
		roomID, _ := strconv.Atoi(roomIDString)
		views.WriteChatPage(w, rooms, roomID)
	}
}

func (s *Server) chatWebSocket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			s.Log.Error(fmt.Sprintf("wsUpgrader.Upgrade: %s", err), "service", "web")
			http.Error(w, "websocket handshake error", http.StatusInternalServerError)
		}
		defer func() {
			if err := ws.Close(); err != nil {
				s.Log.Error(fmt.Sprintf("ws.Close: %s", err), "service", "web")
			}
		}()
		s.ChatWs(ws)
	}
}

func WriteChatMessage(msg message.Message, ws *websocket.Conn) (err error) {
	defer dfrr.Wrap(&err, "WriteChatMessage")
	return ws.WriteMessage(websocket.TextMessage, []byte(views.NewMessage(msg)))
}
