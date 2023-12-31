package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/dfrr"
	"github.com/Jiang-Gianni/chat/message"
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

func index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok, _ := isAuth(r)
		if ok {
			http.Redirect(w, r, config.ChatEndpoint, http.StatusSeeOther)
			return
		}
		views.WriteLoginPage(w)
	}
}

func postLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clearTokenJWT(w)
		w.Header().Add("HX-Redirect", config.IndexEndpoint)
	}
}

func getChatRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomIDString := chi.URLParam(r, "roomID")
		roomID, err := strconv.Atoi(roomIDString)
		if err != nil {
			w.Header().Add("HX-Redirect", config.ChatEndpoint)
			return
		}
		w.Header().Add("HX-Redirect", config.ChatRoomIDEndpoint(roomID))
	}
}

func WriteChatMessage(msg message.Message, ws *websocket.Conn, currentUser string) (err error) {
	defer dfrr.Wrap(&err, "WriteChatMessage")
	return ws.WriteMessage(websocket.TextMessage, []byte(views.NewMessage(msg, currentUser)))
}
