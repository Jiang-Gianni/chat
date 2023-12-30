package room

import (
	"context"
	"testing"

	"github.com/Jiang-Gianni/chat/config"
	grpc "google.golang.org/grpc"
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
	go g.Run(config.RoomServiceAddr)

	ctx := context.Background()
	cr := &CreateRequest{
		RoomName: "grpc-room-test",
	}

	c, err := NewGRPCClient(config.RoomServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer c.Conn.Close()

	resp, err := c.Create(ctx, cr)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.RoomId == 0 {
		t.Fatal("resp is nil or roomID is zero")
	}
}
