package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/gorilla/websocket"
)

func chatws(cookies []*http.Cookie, user string) error {
	URL := &url.URL{Scheme: "ws", Host: "localhost:3000", Path: "/chat/1/ws"}
	HTTPURL := &url.URL{Scheme: "http", Host: "localhost:3000", Path: "/chat/1/ws"}
	cj, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	wsd := websocket.Dialer{
		Jar: cj,
	}
	wsd.Jar.SetCookies(HTTPURL, cookies)
	conn, _, err := wsd.Dial(URL.String(), nil)
	if err != nil {
		return err
	}
	for i := range [200]int{} {
		// time.Sleep(time.Millisecond * 60)
		emoji := "👋"
		err := conn.WriteJSON(map[string]string{
			"message": fmt.Sprintf("Message n %d from %s %s", i+1, user, emoji),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
