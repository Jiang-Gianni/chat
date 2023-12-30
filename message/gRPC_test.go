package message

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Jiang-Gianni/chat/config"
	"golang.org/x/sync/errgroup"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func testClientStream(t *testing.T, ctx context.Context) (Message_StreamMessageClient, func()) {
	c, err := NewGRPCClient(
		config.MessageServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	s, err := c.StreamMessage(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return s, func() {
		c.Conn.Close()
		s.CloseSend()
	}
}

// The test verifies that the sent message is stored in db and that
// it is broadcasted to the clients in the same room
func TestGRPC(t *testing.T) {
	t.Parallel()
	db, cleanup, err := config.GetSqliteTest()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	q := New(db)

	g := GRPCServer{
		Queries:     *q,
		RoomStreams: map[string]map[*Message_StreamMessageServer]struct{}{},
		Lock:        &sync.Mutex{},
	}

	// Start the service
	go g.Run(config.MessageServiceAddr)

	// Client without room_id metadata > should give error
	client, err := NewGRPCClient(
		config.MessageServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Conn.Close()
	stream, err := client.StreamMessage(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	// Time for the server service to check the missing room_id
	// The function retunrs RoomIDError but the error is EOF
	// So I presume that the stream connection gets closed
	time.Sleep(10 * time.Millisecond)
	err = stream.Send(&StreamRequest{})
	if err == nil {
		t.Fatal("room id is missing, err should not be nil")
	}

	// Client 1 and Client 2 in room 1
	ctxRoom1 := context.Background()
	ctxRoom1 = metadata.AppendToOutgoingContext(ctxRoom1, config.RoomIDStringKey, "1")
	stream1, cleanup1 := testClientStream(t, ctxRoom1)
	defer cleanup1()
	stream2, cleanup2 := testClientStream(t, ctxRoom1)
	defer cleanup2()

	msgReq1 := &StreamRequest{RoomId: int32(1), Username: "client-1", Message: "message-1"}
	err = stream1.Send(msgReq1)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for the service to save the message to database
	time.Sleep(time.Millisecond * 10)
	msg, err := q.GetMessageByRoomID(ctxRoom1, int64(1))
	if err != nil {
		t.Fatal(err)
	}
	if len(msg) != 1 {
		t.Fatalf("room 1 expected 1 message, got %d", len(msg))
	}
	if msg[0].Message != msgReq1.Message {
		t.Fatalf("room 1 expected message: %s\n got: %s", msgReq1.Message, msg[0].Message)
	}
	if msg[0].Username != msgReq1.Username {
		t.Fatalf("room 1 expected username: %s\n got: %s", msgReq1.Message, msg[0].Message)
	}

	// Client 1 and 2 should both be broadcasted the same message
	respCh1 := [2]*StreamResponse{}
	eg := errgroup.Group{}
	eg.Go(func() error {
		respCh1[0], err = stream1.Recv()
		return err
	})
	eg.Go(func() error {
		respCh1[1], err = stream2.Recv()
		return err
	})
	err = eg.Wait()
	if err != nil {
		t.Fatal(err)
	}
	if respCh1[0].Message != respCh1[1].Message {
		t.Fatalf(
			"room 1 clients got broadcasted different messages: %s, %s",
			respCh1[0].Message,
			respCh1[1].Message,
		)
	}
	if respCh1[0].Username != respCh1[1].Username {
		t.Fatalf(
			"room 1 clients got broadcasted different usernames: %s, %s",
			respCh1[0].Username,
			respCh1[1].Username,
		)
	}
}
