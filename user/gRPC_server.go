package user

import (
	"context"
	"fmt"
	"net"

	"github.com/Jiang-Gianni/chat/dfrr"
	grpc "google.golang.org/grpc"
)

type GRPCServer struct {
	Queries
	UnimplementedUserServer
}

// Interface Check
var (
	_ UserServer = (*GRPCServer)(nil)
	_ Querier    = (*GRPCServer)(nil)
)

func (g *GRPCServer) Register(
	ctx context.Context,
	req *RegisterRequest,
) (resp *RegisterResponse, derr error) {
	defer dfrr.Wrap(&derr, "g.Register")
	if err := register(ctx, g, req.Username, req.Password); err != nil {
		return nil, err
	}
	return &RegisterResponse{}, nil
}

func (g *GRPCServer) Login(
	ctx context.Context,
	req *LoginRequest,
) (resp *LoginResponse, derr error) {
	defer dfrr.Wrap(&derr, "g.Login")
	if err := login(ctx, g, req.Username, req.Password); err != nil {
		return nil, err
	}
	return &LoginResponse{}, nil
}

func (g *GRPCServer) RunGRPC(addr string, opts ...grpc.ServerOption) (derr error) {
	defer dfrr.Wrap(&derr, "RunGRPCServer")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("s.Serve: %w", err)
	}
	defer func(lis net.Listener) {
		if err != nil {
			derr = lis.Close()
		}
	}(lis)
	s := grpc.NewServer(opts...)
	RegisterUserServer(s, g)
	defer s.Stop()
	err = s.Serve(lis)
	if err != nil {
		return fmt.Errorf("s.Serve: %w", err)
	}
	return nil
}
