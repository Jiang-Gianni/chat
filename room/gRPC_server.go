package room

import (
	"context"
	"fmt"
	"net"

	"github.com/Jiang-Gianni/chat/dfrr"
	grpc "google.golang.org/grpc"
)

type GRPCServer struct {
	Queries
	UnimplementedRoomServer
}

// Interface Check
var (
	_ RoomServer = (*GRPCServer)(nil)
	_ Querier    = (*GRPCServer)(nil)
)

func (g *GRPCServer) Create(
	ctx context.Context,
	req *CreateRequest,
) (resp *CreateResponse, rerr error) {
	defer dfrr.Wrap(&rerr, "g.Create")
	roomID, err := create(ctx, g, req.RoomName)
	if err != nil {
		return nil, err
	}
	return &CreateResponse{RoomId: int32(roomID)}, nil
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
	RegisterRoomServer(s, g)
	defer s.Stop()
	err = s.Serve(lis)
	if err != nil {
		return fmt.Errorf("s.Serve: %w", err)
	}
	return nil
}
