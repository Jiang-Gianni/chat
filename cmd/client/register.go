package main

import (
	"net/http"
	"net/url"
)

func register(username, password string) error {
	resp, err := http.PostForm(localhost+"register", url.Values{
		"username": []string{username},
		"password": []string{password},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
