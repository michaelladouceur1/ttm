package models

import (
	"sort"
	"time"
)

// TODO: These values should be configurable
type Category string
type Priority string
type Status string

const (
	CategoryTask    Category = "task"
	CategoryMeeting Category = "meeting"
)

const (
	StatusOpen   Status = "open"
	StatusClosed Status = "closed"
)

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Task struct {
	ID          int64     `json:"id"`
	ListID      int64     `json:"list_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    Category  `json:"category"`
	Priority    Priority  `json:"priority"`
	Status      Status    `json:"status"`
	Duration    time.Time `json:"duration"`
	OpenedAt    time.Time `json:"opened_at"`
	ClosedAt    time.Time `json:"closed_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Sessions    []Session `json:"sessions"`
}

func (s *Task) CalculateDuration() {
	var totalDuration time.Time
	for _, session := range s.Sessions {
		sessionDuration := session.EndTime.Sub(session.StartTime)
		totalDuration = totalDuration.Add(sessionDuration)
	}
	s.Duration = totalDuration
}

func (t *Task) Validate() error {
	var err error

	err = t.Priority.Validate()
	if err != nil {
		return err
	}

	err = t.Status.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (c Category) Validate() error {
	valid := c == CategoryTask || c == CategoryMeeting || c == ""
	if !valid {
		return &InvalidCategoryError{}
	}
	return nil
}

func (p Priority) Validate() error {
	valid := p == PriorityLow || p == PriorityMedium || p == PriorityHigh || p == ""
	if !valid {
		return &InvalidPriorityError{}
	}
	return nil
}

func (s Status) Validate() error {
	valid := s == StatusOpen || s == StatusClosed || s == ""
	if !valid {
		return &InvalidStatusError{}
	}
	return nil
}

func SortTasksById(tasks *[]Task) {
	sort.Slice(*tasks, func(i, j int) bool {
		return (*tasks)[i].ID < (*tasks)[j].ID
	})
}

func PopulateListIDs(tasks *[]Task) {
	for i := range *tasks {
		(*tasks)[i].ListID = int64(i + 1)
	}
}

type InvalidCategoryError struct{}

func (e *InvalidCategoryError) Error() string {
	return "Invalid category. Please choose from task, meeting"
}

type InvalidPriorityError struct{}

func (e *InvalidPriorityError) Error() string {
	return "Invalid priority. Please choose from low, medium, high"
}

type InvalidStatusError struct{}

func (e *InvalidStatusError) Error() string {
	return "Invalid status. Please choose from open, closed"
}
