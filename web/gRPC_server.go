package web

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/dfrr"
	"github.com/Jiang-Gianni/chat/message"
	"github.com/Jiang-Gianni/chat/room"
	"github.com/Jiang-Gianni/chat/user"
	"github.com/go-chi/chi/v5"
)

type GRPCServer struct {
	MessageQuerier message.Querier
	RoomQuerier    room.Querier
	Addr           string
	Log            *slog.Logger
	UserClient     *user.GRPCClient
	RoomClient     *room.GRPCClient
}

func (g *GRPCServer) Run() (rerr error) {
	defer dfrr.Wrap(&rerr, "g.Run")
	r := chi.NewRouter()
	r.Get(config.DiscardEndpoint, func(w http.ResponseWriter, r *http.Request) {})

	r.Get(config.IndexEndpoint, index())

	r.Post(config.LoginEndpoint, g.postLogin())
	r.Post(config.RegisterEndpoint, g.postRegister())
	r.Get(config.DeniedEndpoint, getDenied())
	r.Post(config.LogoutEndpoint, postLogout())

	// Grouping the endpoints that requires authentication
	r.Route(config.IndexEndpoint, func(cr chi.Router) {
		cr.Use(requireAuth())
		cr.Get(config.ChatEndpoint, g.getChat())
		cr.Get(config.ChatParamEndpoint, g.getChat())
		cr.Get(config.ChatRedirectParamEndpoint, getChatRediect())
		cr.Post(config.RoomEndpoint, g.postRoom())
		cr.Get(config.ChatWsEndParampoint, g.getChatWs())
	})

	// Write timeout removed to support server side events
	srv := http.Server{
		Addr:        g.Addr,
		Handler:     r,
		ReadTimeout: 5 * time.Second,
		IdleTimeout: 5 * time.Second,
	}
	return srv.ListenAndServe()
}
