package web

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/dfrr"
	"github.com/Jiang-Gianni/chat/message"
	"github.com/Jiang-Gianni/chat/room"
	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
)

type NATSServer struct {
	MessageQuerier message.Querier
	RoomQuerier    room.Querier
	Addr           string
	Log            *slog.Logger
	NATS           *nats.Conn
	EC             *nats.EncodedConn
}

func (n *NATSServer) Run() (rerr error) {
	defer dfrr.Wrap(&rerr, "g.Run")
	r := chi.NewRouter()
	r.Get(config.DiscardEndpoint, func(w http.ResponseWriter, r *http.Request) {})

	r.Get(config.IndexEndpoint, index())

	r.Post(config.LoginEndpoint, n.postLogin())
	r.Post(config.RegisterEndpoint, n.postRegister())
	r.Get(config.DeniedEndpoint, getDenied())
	r.Post(config.LogoutEndpoint, postLogout())

	// Grouping the endpoints that requires authentication
	r.Route(config.IndexEndpoint, func(cr chi.Router) {
		cr.Use(requireAuth())
		cr.Get(config.ChatEndpoint, n.getChat())
		cr.Get(config.ChatParamEndpoint, n.getChat())
		cr.Get(config.ChatRedirectParamEndpoint, getChatRedirect())
		cr.Post(config.RoomEndpoint, n.postRoom())
		cr.Get(config.ChatWsEndParampoint, n.getChatWs())
	})

	// Write timeout removed to support server side events
	srv := http.Server{
		Addr:        n.Addr,
		Handler:     r,
		ReadTimeout: 5 * time.Second,
		IdleTimeout: 5 * time.Second,
	}
	return srv.ListenAndServe()
}
