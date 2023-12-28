package user

import (
	"context"
	"fmt"

	"github.com/Jiang-Gianni/chat/dfrr"
)

func login(ctx context.Context, q Querier, username, password string) (rerr error) {
	defer dfrr.Wrap(&rerr, "login")
	user, err := q.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("q.GetUser: %w", err)
	}
	if user.Password != password {
		return InvalidCredentialsError
	}
	return nil
}

func register(ctx context.Context, q Querier, username, password string) (rerr error) {
	defer dfrr.Wrap(&rerr, "register")
	count, err := q.CountUser(ctx, username)
	if err != nil {
		return fmt.Errorf("q.CountUser: %w", err)
	}
	if int(count) > 0 {
		return UsernameTakenError
	}
	err = q.InsertUser(ctx, InsertUserParams{
		Username: username,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("q.InsertUser: %w", err)
	}
	return nil
}
