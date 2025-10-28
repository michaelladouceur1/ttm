package db

import (
	"context"
	"database/sql"
	_ "embed"
	"os"
	"time"
	"ttm/pkg/models"
	"ttm/pkg/paths"

	_ "github.com/lib/pq"
)

type DBType string

const (
	Postgres DBType = "postgres"
	Sqlite   DBType = "sqlite"
)

type DBLocal struct {
	ctx    context.Context
	db     *sql.DB
	dbType DBType
}

type DBQueries interface {
	CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error)
	GetTaskById(ctx context.Context, id int64) (Task, error)
	ListTasks(ctx context.Context, arg ListTasksParams) ([]Task, error)
	UpdateTaskField(ctx context.Context, arg UpdateTaskFieldParams) (Task, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	GetSessionsByTaskID(ctx context.Context, taskID sql.NullInt64) ([]Session, error)
	GetSessionsByTimeRange(ctx context.Context, arg GetSessionsByTimeRangeParams) ([]Session, error)
}

//go:embed schema.sqlite.sql
var ddl string

func NewDBStore(dbType DBType) *DBLocal {
	return &DBLocal{dbType: dbType}
}

func (ts *DBLocal) Init() error {
	var err error

	if os.MkdirAll(paths.GetTTMDirectory(), os.ModePerm); err != nil {
		return err
	}

	ts.ctx = context.Background()

	if ts.dbType == "sqlite" {
		ts.db, err = sql.Open("sqlite3", paths.GetTaskStoreDBPath())
	} else if ts.dbType == "postgres" {
		ts.db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/ttmdb")
	} else {
		ts.db, err = sql.Open("sqlite3", paths.GetTaskStoreDBPath())
	}
	if err != nil {
		return err
	}

	if _, err := ts.db.ExecContext(ts.ctx, ddl); err != nil {
		return err
	}

	return nil
}

func (ts *DBLocal) getQueries() DBQueries {
	if ts.dbType == "postgres" {
		return NewPostgresQueriesAdapter(ts.db)
	} else if ts.dbType == "sqlite" {
		return NewSqliteQueriesAdapter(ts.db)
	}
	return NewSqliteQueriesAdapter(ts.db)
}

func (ts *DBLocal) InsertTask(task models.Task) error {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	queries := ts.getQueries()

	_, err := queries.CreateTask(ts.ctx, CreateTaskParams{
		Title:       toNullString(task.Title),
		Description: toNullString(task.Description),
		Category:    toNullString(string(task.Category)),
		Priority:    toNullString(string(task.Priority)),
		Status:      toNullString(string(task.Status)),
		OpenedAt:    toNullTime(task.OpenedAt),
		ClosedAt:    toNullTime(task.ClosedAt),
		CreatedAt:   toNullTime(task.CreatedAt),
		UpdatedAt:   toNullTime(task.UpdatedAt),
	})

	if err != nil {
		return err
	}

	return nil
}

func (ts *DBLocal) GetTaskByID(taskID int64) (models.Task, error) {
	queries := ts.getQueries()

	dbTask, err := queries.GetTaskById(ts.ctx, taskID)

	if err != nil {
		return models.Task{}, err
	}

	task := dbTasksToTasks([]Task{dbTask})

	return task[0], nil

}

func (ts *DBLocal) ListTasks(titleDescSearch string, category models.Category, status models.Status, priority models.Priority) ([]models.Task, error) {
	queries := ts.getQueries()

	dbTasks, err := queries.ListTasks(ts.ctx, ListTasksParams{
		Title:       toNullString(titleDescSearch),
		Description: toNullString(titleDescSearch),
		Category:    toNullString(string(category)),
		Priority:    toNullString(string(priority)),
		Status:      toNullString(string(status)),
	})

	if err != nil {
		return nil, err
	}

	tasks := dbTasksToTasks(dbTasks)

	return tasks, nil
}

func (ts *DBLocal) UpdateTitle(taskID int64, title string) error {
	return ts.updateTaskField(UpdateTaskFieldParams{
		ID:    taskID,
		Title: toNullString(title),
	})
}

func (ts *DBLocal) UpdateDescription(taskID int64, description string) error {
	return ts.updateTaskField(UpdateTaskFieldParams{
		ID:          taskID,
		Description: toNullString(description),
	})
}

func (ts *DBLocal) UpdateCategory(taskID int64, category models.Category) error {
	return ts.updateTaskField(UpdateTaskFieldParams{
		ID:       taskID,
		Category: toNullString(string(category)),
	})
}

func (ts *DBLocal) UpdatePriority(taskID int64, priority models.Priority) error {
	return ts.updateTaskField(UpdateTaskFieldParams{
		ID:       taskID,
		Priority: toNullString(string(priority)),
	})
}

func (ts *DBLocal) UpdateStatus(taskID int64, status models.Status) error {
	return ts.updateTaskField(UpdateTaskFieldParams{
		ID:     taskID,
		Status: toNullString(string(status)),
	})
}

func (ts *DBLocal) UpdateOpenedAt(taskID int64, openedAt time.Time) error {
	return ts.updateTaskField(UpdateTaskFieldParams{
		ID:       taskID,
		OpenedAt: toNullTime(openedAt),
	})
}

func (ts *DBLocal) UpdateClosedAt(taskID int64, closedAt time.Time) error {
	return ts.updateTaskField(UpdateTaskFieldParams{
		ID:       taskID,
		ClosedAt: toNullTime(closedAt),
	})
}

func (ts *DBLocal) updateTaskField(params UpdateTaskFieldParams) error {
	params.UpdatedAt = toNullTime(time.Now())

	queries := ts.getQueries()

	_, err := queries.UpdateTaskField(ts.ctx, params)

	if err != nil {
		return err
	}

	return nil
}

func (ts *DBLocal) AddSession(session models.Session) error {
	queries := ts.getQueries()

	_, err := queries.CreateSession(ts.ctx, CreateSessionParams{
		TaskID:    toNullInt(int(session.TaskId)),
		StartTime: toNullTime(session.StartTime),
		EndTime:   toNullTime(session.EndTime),
	})

	if err != nil {
		return err
	}

	return nil
}

func (ts *DBLocal) GetSessionsByTaskID(taskID int) ([]models.Session, error) {
	queries := ts.getQueries()

	dbSessions, err := queries.GetSessionsByTaskID(ts.ctx, toNullInt(taskID))

	sessions := []models.Session{}

	if err != nil {
		return sessions, err
	}

	if len(dbSessions) == 0 {
		return sessions, nil
	}

	for _, dbSession := range dbSessions {
		sessions = append(sessions, models.Session{
			ID:        dbSession.ID,
			TaskId:    dbSession.TaskID.Int64,
			StartTime: dbSession.StartTime.Time,
			EndTime:   dbSession.EndTime.Time,
		})
	}

	return sessions, nil

}

func (ts *DBLocal) GetSessionsByTimeRange(startTime time.Time, endTime time.Time) ([]models.Session, error) {
	queries := ts.getQueries()

	dbSessions, err := queries.GetSessionsByTimeRange(ts.ctx, GetSessionsByTimeRangeParams{
		StartTime: toNullTime(startTime),
		EndTime:   toNullTime(endTime),
	})

	sessions := []models.Session{}

	if err != nil {
		return sessions, err
	}

	if len(dbSessions) == 0 {
		return sessions, nil
	}

	for _, dbSession := range dbSessions {
		sessions = append(sessions, models.Session{
			ID:        dbSession.ID,
			TaskId:    dbSession.TaskID.Int64,
			StartTime: dbSession.StartTime.Time,
			EndTime:   dbSession.EndTime.Time,
		})
	}

	return sessions, nil

}

func toNullString(v interface{}) sql.NullString {
	switch val := v.(type) {
	case string:
		return sql.NullString{String: val, Valid: val != ""}
	case sql.NullString:
		return val
	default:
		return sql.NullString{String: "", Valid: false}
	}
}

func toNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{Time: t, Valid: true}
}

func toNullInt(i int) sql.NullInt64 {
	return sql.NullInt64{Int64: int64(i), Valid: true}
}

func dbTasksToTasks(t []Task) []models.Task {
	var tasksList []models.Task
	for _, task := range t {
		tasksList = append(tasksList, models.Task{
			ID:          task.ID,
			Title:       task.Title.String,
			Description: task.Description.String,
			Category:    models.Category(task.Category.String),
			Priority:    models.Priority(task.Priority.String),
			Status:      models.Status(task.Status.String),
			OpenedAt:    task.OpenedAt.Time,
			ClosedAt:    task.ClosedAt.Time,
			CreatedAt:   task.CreatedAt.Time,
			UpdatedAt:   task.UpdatedAt.Time,
		})
	}
	return tasksList
}
