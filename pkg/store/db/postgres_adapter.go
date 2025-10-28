package db

import (
	"context"
	"database/sql"
	ttmpostgres "ttm/pkg/store/db/postgres"
)

type PostgresQueriesAdapter struct {
	q *ttmpostgres.Queries
}

func NewPostgresQueriesAdapter(db *sql.DB) *PostgresQueriesAdapter {
	return &PostgresQueriesAdapter{q: ttmpostgres.New(db)}
}

func (a *PostgresQueriesAdapter) CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error) {
	pgArg := ttmpostgres.CreateTaskParams(arg)
	pgTask, err := a.q.CreateTask(ctx, pgArg)
	if err != nil {
		return Task{}, err
	}
	return Task(pgTask), nil
}

func (a *PostgresQueriesAdapter) GetTaskById(ctx context.Context, id int64) (Task, error) {
	pgTask, err := a.q.GetTaskById(ctx, id)
	if err != nil {
		return Task{}, err
	}
	return Task(pgTask), nil
}

func (a *PostgresQueriesAdapter) ListTasks(ctx context.Context, arg ListTasksParams) ([]Task, error) {
	pgArg := listTasksParamsToPgParams(arg)
	pgTasks, err := a.q.ListTasks(ctx, pgArg)
	if err != nil {
		return nil, err
	}
	tasks := make([]Task, len(pgTasks))
	for i, pgTask := range pgTasks {
		tasks[i] = Task(pgTask)
	}
	return tasks, nil
}

func (a *PostgresQueriesAdapter) UpdateTaskField(ctx context.Context, arg UpdateTaskFieldParams) (Task, error) {
	pgArg := updateTaskFieldParamsToPgParams(arg)
	pgTask, err := a.q.UpdateTaskField(ctx, pgArg)
	if err != nil {
		return Task{}, err
	}
	return Task(pgTask), nil
}

func (a *PostgresQueriesAdapter) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	pgArg := createSessionParamsToPgParams(arg)
	pgSession, err := a.q.CreateSession(ctx, pgArg)
	if err != nil {
		return Session{}, err
	}
	return Session(pgSession), nil
}

func (a *PostgresQueriesAdapter) GetSessionsByTaskID(ctx context.Context, taskID sql.NullInt64) ([]Session, error) {
	pgSessions, err := a.q.GetSessionsByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	sessions := make([]Session, len(pgSessions))
	for i, pgSession := range pgSessions {
		sessions[i] = Session(pgSession)
	}
	return sessions, nil
}

func (a *PostgresQueriesAdapter) GetSessionsByTimeRange(ctx context.Context, arg GetSessionsByTimeRangeParams) ([]Session, error) {
	pgArg := ttmpostgres.GetSessionsByTimeRangeParams{
		StartTime: arg.StartTime,
		EndTime:   arg.EndTime,
	}
	pgSessions, err := a.q.GetSessionsByTimeRange(ctx, pgArg)
	if err != nil {
		return nil, err
	}
	sessions := make([]Session, len(pgSessions))
	for i, pgSession := range pgSessions {
		sessions[i] = Session(pgSession)
	}
	return sessions, nil
}

func listTasksParamsToPgParams(arg ListTasksParams) ttmpostgres.ListTasksParams {
	return ttmpostgres.ListTasksParams{
		Column1: toNullString(arg.Title).String,
		Column2: toNullString(arg.Description).String,
		Column3: toNullString(arg.Category).String,
		Column4: toNullString(arg.Priority).String,
		Column5: toNullString(arg.Status).String,
	}
}

func updateTaskFieldParamsToPgParams(arg UpdateTaskFieldParams) ttmpostgres.UpdateTaskFieldParams {
	return ttmpostgres.UpdateTaskFieldParams{
		ID:          arg.ID,
		Title:       arg.Title,
		Description: arg.Description,
		Category:    arg.Category,
		Priority:    arg.Priority,
		Status:      arg.Status,
		OpenedAt:    arg.OpenedAt,
		ClosedAt:    arg.ClosedAt,
		UpdatedAt:   arg.UpdatedAt,
	}
}

func createSessionParamsToPgParams(arg CreateSessionParams) ttmpostgres.CreateSessionParams {
	return ttmpostgres.CreateSessionParams{
		TaskID:    arg.TaskID,
		StartTime: arg.StartTime,
		EndTime:   arg.EndTime,
	}
}
