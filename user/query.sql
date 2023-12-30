-- name: GetUser :one
select * from user where username = sqlc.arg(username);

-- name: CountUser :one
select count(*) from user where username = sqlc.arg(username);

-- name: InsertUser :exec
insert into user(username, password, last_room_id) values (sqlc.arg(username), sqlc.arg(password), -1);