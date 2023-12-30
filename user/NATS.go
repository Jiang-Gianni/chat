package user

import (
	"context"
	"errors"
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

type LoginEventNATS struct {
	Username string
	Password string
	ReplyTo  string
}

type LoginReplyNATS struct {
	StatusCode int
}

type RegisterEventNATS struct {
	Username string
	Password string
	ReplyTo  string
}

type RegisterReplyNATS struct {
	StatusCode int
}

func (n *NATSServer) Run() (rerr error) {
	defer dfrr.Wrap(&rerr, "n.Run")
	ec, err := nats.NewEncodedConn(n.NATS, nats.JSON_ENCODER)
	if err != nil {
		return fmt.Errorf("nats.NewEncodedConn: %w", err)
	}

	ec.Subscribe(config.NATSUserLogin, func(len *LoginEventNATS) {
		lrn := LoginReplyNATS{StatusCode: StatusOK}
		err := login(context.Background(), &n.Queries, len.Username, len.Password)
		if errors.Is(err, InvalidCredentialsError) {
			lrn.StatusCode = StatusInvalidCredentials
		} else if err != nil {
			lrn.StatusCode = StatusInternalError
		}
		ec.Publish(len.ReplyTo, lrn)
	})

	ec.Subscribe(config.NATSUserRegister, func(ren *RegisterEventNATS) {
		rrn := RegisterReplyNATS{StatusCode: StatusOK}
		err := register(context.Background(), &n.Queries, ren.Username, ren.Password)
		if errors.Is(err, UsernameTakenError) {
			rrn.StatusCode = StatusUsernameTaken
		} else if err != nil {
			rrn.StatusCode = StatusInternalError
		}
		ec.Publish(ren.ReplyTo, rrn)
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	return nil
}
