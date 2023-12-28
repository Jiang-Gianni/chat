-- name: GetRooms :many
select * from room;

-- name: InsertRoom :one
insert into room(name) values (sqlc.arg(room_name)) returning id;