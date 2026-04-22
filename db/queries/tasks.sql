-- name: GetTasks :many
select * from tasks where user_id = $1;

-- name: GetTasksByStatus :many
select * from tasks where user_id = $1 and is_completed = $2;

-- name: CreateTask :exec
insert into tasks(user_id, content) values($1, $2);

-- name: UpdateTaskStatus :exec
update tasks set is_completed = $2 where id = $1;

-- name: DeleteTask :exec
delete from tasks where id = $1;

