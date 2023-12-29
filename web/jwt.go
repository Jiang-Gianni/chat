package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/dfrr"
	"github.com/Jiang-Gianni/chat/views"
	"github.com/golang-jwt/jwt/v5"
)

func tokenJWT(w http.ResponseWriter, userID int, username string) (rerr error) {
	defer dfrr.Wrap(&rerr, "TokenJWT")
	exp := time.Now().Add(time.Hour * 24)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userID,
		"username": username,
		"exp":      exp.Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.JWT_SIGN))
	if err != nil {
		return fmt.Errorf("token.SignedString: %w", err)
	}
	c := &http.Cookie{
		Name:     config.JWT_COOKIE,
		Value:    tokenString,
		SameSite: http.SameSiteStrictMode,
		Expires:  exp,
		HttpOnly: true,
		// TODO TLS
		// Secure: true,
	}
	http.SetCookie(w, c)
	return nil
}

func clearTokenJWT(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:     config.JWT_COOKIE,
		Value:    "",
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-time.Second),
		HttpOnly: true,
		MaxAge:   -1,
	}
	http.SetCookie(w, c)
}

func requireAuth() func(http.Handler) http.Handler {
	denied := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("HX-Redirect", config.DeniedEndpoint)
		views.WriteDeniedPage(w)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ok, userID, username := isAuth(r)
				if !ok {
					denied(w, r)
					return
				}
				ctx := context.WithValue(r.Context(), UserIDCtxKey, userID)
				ctx = context.WithValue(ctx, UsernameCtxKey, username)
				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}

func isAuth(r *http.Request) (ok bool, userID int, username string) {
	cookie, err := r.Cookie(config.JWT_COOKIE)
	if err != nil || cookie == nil {
		return false, userID, username
	}
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWT_SIGN), nil
	})
	if err != nil {
		return false, userID, username
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, userID, username
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return false, userID, username
	}
	if float64(time.Now().Unix()) > exp {
		return false, userID, username
	}
	id, ok := claims["id"].(float64)
	if !ok {
		return false, userID, username
	}
	username, ok = claims["username"].(string)
	if !ok {
		return false, userID, username
	}
	return true, int(id), username
}

func getDenied() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		views.WriteDeniedPage(w)
	}
}

func userIDUserName(r *http.Request) (int, string) {
	userID, _ := r.Context().Value(UserIDCtxKey).(int)
	username, _ := r.Context().Value(UsernameCtxKey).(string)
	return int(userID), username
}
