package user

import (
	"context"
	"testing"

	"github.com/Jiang-Gianni/chat/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGRPC(t *testing.T) {
	t.Parallel()
	db, cleanup, err := config.GetSqliteTest()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	q := New(db)

	g := GRPCServer{
		Queries: *q,
	}

	// Start the service
	go g.Run(config.UserServiceAddr)

	ctx := context.Background()
	lr := &LoginRequest{
		Username: "my-test-username",
		Password: "my-test-password",
	}

	rr := &RegisterRequest{
		Username: "my-test-username",
		Password: "my-test-password",
	}

	c, err := NewGRPCClient(config.UserServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer c.Conn.Close()

	// Invalid Credentials
	_, err = c.Login(ctx, lr)
	if err == nil {
		t.Fatal("should not be able to login without registering")
	}

	// Successful register
	_, err = c.Register(ctx, rr)
	if err != nil {
		t.Fatal("registration should have gone OK", err)
	}

	// Username Taken
	_, err = c.Register(ctx, rr)
	if err == nil {
		t.Fatal("username should already be taken")
	}

	// Wrong Password
	lr.Password = "wrong-pw"
	_, err = c.Login(ctx, lr)
	if err == nil {
		t.Fatal("status should be invalid credentials")
	}

	// Successful login
	lr.Password = "my-test-password"
	_, err = c.Login(ctx, lr)
	if err != nil {
		t.Fatal("failed to login", err)
	}
}
