package message

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

// Interface Check
var _ Querier = (*NATSServer)(nil)

func (n *NATSServer) Run() (rerr error) {
	defer dfrr.Wrap(&rerr, "n.Run")
	ec, err := nats.NewEncodedConn(n.NATS, nats.JSON_ENCODER)
	if err != nil {
		return fmt.Errorf("nats.NewEncodedConn: %w", err)
	}

	ec.Subscribe(config.NATSMessageStreamRoom("*"), func(m *Message) {
		_, _ = n.InsertMessage(context.Background(), InsertMessageParams{
			RoomID:   m.RoomID,
			Username: m.Username,
			Message:  m.Message,
			SentAt:   m.SentAt,
		})
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	return nil
}
