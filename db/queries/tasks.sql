-- name: GetTasks :many
select id, user_id, title, description, is_completed, created_at, updated_at
from tasks
where user_id = $1
order by created_at desc;

-- name: GetTasksByStatus :many
select id, user_id, title, description, is_completed, created_at, updated_at
from tasks
where user_id = $1 and is_completed = $2
order by created_at desc;

-- name: CreateTask :one
insert into tasks (user_id, title, description)
values ($1, $2, $3)
returning id, user_id, title, description, is_completed, created_at, updated_at;

-- name: UpdateTaskStatus :one
update tasks
set is_completed = $2, updated_at = now()
where id = $1 and user_id = $3
returning id, user_id, title, description, is_completed, created_at, updated_at;

-- name: DeleteTask :execrows
delete from tasks where id = $1 and user_id = $2;
