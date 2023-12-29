-- name: GetMessages :many
select * from message;

-- name: GetMessageByID :one
select * from message where id = sqlc.arg(id);

-- name: GetMessageByRoomID :many
select * from message where room_id = sqlc.arg(room_id) limit 10;

-- name: InsertMessage :one
insert into message(room_id, username, message, sent_at) values (
    sqlc.arg(room_id),
    sqlc.arg(username),
    sqlc.arg(message),
    sqlc.arg(sent_at)
) returning *;