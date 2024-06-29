// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package db

import (
	"context"
	"database/sql"
)

const createSession = `-- name: CreateSession :one
INSERT INTO sessions (task_id, start_time, end_time)
VALUES (?, ?, ?)
RETURNING id, task_id, start_time, end_time
`

type CreateSessionParams struct {
	TaskID    sql.NullInt64
	StartTime sql.NullTime
	EndTime   sql.NullTime
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, createSession, arg.TaskID, arg.StartTime, arg.EndTime)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.TaskID,
		&i.StartTime,
		&i.EndTime,
	)
	return i, err
}

const createTask = `-- name: CreateTask :one
INSERT INTO tasks (title, description, category, priority, status, opened_at, closed_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id, title, description, category, priority, status, opened_at, closed_at, created_at, updated_at
`

type CreateTaskParams struct {
	Title       sql.NullString
	Description sql.NullString
	Category    sql.NullString
	Priority    sql.NullString
	Status      sql.NullString
	OpenedAt    sql.NullTime
	ClosedAt    sql.NullTime
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error) {
	row := q.db.QueryRowContext(ctx, createTask,
		arg.Title,
		arg.Description,
		arg.Category,
		arg.Priority,
		arg.Status,
		arg.OpenedAt,
		arg.ClosedAt,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Category,
		&i.Priority,
		&i.Status,
		&i.OpenedAt,
		&i.ClosedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSessionByTaskId = `-- name: GetSessionByTaskId :many
SELECT id, task_id, start_time, end_time FROM sessions
WHERE task_id = ?
`

func (q *Queries) GetSessionByTaskId(ctx context.Context, taskID sql.NullInt64) ([]Session, error) {
	rows, err := q.db.QueryContext(ctx, getSessionByTaskId, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Session
	for rows.Next() {
		var i Session
		if err := rows.Scan(
			&i.ID,
			&i.TaskID,
			&i.StartTime,
			&i.EndTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSessionsByTimeRange = `-- name: GetSessionsByTimeRange :many
SELECT id, task_id, start_time, end_time FROM sessions
WHERE start_time >= ? AND end_time <= ?
`

type GetSessionsByTimeRangeParams struct {
	StartTime sql.NullTime
	EndTime   sql.NullTime
}

func (q *Queries) GetSessionsByTimeRange(ctx context.Context, arg GetSessionsByTimeRangeParams) ([]Session, error) {
	rows, err := q.db.QueryContext(ctx, getSessionsByTimeRange, arg.StartTime, arg.EndTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Session
	for rows.Next() {
		var i Session
		if err := rows.Scan(
			&i.ID,
			&i.TaskID,
			&i.StartTime,
			&i.EndTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTaskById = `-- name: GetTaskById :one
SELECT id, title, description, category, priority, status, opened_at, closed_at, created_at, updated_at FROM tasks
WHERE id = ?
`

func (q *Queries) GetTaskById(ctx context.Context, id int64) (Task, error) {
	row := q.db.QueryRowContext(ctx, getTaskById, id)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Category,
		&i.Priority,
		&i.Status,
		&i.OpenedAt,
		&i.ClosedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listTasks = `-- name: ListTasks :many
SELECT id, title, description, category, priority, status, opened_at, closed_at, created_at, updated_at FROM tasks
WHERE 
    (
        (?1 IS NULL OR title LIKE '%' || ?1 || '%')
        OR (?2 IS NULL OR description LIKE '%' || ?2 || '%')
    )
    AND (?3 IS NULL OR category = ?3)
    AND (?4 IS NULL OR priority = ?4)
    AND (?5 IS NULL OR status = ?5)
`

type ListTasksParams struct {
	Title       interface{}
	Description interface{}
	Category    interface{}
	Priority    interface{}
	Status      interface{}
}

func (q *Queries) ListTasks(ctx context.Context, arg ListTasksParams) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, listTasks,
		arg.Title,
		arg.Description,
		arg.Category,
		arg.Priority,
		arg.Status,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Category,
			&i.Priority,
			&i.Status,
			&i.OpenedAt,
			&i.ClosedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTaskField = `-- name: UpdateTaskField :one
UPDATE tasks
SET 
    title = COALESCE(?, title),
    description = COALESCE(?, description),
    category = COALESCE(?, category),
    priority = COALESCE(?, priority),
    status = COALESCE(?, status),
    opened_at = COALESCE(?, opened_at),
    closed_at = COALESCE(?, closed_at),
    updated_at = ?
WHERE id = ?
RETURNING id, title, description, category, priority, status, opened_at, closed_at, created_at, updated_at
`

type UpdateTaskFieldParams struct {
	Title       sql.NullString
	Description sql.NullString
	Category    sql.NullString
	Priority    sql.NullString
	Status      sql.NullString
	OpenedAt    sql.NullTime
	ClosedAt    sql.NullTime
	UpdatedAt   sql.NullTime
	ID          int64
}

func (q *Queries) UpdateTaskField(ctx context.Context, arg UpdateTaskFieldParams) (Task, error) {
	row := q.db.QueryRowContext(ctx, updateTaskField,
		arg.Title,
		arg.Description,
		arg.Category,
		arg.Priority,
		arg.Status,
		arg.OpenedAt,
		arg.ClosedAt,
		arg.UpdatedAt,
		arg.ID,
	)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Category,
		&i.Priority,
		&i.Status,
		&i.OpenedAt,
		&i.ClosedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
