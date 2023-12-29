package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/dfrr"
	"github.com/Jiang-Gianni/chat/message"
	"github.com/Jiang-Gianni/chat/views"
	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func (g *GRPCServer) getChat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomIDString := chi.URLParam(r, "roomID")
		roomID, _ := strconv.Atoi(roomIDString)
		rooms, err := g.RoomQuerier.GetRooms(r.Context())
		if err != nil {
			g.Log.Error(fmt.Sprintf("cs.GetRooms: %s", err), "service", "web")
		}
		messages, err := g.MessageQuerier.GetMessageByRoomID(r.Context(), int64(roomID))
		if roomID > 0 && err != nil {
			g.Log.Error(fmt.Sprintf("cs.GetMessageByRoomID: %s", err), "service", "web")
		}
		_, username := userIDUserName(r)
		views.WriteChatPage(w, rooms, roomID, messages, username)
	}
}

func (g *GRPCServer) getChatWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			g.Log.Error(fmt.Sprintf("wsUpgrader.Upgrade: %s", err), "service", "web")
			http.Error(w, "websocket handshake error", http.StatusInternalServerError)
		}
		defer func() {
			if err := ws.Close(); err != nil {
				g.Log.Error(fmt.Sprintf("ws.Close: %s", err), "service", "web")
			}
		}()

		roomIDString := chi.URLParam(r, "roomID")
		roomID, err := strconv.Atoi(roomIDString)
		if err != nil {
			g.Log.Error(fmt.Sprintf("strconv.Atoi: %s", err), "service", "web")
			http.Error(w, "room ID error", http.StatusInternalServerError)
		}

		// Message Service Client
		msgClient, err := message.NewGRPCClient(
			config.MessageServiceAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatal(err)
		}
		defer msgClient.Conn.Close()

		// Service stream
		ctx := context.Background()
		ctx = metadata.AppendToOutgoingContext(ctx, config.RoomIDStringKey, roomIDString)
		stream, err := msgClient.StreamMessage(ctx)
		if err != nil {
			g.Log.Error(fmt.Sprintf("msgClient.StreamMessage: %s", err), "service", "web")
		}
		defer func() {
			if err := stream.CloseSend(); err != nil {
				g.Log.Info(fmt.Sprintf("stream.CloseSend: %s", err), "service", "web")
			}
		}()

		_, username := userIDUserName(r)
		msr := &message.StreamRequest{
			RoomId:   int32(roomID),
			Username: username,
		}

		readLoop := func() (rerr error) {
			defer dfrr.Wrap(&rerr, "readLoop")
			for {
				_, b, err := ws.ReadMessage()
				if err != nil {
					return fmt.Errorf("ws.ReadMessage: %w", err)
				}
				// The browser client sends the data with a `message` field
				// Unmarshal into `msr` keeps the previously set `RoomId` and Username
				if err := json.Unmarshal(b, &msr); err != nil {
					return fmt.Errorf("json.Unmarshal: %w", err)
				}
				err = stream.Send(msr)
				if err != nil {
					return fmt.Errorf("stream.Send: %w", err)
				}
			}
		}

		var resp *message.StreamResponse
		writeLoop := func() (rerr error) {
			defer dfrr.Wrap(&rerr, "writeLoop")
			for {
				resp, err = stream.Recv()
				if err != nil {
					return fmt.Errorf("stream.Recv: %w", err)
				}
				if err := WriteChatMessage(message.Message{
					Username: resp.Username,
					Message:  resp.Message,
					SentAt:   resp.SentAt.AsTime(),
				}, ws, username); err != nil {
					return fmt.Errorf("WriteChatMessage: %w", err)
				}
			}
		}

		errg := errgroup.Group{}
		errg.Go(writeLoop)
		errg.Go(readLoop)

		g.Log.Info(errg.Wait().Error(), "service", "web")
	}
}
