package message

import (
	"fmt"

	"github.com/Jiang-Gianni/chat/dfrr"
	grpc "google.golang.org/grpc"
)

type GRPCClient struct {
	Conn *grpc.ClientConn
	MessageClient
}

// Remember to close *grpc.ClientConn
func NewGRPCClient(addr string, opts ...grpc.DialOption) (c *GRPCClient, derr error) {
	defer dfrr.Wrap(&derr, "NewGRPCClient")
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return c, fmt.Errorf("grpc.Dial: %w", err)
	}
	c = &GRPCClient{
		Conn:          conn,
		MessageClient: NewMessageClient(conn),
	}
	return c, nil
}
