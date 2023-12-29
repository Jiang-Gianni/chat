package web

import (
	"net/http"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/user"
	"github.com/Jiang-Gianni/chat/views"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *GRPCServer) postLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password := r.FormValue("username"), r.FormValue("password")
		if username == "" || password == "" {
			views.WriteLoginRegisterError(w, "Invalid username or password")
			return
		}
		req := &user.LoginRequest{
			Username: username,
			Password: password,
		}
		resp, err := g.UserClient.Login(r.Context(), req)
		if err != nil {
			g.Log.Info(err.Error(), "service", "web")
			if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
				views.WriteLoginRegisterError(w, "Invalid credentials")
				return
			}
			views.WriteLoginRegisterError(w, "Internal Server Error")
			return
		}
		if err := tokenJWT(w, int(resp.UserId), username); err != nil {
			g.Log.Error(err.Error(), "service", "web")
			views.WriteLoginRegisterError(w, "Internal Server Error")
		}
		w.Header().Add("HX-Redirect", config.ChatEndpoint)
	}
}

func (g *GRPCServer) postRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password := r.FormValue("username"), r.FormValue("password")
		if username == "" || password == "" {
			views.WriteLoginRegisterError(w, "Invalid username or password")
			return
		}
		req := &user.RegisterRequest{
			Username: username,
			Password: password,
		}
		resp, err := g.UserClient.Register(r.Context(), req)
		if err != nil {
			g.Log.Info(err.Error(), "service", "web")
			if status, ok := status.FromError(err); ok && status.Code() == codes.AlreadyExists {
				views.WriteLoginRegisterError(w, "Username already taken")
				return
			}
			views.WriteLoginRegisterError(w, "Internal Server Error")
			return
		}
		if err := tokenJWT(w, int(resp.UserId), username); err != nil {
			g.Log.Error(err.Error(), "service", "web")
			views.WriteLoginRegisterError(w, "Internal Server Error")
		}
		w.Header().Add("HX-Redirect", config.ChatEndpoint)
	}
}
