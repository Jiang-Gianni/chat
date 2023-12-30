-- name: GetMessages :many
select * from message;

-- name: GetMessageByID :one
select * from message where id = sqlc.arg(id);

-- name: GetMessageByRoomID :many
select A.* from (
    select * from message where room_id = sqlc.arg(room_id) order by id desc limit 10
) A order by A.id asc;

-- name: InsertMessage :one
insert into message(room_id, username, message, sent_at) values (
    sqlc.arg(room_id),
    sqlc.arg(username),
    sqlc.arg(message),
    sqlc.arg(sent_at)
) returning *;