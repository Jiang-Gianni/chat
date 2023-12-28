// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: query.sql

package room

import (
	"context"
)

const getRooms = `-- name: GetRooms :many
select id, name from room
`

func (q *Queries) GetRooms(ctx context.Context) ([]Room, error) {
	rows, err := q.db.QueryContext(ctx, getRooms)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Room
	for rows.Next() {
		var i Room
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertRoom = `-- name: InsertRoom :one
insert into room(name) values (?1) returning id
`

func (q *Queries) InsertRoom(ctx context.Context, roomName string) (int64, error) {
	row := q.db.QueryRowContext(ctx, insertRoom, roomName)
	var id int64
	err := row.Scan(&id)
	return id, err
}