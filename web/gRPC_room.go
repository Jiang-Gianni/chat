package web

import (
	"fmt"
	"net/http"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/room"
	"github.com/Jiang-Gianni/chat/views"
)

func (g *GRPCServer) postRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomName := r.FormValue("room-name")
		req := &room.CreateRequest{
			RoomName: roomName,
		}
		resp, err := g.RoomClient.Create(r.Context(), req)
		if err != nil {
			g.Log.Error(fmt.Sprintf("roomClient.Create: %s", err), "service", "web")
			views.WriteNewChatError(w, "post error")
			return
		}
		w.Header().Add("HX-Redirect", config.ChatRoomIDEndpoint(int(resp.RoomId)))
	}
}
