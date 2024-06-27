package store

import (
	_ "embed"
	"sort"
	"time"
	"ttm/pkg/models"

	_ "github.com/mattn/go-sqlite3"
)

type StoreStrategy interface {
	Init() error
	InsertTask(task models.Task) error
	ListTasks(titleDescSearch string, category models.Category, status models.Status, priority models.Priority) ([]models.Task, error)
	UpdateTitle(taskID int64, title string) error
	UpdateDescription(taskID int64, description string) error
	UpdateCategory(taskID int64, category models.Category) error
	UpdatePriority(taskID int64, priority models.Priority) error
	UpdateStatus(taskID int64, status models.Status) error
	UpdateOpenedAt(taskID int64, openedAt time.Time) error
	UpdateClosedAt(taskID int64, closedAt time.Time) error
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
	tasks, err := s.strategy.ListTasks(titleDescSearch, category, status, priority)
	if err != nil {
		return nil, err
	}

	err = s.getTasksDuration(&tasks)
	if err != nil {
		return nil, err
	}

	s.sortTasksByID(&tasks)
	s.populateListIDs(&tasks)

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

func (s *Store) GetSessionByTaskID(taskID int) ([]models.Session, error) {
	return s.strategy.GetSessionByTaskID(taskID)
}

// TODO: add result type to handle errors
func (s *Store) getTasksDuration(tasks *[]models.Task) error {
	taskChannel := make(chan models.Task)
	for _, task := range *tasks {
		go s.getTaskDuration(task, taskChannel)
	}

	var tasksWithDuration []models.Task
	for range *tasks {
		taskWithDuration := <-taskChannel
		tasksWithDuration = append(tasksWithDuration, taskWithDuration)
	}

	*tasks = tasksWithDuration

	return nil
}

func (s *Store) getTaskDuration(task models.Task, taskChannel chan models.Task) {
	sessions, err := s.GetSessionByTaskID(int(task.ID))
	if err != nil {
		taskChannel <- task
		return
	}

	var totalDuration time.Time
	for _, session := range sessions {
		sessionDuration := session.EndTime.Sub(session.StartTime)
		totalDuration = totalDuration.Add(sessionDuration)
	}

	task.Duration = totalDuration
	taskChannel <- task
}

func (s *Store) sortTasksByID(tasks *[]models.Task) {
	sort.Slice(*tasks, func(i, j int) bool {
		return (*tasks)[i].ID < (*tasks)[j].ID
	})
}

func (s *Store) populateListIDs(tasks *[]models.Task) {
	for i := range *tasks {
		(*tasks)[i].ListID = int64(i + 1)
	}
}
