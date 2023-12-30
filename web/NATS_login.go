package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/user"
	"github.com/Jiang-Gianni/chat/views"
	"github.com/nats-io/nuid"
)

func (n *NATSServer) postLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password := r.FormValue("username"), r.FormValue("password")
		if username == "" || password == "" {
			views.WriteLoginRegisterError(w, InvalidUsernamePassword)
			return
		}
		len := user.LoginEventNATS{
			Username: username,
			Password: password,
			ReplyTo:  nuid.New().Next(),
		}
		err := n.EC.Publish(config.NATSUserLogin, len)
		if err != nil {
			n.Log.Error(fmt.Sprintf("ec.Publish: %s", err), "service", "web", "post", "login")
			views.WriteLoginRegisterError(w, InternalServerError)
			return
		}
		reply := &user.LoginReplyNATS{StatusCode: user.StatusOK}
		ctx, cancel := context.WithTimeout(r.Context(), time.Second)
		defer cancel()
		sub, err := n.EC.Subscribe(len.ReplyTo, func(lrn *user.LoginReplyNATS) {
			reply = lrn
			cancel()
		})
		if err != nil {
			n.Log.Error(fmt.Sprintf("ec.Subscribe: %s", err), "service", "web", "post", "login")
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
					"login",
				)
			}
		}()
		<-ctx.Done()
		n.Log.Info(
			"login",
			"service",
			"web",
			"username",
			username,
			"code",
			reply.StatusCode,
			"post",
			"login",
		)
		switch reply.StatusCode {
		case user.StatusInvalidCredentials:
			views.WriteLoginRegisterError(w, InvalidCredentials)
			return
		case user.StatusInternalError:
			views.WriteLoginRegisterError(w, InternalServerError)
			return
		default:
			if err := tokenJWT(w, username); err != nil {
				n.Log.Error(err.Error(), "service", "web", "post", "login")
				views.WriteLoginRegisterError(w, InternalServerError)
				return
			}
			w.Header().Add("HX-Redirect", config.ChatEndpoint)
			return
		}
	}
}

func (n *NATSServer) postRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password := r.FormValue("username"), r.FormValue("password")
		if username == "" || password == "" {
			views.WriteLoginRegisterError(w, InvalidUsernamePassword)
			return
		}
		ren := user.RegisterEventNATS{
			Username: username,
			Password: password,
			ReplyTo:  nuid.New().Next(),
		}
		err := n.EC.Publish(config.NATSUserRegister, ren)
		if err != nil {
			n.Log.Error(fmt.Sprintf("ec.Publish: %s", err), "service", "web", "post", "register")
			views.WriteLoginRegisterError(w, InternalServerError)
			return
		}
		reply := &user.RegisterReplyNATS{StatusCode: user.StatusOK}
		ctx, cancel := context.WithTimeout(r.Context(), time.Second)
		defer cancel()
		sub, err := n.EC.Subscribe(ren.ReplyTo, func(lrn *user.RegisterReplyNATS) {
			reply = lrn
			cancel()
		})
		if err != nil {
			n.Log.Error(fmt.Sprintf("ec.Subscribe: %s", err), "service", "web", "post", "register")
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
					"register",
				)
			}
		}()
		<-ctx.Done()
		n.Log.Info(
			"login",
			"service",
			"web",
			"username",
			username,
			"code",
			reply.StatusCode,
			"post",
			"register",
		)
		switch reply.StatusCode {
		case user.StatusUsernameTaken:
			views.WriteLoginRegisterError(w, UsernameAlreadyTaken)
			return
		case user.StatusInternalError:
			views.WriteLoginRegisterError(w, InternalServerError)
			return
		default:
			if err := tokenJWT(w, username); err != nil {
				n.Log.Error(err.Error(), "service", "web", "post", "login")
				views.WriteLoginRegisterError(w, InternalServerError)
				return
			}
			w.Header().Add("HX-Redirect", config.ChatEndpoint)
			return
		}
	}
}
