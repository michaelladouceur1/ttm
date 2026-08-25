package store

import (
	_ "embed"
	"time"
	"ttm/pkg/models"

	_ "github.com/mattn/go-sqlite3"
)

type StoreStrategy interface {
	Init() error
	// Tasks
	InsertTask(task models.Task) (models.Task, error)
	GetTaskByID(taskID int64) (models.Task, error)
	ListTasks(statuses []models.Status, priorities []models.Priority) ([]models.Task, error)
	SearchTasks(search models.TaskSearch) ([]models.Task, error)
	UpdateTitle(taskID int64, title string) error
	UpdateDescription(taskID int64, description string) error
	UpdatePriority(taskID int64, priority models.Priority) error
	UpdateStatus(taskID int64, status models.Status) error
	UpdateTags(taskID int64, tags []string) error
	UpdateOpenedAt(taskID int64, openedAt time.Time) error
	UpdateClosedAt(taskID int64, closedAt time.Time) error
	// Sessions
	AddSession(session models.Session) error
	GetSessionsByTaskID(taskID int) ([]models.Session, error)
	GetSessionsByTimeRange(startTime time.Time, endTime time.Time) ([]models.Session, error)
	// Tags
	InsertTags(taskID int64, tags []string) error
	GetTagsByTaskID(taskID int64) ([]string, error)
	ListTagCounts() ([]models.TagCount, error)
	// Notes
	InsertNote(note models.Note) (models.Note, error)
	GetNotesByTaskID(taskID int64) ([]models.Note, error)
}

type Store struct {
	strategy StoreStrategy
}

func NewStore(strategy StoreStrategy) *Store {
	return &Store{
		strategy: strategy,
	}
}

func Init(strategy StoreStrategy) error {
	store := NewStore(strategy)
	return store.strategy.Init()
}

func (s *Store) UpdateStoreStrategy(strategy StoreStrategy) {
	s.strategy = strategy
}

func (s *Store) Init() error {
	return s.strategy.Init()
}

func (s *Store) InsertTask(task models.Task) error {
	newTask, err := s.strategy.InsertTask(task)
	if err != nil {
		return err
	}

	if len(task.Tags) > 0 {
		if err := s.strategy.InsertTags(newTask.ID, task.Tags); err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) GetTaskByID(taskID int64) (models.Task, error) {
	task, err := s.strategy.GetTaskByID(taskID)
	if err != nil {
		return models.Task{}, err
	}

	sessions, err := s.GetSessionsByTaskID(int(task.ID))
	if err != nil {
		return models.Task{}, err
	}

	task.Sessions = sessions
	task.CalculateDuration()

	tags, err := s.strategy.GetTagsByTaskID(task.ID)
	if err != nil {
		return models.Task{}, err
	}

	task.Tags = tags

	return task, nil
}

func (s *Store) ListTasks(statuses []models.Status, priorities []models.Priority) ([]models.Task, error) {
	tasks, err := s.strategy.ListTasks(statuses, priorities)
	if err != nil {
		return nil, err
	}

	return s.populateTasks(tasks)
}

func (s *Store) SearchTasks(search models.TaskSearch) ([]models.Task, error) {
	tasks, err := s.strategy.SearchTasks(search)
	if err != nil {
		return nil, err
	}

	return s.populateTasks(tasks)
}

func (s *Store) populateTasks(tasks []models.Task) ([]models.Task, error) {
	// TODO: Refactor to run in parallel
	for i, task := range tasks {
		sessions, err := s.GetSessionsByTaskID(int(task.ID))
		if err != nil {
			return nil, err
		}

		tasks[i].Sessions = sessions
		tasks[i].CalculateDuration()

		tags, err := s.strategy.GetTagsByTaskID(task.ID)
		if err != nil {
			return nil, err
		}

		tasks[i].Tags = tags
	}

	models.SortTasksByID(tasks)
	models.PopulateListIDs(tasks)

	return tasks, nil
}

func (s *Store) UpdateTitle(taskID int64, title string) error {
	return s.strategy.UpdateTitle(taskID, title)
}

func (s *Store) UpdateDescription(taskID int64, description string) error {
	return s.strategy.UpdateDescription(taskID, description)
}

func (s *Store) UpdatePriority(taskID int64, priority models.Priority) error {
	return s.strategy.UpdatePriority(taskID, priority)
}

func (s *Store) UpdateStatus(taskID int64, status models.Status) error {
	return s.strategy.UpdateStatus(taskID, status)
}

func (s *Store) UpdateTags(taskID int64, tags []string) error {
	return s.strategy.UpdateTags(taskID, tags)
}

func (s *Store) UpdateOpenedAt(taskID int64, openedAt time.Time) error {
	return s.strategy.UpdateOpenedAt(taskID, openedAt)
}

func (s *Store) UpdateClosedAt(taskID int64, closedAt time.Time) error {
	return s.strategy.UpdateClosedAt(taskID, closedAt)
}

func (s *Store) InsertTags(taskID int64, tags []string) error {
	return s.strategy.InsertTags(taskID, tags)
}

func (s *Store) ListTagCounts() ([]models.TagCount, error) {
	return s.strategy.ListTagCounts()
}

func (s *Store) AddSession(session models.Session) error {
	return s.strategy.AddSession(session)
}

func (s *Store) GetSessionsByTaskID(taskID int) ([]models.Session, error) {
	return s.strategy.GetSessionsByTaskID(taskID)
}

func (s *Store) GetSessionsByTimeRange(startTime time.Time, endTime time.Time) ([]models.Session, error) {
	return s.strategy.GetSessionsByTimeRange(startTime, endTime)
}

func (s *Store) InsertNote(taskID int64, content string) (models.Note, error) {
	return s.strategy.InsertNote(models.Note{
		ID:        0,
		TaskID:    taskID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

func (s *Store) GetNotesByTaskID(taskID int64) ([]models.Note, error) {
	return s.strategy.GetNotesByTaskID(taskID)
}
