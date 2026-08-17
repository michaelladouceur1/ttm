-- name: GetTaskById :one
SELECT * FROM tasks
WHERE id = $1;

-- name: ListTasks :many
SELECT * FROM tasks
WHERE
  (
    $1 IS NULL
    OR $1 = ''
    OR category IN (SELECT value FROM json_each($1))
  )
  AND (
    $2 IS NULL
    OR $2 = ''
    OR priority IN (SELECT value FROM json_each($2))
  )
  AND (
    $3 IS NULL
    OR $3 = ''
    OR status IN (SELECT value FROM json_each($3))
  );

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

-- name: ListTags :many
SELECT tag, COUNT(DISTINCT task_id) AS task_count
FROM tags
WHERE tag IS NOT NULL AND tag <> ''
GROUP BY tag
ORDER BY LOWER(tag), tag;

-- name: GetTagsByTaskID :many
SELECT * FROM tags
WHERE task_id = $1;

-- name: CreateTag :one
INSERT INTO tags (task_id, tag)
VALUES ($1, $2)
RETURNING *;

-- name: CreateTags :many
INSERT INTO tags (task_id, tag)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteTagsByTaskID :exec
DELETE FROM tags
WHERE task_id = $1;

-- name: CreateNote :one
INSERT INTO notes (task_id, content, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetNotesByTaskID :many
SELECT * FROM notes
WHERE task_id = $1;