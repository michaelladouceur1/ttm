package db

import (
	"context"
	"database/sql"
	_ "embed"
	"os"
	"time"
	"ttm/pkg/models"
	"ttm/pkg/paths"
)

type DBStore struct {
	ctx context.Context
	db  *sql.DB
}

//go:embed schema.sql
var ddl string

func NewDBStore() *DBStore {
	return &DBStore{}
}

func (ts *DBStore) Init() error {
	var err error

	if os.MkdirAll(paths.GetTaskStorePath(), os.ModePerm); err != nil {
		return err
	}

	ts.ctx = context.Background()

	ts.db, err = sql.Open("sqlite3", paths.GetTaskStoreDBPath())
	if err != nil {
		return err
	}

	if _, err := ts.db.ExecContext(ts.ctx, ddl); err != nil {
		return err
	}

	return nil
}

func (ts *DBStore) InsertTask(task models.Task) error {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	queries := New(ts.db)

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

func (ts *DBStore) ListTasks(titleDescSearch string, category models.Category, status models.Status, priority models.Priority) ([]models.Task, error) {
	queries := New(ts.db)

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

func (ts *DBStore) UpdateTitle(taskID int64, title string) error {
	return ts.updateTaskField(UpdateTaskFieldParams{
		ID:    taskID,
		Title: toNullString(title),
	})
}

func (ts *DBStore) UpdateDescription(taskID int64, description string) error {
	return ts.updateTaskField(UpdateTaskFieldParams{
		ID:          taskID,
		Description: toNullString(description),
	})
}

func (ts *DBStore) UpdateCategory(taskID int64, category models.Category) error {
	return ts.updateTaskField(UpdateTaskFieldParams{
		ID:       taskID,
		Category: toNullString(string(category)),
	})
}

func (ts *DBStore) UpdatePriority(taskID int64, priority models.Priority) error {
	return ts.updateTaskField(UpdateTaskFieldParams{
		ID:       taskID,
		Priority: toNullString(string(priority)),
	})
}

func (ts *DBStore) UpdateStatus(taskID int64, status models.Status) error {
	return ts.updateTaskField(UpdateTaskFieldParams{
		ID:     taskID,
		Status: toNullString(string(status)),
	})
}

func (ts *DBStore) UpdateOpenedAt(taskID int64, openedAt time.Time) error {
	return ts.updateTaskField(UpdateTaskFieldParams{
		ID:       taskID,
		OpenedAt: toNullTime(openedAt),
	})
}

func (ts *DBStore) UpdateClosedAt(taskID int64, closedAt time.Time) error {
	return ts.updateTaskField(UpdateTaskFieldParams{
		ID:       taskID,
		ClosedAt: toNullTime(closedAt),
	})
}

func (ts *DBStore) updateTaskField(params UpdateTaskFieldParams) error {
	params.UpdatedAt = toNullTime(time.Now())

	queries := New(ts.db)

	_, err := queries.UpdateTaskField(ts.ctx, params)

	if err != nil {
		return err
	}

	return nil
}

func (ts *DBStore) AddSession(session models.Session) error {
	queries := New(ts.db)

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

func (ts *DBStore) GetSessionByTaskID(taskID int) ([]models.Session, error) {
	queries := New(ts.db)

	dbSessions, err := queries.GetSessionByTaskId(ts.ctx, toNullInt(taskID))

	sessions := []models.Session{}

	if err != nil {
		return sessions, err
	}

	if len(dbSessions) == 0 {
		return sessions, nil
	}

	for _, dbSession := range dbSessions {
		sessions = append(sessions, models.Session{
			TaskId:    dbSession.ID,
			StartTime: dbSession.StartTime.Time,
			EndTime:   dbSession.EndTime.Time,
		})
	}

	return sessions, nil

}

func toNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
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
