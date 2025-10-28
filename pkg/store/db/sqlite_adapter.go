package db

import (
	"context"
	"database/sql"
	ttmsqlite "ttm/pkg/store/db/sqlite"
)

type SqliteQueriesAdapter struct {
	q *ttmsqlite.Queries
}

func NewSqliteQueriesAdapter(db *sql.DB) *SqliteQueriesAdapter {
	return &SqliteQueriesAdapter{q: ttmsqlite.New(db)}
}

func (a *SqliteQueriesAdapter) CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error) {
	sqliteArg := ttmsqlite.CreateTaskParams(arg)
	sqliteTask, err := a.q.CreateTask(ctx, sqliteArg)
	if err != nil {
		return Task{}, err
	}
	return Task(sqliteTask), nil
}

func (a *SqliteQueriesAdapter) GetTaskById(ctx context.Context, id int64) (Task, error) {
	sqliteTask, err := a.q.GetTaskById(ctx, id)
	if err != nil {
		return Task{}, err
	}
	return Task(sqliteTask), nil
}

func (a *SqliteQueriesAdapter) ListTasks(ctx context.Context, arg ListTasksParams) ([]Task, error) {
	sqliteArg := ttmsqlite.ListTasksParams(arg)
	sqliteTasks, err := a.q.ListTasks(ctx, sqliteArg)
	if err != nil {
		return nil, err
	}
	tasks := make([]Task, len(sqliteTasks))
	for i, sqliteTask := range sqliteTasks {
		tasks[i] = Task(sqliteTask)
	}
	return tasks, nil
}

func (a *SqliteQueriesAdapter) UpdateTaskField(ctx context.Context, arg UpdateTaskFieldParams) (Task, error) {
	sqliteArg := ttmsqlite.UpdateTaskFieldParams(arg)
	sqliteTask, err := a.q.UpdateTaskField(ctx, sqliteArg)
	if err != nil {
		return Task{}, err
	}
	return Task(sqliteTask), nil
}

func (a *SqliteQueriesAdapter) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	sqliteArg := ttmsqlite.CreateSessionParams(arg)
	sqliteSession, err := a.q.CreateSession(ctx, sqliteArg)
	if err != nil {
		return Session{}, err
	}
	return Session(sqliteSession), nil
}

func (a *SqliteQueriesAdapter) GetSessionsByTaskID(ctx context.Context, taskID sql.NullInt64) ([]Session, error) {
	sqliteSessions, err := a.q.GetSessionsByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	sessions := make([]Session, len(sqliteSessions))
	for i, sqliteSession := range sqliteSessions {
		sessions[i] = Session(sqliteSession)
	}
	return sessions, nil
}

func (a *SqliteQueriesAdapter) GetSessionsByTimeRange(ctx context.Context, arg GetSessionsByTimeRangeParams) ([]Session, error) {
	sqliteArg := ttmsqlite.GetSessionsByTimeRangeParams{
		StartTime: arg.StartTime,
		EndTime:   arg.EndTime,
	}
	sqliteSessions, err := a.q.GetSessionsByTimeRange(ctx, sqliteArg)
	if err != nil {
		return nil, err
	}
	sessions := make([]Session, len(sqliteSessions))
	for i, sqliteSession := range sqliteSessions {
		sessions[i] = Session(sqliteSession)
	}
	return sessions, nil
}
