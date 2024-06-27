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
	ListTasks(titleDescSearch string, category models.Category, status models.Status, priority models.Priority) ([]models.Task, error)
	UpdateTitle(taskID int, title string) error
	UpdateDescription(taskID int, description string) error
	UpdateCategory(taskID int, category models.Category) error
	UpdatePriority(taskID int, priority models.Priority) error
	UpdateStatus(taskID int, status models.Status) error
	UpdateOpenedAt(taskID int, openedAt time.Time) error
	UpdateClosedAt(taskID int, closedAt time.Time) error
	AddSession(session models.Session) error
	GetSessionByTaskID(taskID int) ([]models.Session, error)
}

type Store struct {
	strategy StoreStrategy
}

func NewStore(strategy StoreStrategy) *Store {
	return &Store{
		strategy: strategy,
	}
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

func (s *Store) ListTasks(titleDescSearch string, category models.Category, status models.Status, priority models.Priority) ([]models.Task, error) {
	return s.strategy.ListTasks(titleDescSearch, category, status, priority)
}

func (s *Store) UpdateTitle(taskID int, title string) error {
	return s.strategy.UpdateTitle(taskID, title)
}

func (s *Store) UpdateDescription(taskID int, description string) error {
	return s.strategy.UpdateDescription(taskID, description)
}

func (s *Store) UpdateCategory(taskID int, category models.Category) error {
	return s.strategy.UpdateCategory(taskID, category)
}

func (s *Store) UpdatePriority(taskID int, priority models.Priority) error {
	return s.strategy.UpdatePriority(taskID, priority)
}

func (s *Store) UpdateStatus(taskID int, status models.Status) error {
	return s.strategy.UpdateStatus(taskID, status)
}

func (s *Store) UpdateOpenedAt(taskID int, openedAt time.Time) error {
	return s.strategy.UpdateOpenedAt(taskID, openedAt)
}

func (s *Store) UpdateClosedAt(taskID int, closedAt time.Time) error {
	return s.strategy.UpdateClosedAt(taskID, closedAt)
}

func (s *Store) AddSession(session models.Session) error {
	return s.strategy.AddSession(session)
}

func (s *Store) GetSessionByTaskID(taskID int) ([]models.Session, error) {
	return s.strategy.GetSessionByTaskID(taskID)
}
