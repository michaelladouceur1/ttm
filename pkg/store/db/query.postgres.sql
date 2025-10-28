-- name: GetTaskById :one
SELECT * FROM tasks
WHERE id = $1;

-- name: ListTasks :many
SELECT * FROM tasks
WHERE 
    ($1::text IS NULL OR $1::text = '' OR title ILIKE '%' || $1::text || '%')
    AND ($2::text IS NULL OR $2::text = '' OR description ILIKE '%' || $2::text || '%')
    AND ($3::text IS NULL OR $3::text = '' OR category = $3::text)
    AND ($4::text IS NULL OR $4::text = '' OR priority = $4::text)
    AND ($5::text IS NULL OR $5::text = '' OR status = $5::text);

-- name: CreateTask :one
INSERT INTO tasks (title, description, category, priority, status, opened_at, closed_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateTaskField :one
UPDATE tasks
SET 
    title = COALESCE($1, title),
    description = COALESCE($2, description),
    category = COALESCE($3, category),
    priority = COALESCE($4, priority),
    status = COALESCE($5, status),
    opened_at = COALESCE($6, opened_at),
    closed_at = COALESCE($7, closed_at),
    updated_at = $8
WHERE id = $9
RETURNING *;

-- name: CreateSession :one
INSERT INTO sessions (task_id, start_time, end_time)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSessionsByTaskID :many
SELECT * FROM sessions
WHERE task_id = $1;

-- name: GetSessionsByTimeRange :many
SELECT * FROM sessions
WHERE start_time >= $1 AND end_time <= $2;