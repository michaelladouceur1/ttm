-- TASKS

-- name: GetTaskById :one
SELECT id, title, description, priority, status, opened_at, closed_at, created_at, updated_at FROM tasks
WHERE id = ?;

-- name: ListTasks :many
SELECT id, title, description, priority, status, opened_at, closed_at, created_at, updated_at FROM tasks
WHERE
  (
    @priorities_json IS NULL
    OR @priorities_json = ''
    OR priority IN (SELECT value FROM json_each(@priorities_json))
  )
  AND (
    @statuses_json IS NULL
    OR @statuses_json = ''
    OR status IN (SELECT value FROM json_each(@statuses_json))
  );

-- name: CreateTask :one
INSERT INTO tasks (title, description, priority, status, opened_at, closed_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id, title, description, priority, status, opened_at, closed_at, created_at, updated_at;

-- name: UpdateTaskField :one
UPDATE tasks
SET
    title = COALESCE(?, title),
    description = COALESCE(?, description),
    priority = COALESCE(?, priority),
    status = COALESCE(?, status),
    opened_at = COALESCE(?, opened_at),
    closed_at = COALESCE(?, closed_at),
    updated_at = ?
WHERE id = ?
RETURNING id, title, description, priority, status, opened_at, closed_at, created_at, updated_at;

-- SESSIONS

-- name: CreateSession :one
INSERT INTO sessions (task_id, start_time, end_time)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetSessionsByTaskID :many
SELECT * FROM sessions
WHERE task_id = ?;

-- name: GetSessionsByTimeRange :many
SELECT * FROM sessions
WHERE start_time >= ? AND end_time <= ?;

-- TAGS

-- name: ListTags :many
SELECT tag, COUNT(DISTINCT task_id) AS task_count
FROM tags
WHERE tag IS NOT NULL AND tag <> ''
GROUP BY tag
ORDER BY LOWER(tag), tag;

-- name: GetTagsByTaskID :many
SELECT * FROM tags
WHERE task_id = ?;

-- name: CreateTag :one
INSERT INTO tags (task_id, tag)
VALUES (?, ?)
RETURNING *;

-- name: CreateTags :many
INSERT INTO tags (task_id, tag)
VALUES (?, ?)
RETURNING *;

-- name: DeleteTagsByTaskID :exec
DELETE FROM tags
WHERE task_id = ?;

-- NOTES

-- name: CreateNote :one
INSERT INTO notes (task_id, content, created_at, updated_at)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetNotesByTaskID :many
SELECT * FROM notes
WHERE task_id = ?;
