package main

import (
	"net/http"
	"net/url"
)

func login(username, password string) ([]*http.Cookie, error) {
	resp, err := http.PostForm(localhost+"login", url.Values{
		"username": []string{username},
		"password": []string{password},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp.Cookies(), nil
}
