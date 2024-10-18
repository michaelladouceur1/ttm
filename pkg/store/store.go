package store

import (
	_ "embed"
	"time"
	"ttm/pkg/models"

	_ "github.com/mattn/go-sqlite3"
)

type StoreStrategy interface {
	Init() error
	InsertTask(task models.Task) error
	GetTaskByID(taskID int64) (models.Task, error)
	ListTasks(titleDescSearch string, category models.Category, status models.Status, priority models.Priority) ([]models.Task, error)
	UpdateTitle(taskID int64, title string) error
	UpdateDescription(taskID int64, description string) error
	UpdateCategory(taskID int64, category models.Category) error
	UpdatePriority(taskID int64, priority models.Priority) error
	UpdateStatus(taskID int64, status models.Status) error
	UpdateOpenedAt(taskID int64, openedAt time.Time) error
	UpdateClosedAt(taskID int64, closedAt time.Time) error
	AddSession(session models.Session) error
	GetSessionsByTaskID(taskID int) ([]models.Session, error)
	GetSessionsByTimeRange(startTime time.Time, endTime time.Time) ([]models.Session, error)
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
	return store.Init()
}

func (s *Store) UpdateStoreStrategy(strategy StoreStrategy) {
	s.strategy = strategy
}

func (s *Store) Init() error {
	return s.strategy.Init()
}

func (s *Store) InsertTask(task models.Task) error {
	return s.strategy.InsertTask(task)
}

func (s *Store) GetTaskByID(taskID int64) (models.Task, error) {
	return s.strategy.GetTaskByID(taskID)
}

func (s *Store) ListTasks(titleDescSearch string, category models.Category, status models.Status, priority models.Priority) ([]models.Task, error) {
	tasks, err := s.strategy.ListTasks(titleDescSearch, category, status, priority)
	if err != nil {
		return nil, err
	}

	// TODO: Refactor to run in parallel
	for i, task := range tasks {
		sessions, err := s.GetSessionsByTaskID(int(task.ID))
		if err != nil {
			return nil, err
		}

		tasks[i].Sessions = sessions
		tasks[i].CalculateDuration()
	}

	models.SortTasksById(&tasks)
	models.PopulateListIDs(&tasks)

	return tasks, nil
}

func (s *Store) UpdateTitle(taskID int64, title string) error {
	return s.strategy.UpdateTitle(taskID, title)
}

func (s *Store) UpdateDescription(taskID int64, description string) error {
	return s.strategy.UpdateDescription(taskID, description)
}

func (s *Store) UpdateCategory(taskID int64, category models.Category) error {
	return s.strategy.UpdateCategory(taskID, category)
}

func (s *Store) UpdatePriority(taskID int64, priority models.Priority) error {
	return s.strategy.UpdatePriority(taskID, priority)
}

func (s *Store) UpdateStatus(taskID int64, status models.Status) error {
	return s.strategy.UpdateStatus(taskID, status)
}

func (s *Store) UpdateOpenedAt(taskID int64, openedAt time.Time) error {
	return s.strategy.UpdateOpenedAt(taskID, openedAt)
}

func (s *Store) UpdateClosedAt(taskID int64, closedAt time.Time) error {
	return s.strategy.UpdateClosedAt(taskID, closedAt)
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
