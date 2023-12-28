package room

import (
	"context"
	"fmt"

	"github.com/Jiang-Gianni/chat/dfrr"
)

func create(ctx context.Context, q Querier, roomName string) (roomID int64, rerr error) {
	defer dfrr.Wrap(&rerr, "create")
	roomID, err := q.InsertRoom(ctx, roomName)
	if err != nil {
		return 0, fmt.Errorf("q.InsertRoom: %w", err)
	}
	return roomID, nil
}
