-- name: GetTaskById :one
SELECT * FROM tasks
WHERE id = $1;

-- name: ListTasks :many
SELECT * FROM tasks
WHERE 
    (
        ($1 IS NULL OR title LIKE '%' || $1 || '%')
        OR ($2 IS NULL OR description LIKE '%' || $2 || '%')
    )
    AND ($3 IS NULL OR category = $3)
    AND ($4 IS NULL OR priority = $4)
    AND ($5 IS NULL OR status = $5);

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