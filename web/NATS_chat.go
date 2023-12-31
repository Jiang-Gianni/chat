package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/dfrr"
	"github.com/Jiang-Gianni/chat/message"
	"github.com/Jiang-Gianni/chat/views"
	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
)

func (n *NATSServer) getChat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomIDString := chi.URLParam(r, "roomID")
		roomID, _ := strconv.Atoi(roomIDString)
		rooms, err := n.RoomQuerier.GetRooms(r.Context())
		if err != nil {
			n.Log.Error(fmt.Sprintf("n.GetRooms: %s", err), "service", "web")
		}
		messages, err := n.MessageQuerier.GetMessageByRoomID(r.Context(), int64(roomID))
		if roomID > 0 && err != nil {
			n.Log.Error(fmt.Sprintf("n.GetMessageByRoomID: %s", err), "service", "web")
		}
		username := ctxUsername(r)
		views.WriteChatPage(w, rooms, roomID, messages, username)
	}
}

func (n *NATSServer) getChatWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			n.Log.Error(fmt.Sprintf("wsUpgrader.Upgrade: %s", err), "service", "web")
			http.Error(w, "websocket handshake error", http.StatusInternalServerError)
		}
		defer func() {
			if err := ws.Close(); err != nil {
				n.Log.Error(fmt.Sprintf("ws.Close: %s", err), "service", "web")
			}
		}()

		roomIDString := chi.URLParam(r, "roomID")
		roomID, err := strconv.Atoi(roomIDString)
		if err != nil {
			n.Log.Error(fmt.Sprintf("strconv.Atoi: %s", err), "service", "web")
			http.Error(w, "room ID error", http.StatusInternalServerError)
		}

		username := ctxUsername(r)
		msg := message.Message{
			RoomID:   int64(roomID),
			Username: username,
		}
		roomSubject := config.NATSMessageStreamRoom(roomIDString)

		readLoop := func() (rerr error) {
			defer dfrr.Wrap(&rerr, "readLoop")
			for {
				_, b, err := ws.ReadMessage()
				if err != nil {
					return fmt.Errorf("ws.ReadMessage: %w", err)
				}
				// The browser client sends the data with a `message` field
				// Unmarshal into `msr` keeps the previously set RoomId and Username
				if err := json.Unmarshal(b, &msg); err != nil {
					return fmt.Errorf("json.Unmarshal: %w", err)
				}
				msg.SentAt = time.Now()
				err = n.EC.Publish(roomSubject, msg)
				if err != nil {
					return fmt.Errorf("n.EC.Publish: %w", err)
				}
			}
		}

		writeLoop := func() (rerr error) {
			defer dfrr.Wrap(&rerr, "writeLoop")
			quitError := make(chan error)
			defer func() {
				// Set to nil so that the select with default doesn't block or panic
				quitError = nil
			}()
			sub, err := n.EC.Subscribe(
				roomSubject,
				func(msg *message.Message) {
					if err := WriteChatMessage(*msg, ws, username); err != nil {
						select {
						case quitError <- fmt.Errorf("WriteChatMessage: %w", err):
						default:
						}
					}
				},
			)
			if err != nil {
				return fmt.Errorf("n.EC.Subscribe: %w", err)
			}
			defer func() {
				unsubErr := sub.Unsubscribe()
				if rerr == nil {
					rerr = unsubErr
				}
			}()
			return <-quitError
		}

		errg := errgroup.Group{}
		errg.Go(writeLoop)
		errg.Go(readLoop)

		n.Log.Info(errg.Wait().Error(), "service", "web")
	}
}
