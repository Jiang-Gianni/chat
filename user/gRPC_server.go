package user

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/Jiang-Gianni/chat/dfrr"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
	err := register(ctx, g, req.Username, req.Password)
	if errors.Is(err, UsernameTakenError) {
		return nil, status.Errorf(codes.AlreadyExists, err.Error())
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}

func (g *GRPCServer) Login(
	ctx context.Context,
	req *LoginRequest,
) (resp *LoginResponse, derr error) {
	defer dfrr.Wrap(&derr, "g.Login")
	resp = &LoginResponse{}
	err := login(ctx, g, req.Username, req.Password)
	if errors.Is(err, InvalidCredentialsError) {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}

func (g *GRPCServer) Run(addr string, opts ...grpc.ServerOption) (derr error) {
	defer dfrr.Wrap(&derr, "g.Run")
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
