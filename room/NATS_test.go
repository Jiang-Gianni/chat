package room

import (
	"context"
	"testing"
	"time"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nuid"
)

func TestNATS(t *testing.T) {
	t.Parallel()
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatal(err)
	}
	db, cleanup, err := config.GetSqliteTest()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	q := *New(db)
	ns := NATSServer{
		NATS:    nc,
		Queries: q,
	}

	// Start the service
	go ns.Run()

	ctx := context.Background()

	// A new NATS connection has to be instantiated
	// If both client and server share the same connection
	// then the same connection is publishing and subscribing
	// to the same topic -> no message is received
	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatal(err)
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		t.Fatal(err)
	}

	cen := CreateEventNATS{
		RoomName: "my-test-room",
		ReplyTo:  nuid.Next(),
	}

	err = ec.Publish(config.NATSRoomCreate, cen)
	if err != nil {
		t.Fatal(err)
	}
	reply := &CreateReplyNATS{StatusCode: StatusOK}
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	sub, err := ec.Subscribe(cen.ReplyTo, func(lrn *CreateReplyNATS) {
		reply = lrn
		cancel()
	})
	defer sub.Unsubscribe()
	if err != nil {
		t.Fatal(err)
	}
	<-ctxTimeout.Done()
	if reply.RoomID == 0 {
		t.Fatalf("expected an auto assigned room id but got 0")
	}
	if reply.StatusCode != StatusOK {
		t.Fatalf("status ok is not %d but: %d", StatusOK, reply.StatusCode)
	}
	r, err := q.GetRooms(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(r) != 1 {
		t.Fatalf("expected 1 room, got %d", len(r))
	}
}
