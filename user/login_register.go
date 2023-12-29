package user

import (
	"context"
	"fmt"

	"github.com/Jiang-Gianni/chat/dfrr"
	"golang.org/x/crypto/bcrypt"
)

func login(ctx context.Context, q Querier, username, password string) (userID int, rerr error) {
	defer dfrr.Wrap(&rerr, "login")
	count, err := q.CountUser(ctx, username)
	if err != nil {
		return userID, fmt.Errorf("q.CountUser: %w", err)
	}
	if int(count) == 0 {
		return userID, InvalidCredentialsError
	}
	user, err := q.GetUser(ctx, username)
	if err != nil {
		return userID, fmt.Errorf("q.GetUser: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return userID, InvalidCredentialsError
	}
	return int(user.ID), nil
}

func register(ctx context.Context, q Querier, username, password string) (userID int, rerr error) {
	defer dfrr.Wrap(&rerr, "register")
	count, err := q.CountUser(ctx, username)
	if err != nil {
		return userID, fmt.Errorf("q.CountUser: %w", err)
	}
	if int(count) > 0 {
		return userID, UsernameTakenError
	}
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return userID, fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}
	id, err := q.InsertUser(ctx, InsertUserParams{
		Username: username,
		Password: string(hashedPw),
	})
	if err != nil {
		return userID, fmt.Errorf("q.InsertUser: %w", err)
	}
	return int(id), nil
}
