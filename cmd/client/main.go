package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

const localhost = "http://localhost:3000/"

func main() {
	defer func(start time.Time) {
		fmt.Println("Done sending messages")
		fmt.Println("Time since start: ", time.Since(start))
	}(time.Now())
	wg := &sync.WaitGroup{}
	for i := range [10]int{} {
		wg.Add(1)
		go loop(i+1, wg)
	}
	wg.Wait()
}

func loop(i int, wg *sync.WaitGroup) (err error) {
	defer func() {
		if err != nil {
			fmt.Println("loop ", i, " error ", err)
		}
	}()
	defer wg.Done()
	var usernamePassword = "User " + strconv.Itoa(i)
	err = register(usernamePassword, usernamePassword)
	if err != nil {
		return err
	}
	cookies, err := login(usernamePassword, usernamePassword)
	if err != nil {
		return err
	}
	err = chatws(cookies, usernamePassword)
	if err != nil {
		return err
	}
	return nil
}
