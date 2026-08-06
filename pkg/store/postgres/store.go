package ttmpostgres

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"time"
	"ttm/pkg/models"
	"ttm/pkg/paths"

	_ "github.com/lib/pq"
)

type DBLocal struct {
	ctx context.Context
	db  *sql.DB
}

//go:embed schema.postgres.sql
var ddl string

func NewStore() *DBLocal {
	return &DBLocal{}
}

func (ts *DBLocal) Init() error {
	var err error

	if os.MkdirAll(paths.GetTTMDirectory(), os.ModePerm); err != nil {
		return err
	}

	ts.ctx = context.Background()

	ts.db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/ttmdb")
	if err != nil {
		return err
	}

	if _, err := ts.db.ExecContext(ts.ctx, ddl); err != nil {
		return err
	}

	return nil
}

func (ts *DBLocal) InsertTask(task models.Task) (models.Task, error) {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	queries := New(ts.db)

	newTask, err := queries.CreateTask(ts.ctx, CreateTaskParams{
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
		return models.Task{}, err
	}

	return dbTaskToTask(newTask), nil
}

func (ts *DBLocal) GetTaskByID(taskID int64) (models.Task, error) {
	queries := New(ts.db)

	dbTask, err := queries.GetTaskById(ts.ctx, taskID)

	if err != nil {
		return models.Task{}, err
	}

	task := dbTasksToTasks([]Task{dbTask})

	return task[0], nil

}

func (ts *DBLocal) ListTasks(titleDescSearch string, category models.Category, status models.Status, priority models.Priority) ([]models.Task, error) {
	queries := New(ts.db)

	dbTasks, err := queries.ListTasks(ts.ctx, ListTasksParams{
		Column1: toNullString(titleDescSearch).String,
		Column2: toNullString(titleDescSearch).String,
		Column3: toNullString(string(category)).String,
		Column4: toNullString(string(priority)).String,
		Column5: toNullString(string(status)).String,
	})

	if err != nil {
		return nil, err
	}

	tasks := dbTasksToTasks(dbTasks)

	return tasks, nil
}

func (ts *DBLocal) SearchTasks(search models.TaskSearch) ([]models.Task, error) {
	query := `SELECT t.id, t.title, t.description, t.category, t.priority, t.status, t.opened_at, t.closed_at, t.created_at, t.updated_at FROM tasks t`
	clauses := []string{}
	args := []any{}
	placeholder := func(value string) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}
	addMatch := func(column, value string) {
		if value == "" {
			return
		}
		clauses = append(clauses, "COALESCE(t."+column+", '') ILIKE '%' || "+placeholder(value)+" || '%'")
	}
	addTagMatch := func(value string) string {
		return "EXISTS (SELECT 1 FROM tags tag WHERE tag.task_id = t.id AND COALESCE(tag.tag, '') ILIKE '%' || " + placeholder(value) + " || '%')"
	}

	if search.General != "" {
		title := placeholder(search.General)
		description := placeholder(search.General)
		category := placeholder(search.General)
		priority := placeholder(search.General)
		status := placeholder(search.General)
		tag := placeholder(search.General)
		clauses = append(clauses, "(COALESCE(t.title, '') ILIKE '%' || "+title+" || '%' OR "+
			"COALESCE(t.description, '') ILIKE '%' || "+description+" || '%' OR "+
			"COALESCE(t.category, '') ILIKE '%' || "+category+" || '%' OR "+
			"COALESCE(t.priority, '') ILIKE '%' || "+priority+" || '%' OR "+
			"COALESCE(t.status, '') ILIKE '%' || "+status+" || '%' OR "+
			"EXISTS (SELECT 1 FROM tags tag WHERE tag.task_id = t.id AND COALESCE(tag.tag, '') ILIKE '%' || "+tag+" || '%'))")
	}
	addMatch("title", search.Title)
	addMatch("description", search.Description)
	addMatch("category", search.Category)
	addMatch("priority", search.Priority)
	addMatch("status", search.Status)
	for _, tag := range search.Tags {
		clauses = append(clauses, addTagMatch(tag))
	}
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	rows, err := ts.db.QueryContext(ts.ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dbTasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Category, &task.Priority, &task.Status, &task.OpenedAt, &task.ClosedAt, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		dbTasks = append(dbTasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return dbTasksToTasks(dbTasks), nil
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

func (ts *DBLocal) UpdateTags(taskID int64, tags []string) error {
	queries := New(ts.db)

	err := queries.DeleteTagsByTaskID(ts.ctx, toNullInt(int(taskID)))
	if err != nil {
		return err
	}

	for _, tag := range tags {
		_, err := queries.CreateTag(ts.ctx, CreateTagParams{
			TaskID: toNullInt(int(taskID)),
			Tag:    toNullString(tag),
		})
		if err != nil {
			return err
		}
	}

	return nil
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

	queries := New(ts.db)

	_, err := queries.UpdateTaskField(ts.ctx, params)

	if err != nil {
		return err
	}

	return nil
}

func (ts *DBLocal) AddSession(session models.Session) error {
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

func (ts *DBLocal) GetSessionsByTaskID(taskID int) ([]models.Session, error) {
	queries := New(ts.db)

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
	queries := New(ts.db)

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

func (ts *DBLocal) InsertTags(taskID int64, tags []string) error {
	queries := New(ts.db)

	for _, tag := range tags {
		_, err := queries.CreateTag(ts.ctx, CreateTagParams{
			TaskID: toNullInt(int(taskID)),
			Tag:    toNullString(tag),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (ts *DBLocal) GetTagsByTaskID(taskID int64) ([]string, error) {
	queries := New(ts.db)

	dbTags, err := queries.GetTagsByTaskID(ts.ctx, toNullInt(int(taskID)))
	if err != nil {
		return nil, err
	}

	tags := []string{}
	for _, dbTag := range dbTags {
		tags = append(tags, dbTag.Tag.String)
	}

	return tags, nil
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

func dbTaskToTask(t Task) models.Task {
	return models.Task{
		ID:          t.ID,
		Title:       t.Title.String,
		Description: t.Description.String,
		Category:    models.Category(t.Category.String),
		Priority:    models.Priority(t.Priority.String),
		Status:      models.Status(t.Status.String),
		OpenedAt:    t.OpenedAt.Time,
		ClosedAt:    t.ClosedAt.Time,
		CreatedAt:   t.CreatedAt.Time,
		UpdatedAt:   t.UpdatedAt.Time,
	}
}

func dbTasksToTasks(t []Task) []models.Task {
	var tasksList []models.Task
	for _, task := range t {
		tasksList = append(tasksList, dbTaskToTask(task))
	}
	return tasksList
}
