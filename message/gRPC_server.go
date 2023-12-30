package message

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/Jiang-Gianni/chat/dfrr"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCServer struct {
	Queries
	UnimplementedMessageServer

	// List of open streams
	RoomStreams map[string]map[*Message_StreamMessageServer]struct{}

	// Locker to avoid race conditions on RoomStreams updates
	Lock sync.Locker
}

// Interface Check
var (
	_ MessageServer = (*GRPCServer)(nil)
	_ Querier       = (*GRPCServer)(nil)
)

func (g *GRPCServer) StreamMessage(stream Message_StreamMessageServer) (rerr error) {
	defer dfrr.Wrap(&rerr, "g.StreamMessage")

	// Retrieve roomID from metadata
	ctx := stream.Context()
	md, _ := metadata.FromIncomingContext(ctx)
	roomIDString, ok := md[config.RoomIDStringKey]
	if !ok || len(roomIDString) != 1 {
		return RoomIDError
	}

	// Initialize the room connections map if not presentand add the stream
	g.Lock.Lock()
	if _, ok := g.RoomStreams[roomIDString[0]]; !ok {
		g.RoomStreams[roomIDString[0]] = map[*Message_StreamMessageServer]struct{}{}
	}
	g.RoomStreams[roomIDString[0]][&stream] = struct{}{}
	g.Lock.Unlock()

	// Remove the stream from the map once done
	defer func() {
		g.Lock.Lock()
		delete(g.RoomStreams[roomIDString[0]], &stream)
		g.Lock.Unlock()
	}()

	// Save to db and broadcast the message to the room streams
	for {
		r, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("stream.Recv: %w", err)
		}
		msg, err := g.InsertMessage(stream.Context(), InsertMessageParams{
			RoomID:   int64(r.RoomId),
			Username: r.Username,
			Message:  r.Message,
			SentAt:   time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("g.InsertMessage: %w", err)
		}
		g.BroadcastMsg(msg, roomIDString[0])
	}
}

func (g *GRPCServer) BroadcastMsg(msg Message, roomKey string) {
	for stream := range g.RoomStreams[roomKey] {
		(*stream).Send(&StreamResponse{
			Username: msg.Username,
			Message:  msg.Message,
			SentAt:   timestamppb.New(msg.SentAt),
		})
	}
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
	RegisterMessageServer(s, g)
	defer s.Stop()
	err = s.Serve(lis)
	if err != nil {
		return fmt.Errorf("s.Serve: %w", err)
	}
	return nil
}
