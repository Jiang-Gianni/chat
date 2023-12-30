package user

import (
	"context"
	"errors"
	"testing"

	"github.com/Jiang-Gianni/chat/config"
)

// Test both login and register
func TestLoginRegister(t *testing.T) {
	t.Parallel()
	db, cleanup, err := config.GetSqliteTest()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	q := New(db)

	username := "my-test-username"
	password := "my-test-password"

	ctx := context.Background()

	err = login(ctx, q, username, password)
	if !errors.Is(err, InvalidCredentialsError) {
		t.Fatal("should not be able to login without registering", err)
	}

	err = register(ctx, q, username, password)
	if err != nil {
		t.Fatal(err)
	}

	err = register(ctx, q, username, password)
	if !errors.Is(err, UsernameTakenError) {
		t.Fatal("username should already be taken ", err)
	}

	err = login(ctx, q, username, "wrong-password-test")
	if !errors.Is(err, InvalidCredentialsError) {
		t.Fatal("successful login with wrong password")
	}

	err = login(ctx, q, username, password)
	if err != nil {
		t.Fatal(err)
	}
}
