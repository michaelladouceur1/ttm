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

type Store struct {
	ctx context.Context
	db  *sql.DB
}

//go:embed schema.postgres.sql
var ddl string

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Init() error {
	var err error

	if os.MkdirAll(paths.GetTTMDirectory(), os.ModePerm); err != nil {
		return err
	}

	s.ctx = context.Background()

	s.db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/ttmdb")
	if err != nil {
		return err
	}

	if _, err := s.db.ExecContext(s.ctx, ddl); err != nil {
		return err
	}

	return nil
}

func (s *Store) InsertTask(task models.Task) (models.Task, error) {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	queries := New(s.db)

	newTask, err := queries.CreateTask(s.ctx, CreateTaskParams{
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

func (s *Store) GetTaskByID(taskID int64) (models.Task, error) {
	queries := New(s.db)

	dbTask, err := queries.GetTaskById(s.ctx, taskID)

	if err != nil {
		return models.Task{}, err
	}

	task := dbTasksToTasks([]Task{dbTask})

	return task[0], nil

}

func (s *Store) ListTasks(titleDescSearch string, category models.Category, status models.Status, priority models.Priority) ([]models.Task, error) {
	queries := New(s.db)

	dbTasks, err := queries.ListTasks(s.ctx, ListTasksParams{
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

func (s *Store) SearchTasks(search models.TaskSearch) ([]models.Task, error) {
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

	rows, err := s.db.QueryContext(s.ctx, query, args...)
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

func (s *Store) UpdateTitle(taskID int64, title string) error {
	return s.updateTaskField(UpdateTaskFieldParams{
		ID:    taskID,
		Title: toNullString(title),
	})
}

func (s *Store) UpdateDescription(taskID int64, description string) error {
	return s.updateTaskField(UpdateTaskFieldParams{
		ID:          taskID,
		Description: toNullString(description),
	})
}

func (s *Store) UpdateCategory(taskID int64, category models.Category) error {
	return s.updateTaskField(UpdateTaskFieldParams{
		ID:       taskID,
		Category: toNullString(string(category)),
	})
}

func (s *Store) UpdatePriority(taskID int64, priority models.Priority) error {
	return s.updateTaskField(UpdateTaskFieldParams{
		ID:       taskID,
		Priority: toNullString(string(priority)),
	})
}

func (s *Store) UpdateStatus(taskID int64, status models.Status) error {
	return s.updateTaskField(UpdateTaskFieldParams{
		ID:     taskID,
		Status: toNullString(string(status)),
	})
}

func (s *Store) UpdateTags(taskID int64, tags []string) error {
	queries := New(s.db)

	err := queries.DeleteTagsByTaskID(s.ctx, toNullInt(int(taskID)))
	if err != nil {
		return err
	}

	for _, tag := range tags {
		_, err := queries.CreateTag(s.ctx, CreateTagParams{
			TaskID: toNullInt(int(taskID)),
			Tag:    toNullString(tag),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) UpdateOpenedAt(taskID int64, openedAt time.Time) error {
	return s.updateTaskField(UpdateTaskFieldParams{
		ID:       taskID,
		OpenedAt: toNullTime(openedAt),
	})
}

func (s *Store) UpdateClosedAt(taskID int64, closedAt time.Time) error {
	return s.updateTaskField(UpdateTaskFieldParams{
		ID:       taskID,
		ClosedAt: toNullTime(closedAt),
	})
}

func (s *Store) updateTaskField(params UpdateTaskFieldParams) error {
	params.UpdatedAt = toNullTime(time.Now())

	queries := New(s.db)

	_, err := queries.UpdateTaskField(s.ctx, params)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) AddSession(session models.Session) error {
	queries := New(s.db)

	_, err := queries.CreateSession(s.ctx, CreateSessionParams{
		TaskID:    toNullInt(int(session.TaskId)),
		StartTime: toNullTime(session.StartTime),
		EndTime:   toNullTime(session.EndTime),
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetSessionsByTaskID(taskID int) ([]models.Session, error) {
	queries := New(s.db)

	dbSessions, err := queries.GetSessionsByTaskID(s.ctx, toNullInt(taskID))

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

func (s *Store) GetSessionsByTimeRange(startTime time.Time, endTime time.Time) ([]models.Session, error) {
	queries := New(s.db)

	dbSessions, err := queries.GetSessionsByTimeRange(s.ctx, GetSessionsByTimeRangeParams{
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

func (s *Store) InsertTags(taskID int64, tags []string) error {
	queries := New(s.db)

	for _, tag := range tags {
		_, err := queries.CreateTag(s.ctx, CreateTagParams{
			TaskID: toNullInt(int(taskID)),
			Tag:    toNullString(tag),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) GetTagsByTaskID(taskID int64) ([]string, error) {
	queries := New(s.db)

	dbTags, err := queries.GetTagsByTaskID(s.ctx, toNullInt(int(taskID)))
	if err != nil {
		return nil, err
	}

	tags := []string{}
	for _, dbTag := range dbTags {
		tags = append(tags, dbTag.Tag.String)
	}

	return tags, nil
}

func (s *Store) ListTagCounts() ([]models.TagCount, error) {
	rows, err := s.db.QueryContext(s.ctx, `
		SELECT tag, COUNT(DISTINCT task_id)
		FROM tags
		WHERE tag IS NOT NULL AND tag <> ''
		GROUP BY tag
		ORDER BY LOWER(tag), tag`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.TagCount
	for rows.Next() {
		var tag models.TagCount
		if err := rows.Scan(&tag.Tag, &tag.Count); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tags, nil
}

func (s *Store) InsertNote(note models.Note) (models.Note, error) {
	queries := New(s.db)

	newNote, err := queries.CreateNote(s.ctx, CreateNoteParams{
		TaskID:    toNullInt(int(note.TaskID)),
		Content:   toNullString(note.Content),
		CreatedAt: toNullTime(note.CreatedAt),
		UpdatedAt: toNullTime(note.UpdatedAt),
	})

	if err != nil {
		return models.Note{}, err
	}

	return models.Note{
		ID:        newNote.ID,
		TaskID:    newNote.TaskID.Int64,
		Content:   newNote.Content.String,
		CreatedAt: newNote.CreatedAt.Time,
		UpdatedAt: newNote.UpdatedAt.Time,
	}, nil
}

func (s *Store) GetNotesByTaskID(taskID int64) ([]models.Note, error) {
	queries := New(s.db)

	dbNotes, err := queries.GetNotesByTaskID(s.ctx, toNullInt(int(taskID)))
	if err != nil {
		return nil, err
	}

	notes := []models.Note{}
	for _, dbNote := range dbNotes {
		notes = append(notes, models.Note{
			ID:        dbNote.ID,
			TaskID:    dbNote.TaskID.Int64,
			Content:   dbNote.Content.String,
			CreatedAt: dbNote.CreatedAt.Time,
			UpdatedAt: dbNote.UpdatedAt.Time,
		})
	}

	return notes, nil
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
