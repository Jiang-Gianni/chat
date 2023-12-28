-- name: GetMessages :many
select * from message;

-- name: GetMessageByID :one
select * from message where id = ?;