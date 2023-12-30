package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/room"
	"github.com/Jiang-Gianni/chat/views"
	"github.com/nats-io/nuid"
)

func (n *NATSServer) postRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomName := r.FormValue("room-name")
		if roomName == "" {
			views.WriteNewChatError(w, InvalidRoomName)
			return
		}
		cen := room.CreateEventNATS{
			RoomName: roomName,
			ReplyTo:  nuid.New().Next(),
		}

		err := n.EC.Publish(config.NATSRoomCreate, cen)
		if err != nil {
			n.Log.Error(fmt.Sprintf("ec.Publish: %s", err), "service", "web", "post", "register")
			views.WriteLoginRegisterError(w, InternalServerError)
			return
		}
		reply := &room.CreateReplyNATS{StatusCode: room.StatusOK}
		ctx, cancel := context.WithTimeout(r.Context(), time.Second)
		defer cancel()
		sub, err := n.EC.Subscribe(cen.ReplyTo, func(lrn *room.CreateReplyNATS) {
			reply = lrn
			cancel()
		})
		if err != nil {
			n.Log.Error(fmt.Sprintf("ec.Subscribe: %s", err), "service", "web", "post", "room")
			views.WriteLoginRegisterError(w, InternalServerError)
			return
		}
		defer func() {
			if err := sub.Unsubscribe(); err != nil {
				n.Log.Info(
					fmt.Sprintf("sub.Unsubscribe: %s", err),
					"service",
					"web",
					"post",
					"room",
				)
			}
		}()
		<-ctx.Done()
		n.Log.Info(
			"login",
			"service",
			"web",
			"room",
			roomName,
			"code",
			reply.StatusCode,
			"post",
			"register",
		)
		switch reply.StatusCode {
		case room.StatusInternalError:
			views.WriteLoginRegisterError(w, InternalServerError)
			return
		default:
			w.Header().Add("HX-Redirect", config.ChatRoomIDEndpoint(reply.RoomID))
		}
	}
}
