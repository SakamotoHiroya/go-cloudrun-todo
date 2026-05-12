-- name: GetUser :one
select * from users where id = $1;

-- name: GetUserByGoogleSub :one
select * from users where google_sub = $1;

-- name: CreateUserByGoogleSub :one
insert into users (google_sub, name)
values ($1, $2)
returning *;
