// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package message

import (
	"time"
)

type Message struct {
	ID       int64     `json:"id"`
	RoomID   int64     `json:"room_id"`
	Username string    `json:"username"`
	Message  string    `json:"message"`
	SentAt   time.Time `json:"sent_at"`
}