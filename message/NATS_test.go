package message

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/nats-io/nats.go"
)

// Require an active NATS server running on the default URL
// The test resets the test db and assumes that only a maximum of 1 message is
// present for each roomID
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

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// Needed for test number 3
	_, err = q.InsertMessage(ctx, InsertMessageParams{
		RoomID:   int64(3),
		Username: "test-3",
		Message:  "message-for-test-3",
	})
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name    string
		subject string
		msg     Message
		roomID  int64
		count   int
	}{
		{
			name:    "Should be written",
			subject: config.NATSMessageStreamRoom(strconv.Itoa(1)),
			msg:     Message{RoomID: int64(1), Username: "test", Message: "my test message"},
			roomID:  int64(1),
			count:   1,
		},
		{
			name:    "Should NOT be written",
			subject: "not-a-valid-subject",
			msg:     Message{RoomID: int64(2), Username: "test", Message: "my test message"},
			roomID:  int64(2),
			count:   0,
		},
		{
			name:    "Should be written and return count = 2",
			subject: config.NATSMessageStreamRoom(strconv.Itoa(3)),
			msg:     Message{RoomID: int64(3), Username: "test-3", Message: "test number 3"},
			roomID:  int64(3),
			count:   2,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if err := ec.Publish(tc.subject, tc.msg); err != nil {
				t.Fatal(err)
			}

			// Wait for the service to process the message
			time.Sleep(time.Millisecond * 10)

			msg, err := q.GetMessageByRoomID(ctx, tc.roomID)
			if err != nil {
				t.Fatal(err)
			}
			if len(msg) != tc.count {
				t.Fatalf("\nexpected count %d, got %d\n", tc.count, len(msg))
			}
			// Check if same message field values
			if len(msg) == 1 {
				if msg[0].Message != tc.msg.Message {
					t.Fatalf("\nexpected message %s, got %s\n", msg[0].Message, tc.msg.Message)
				}
				if msg[0].Username != tc.msg.Username {
					t.Fatalf("\nexpected username %s, got %s\n", msg[0].Username, tc.msg.Username)
				}
				if msg[0].RoomID != tc.msg.RoomID {
					t.Fatalf("\nexpected roomID %d, got %d\n", msg[0].RoomID, tc.msg.RoomID)
				}
			}
		})
	}
}
