package room

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/dfrr"
	"github.com/nats-io/nats.go"
)

type NATSServer struct {
	Queries
	NATS *nats.Conn
}

type CreateEventNATS struct {
	RoomName string
	ReplyTo  string
}

type CreateReplyNATS struct {
	RoomID     int
	StatusCode int
}

// Interface Check
var _ Querier = (*NATSServer)(nil)

func (n *NATSServer) Run() (rerr error) {
	defer dfrr.Wrap(&rerr, "n.Run")
	ec, err := nats.NewEncodedConn(n.NATS, nats.JSON_ENCODER)
	if err != nil {
		return fmt.Errorf("nats.NewEncodedConn: %w", err)
	}

	ec.Subscribe(config.NATSRoomCreate, func(cen *CreateEventNATS) {
		roomID, err := create(context.Background(), &n.Queries, cen.RoomName)
		crn := CreateReplyNATS{StatusCode: StatusOK, RoomID: int(roomID)}
		if err != nil {
			crn.StatusCode = StatusInternalError
		}
		ec.Publish(cen.ReplyTo, crn)
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	return nil
}
