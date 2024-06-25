package models

import (
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
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    Category  `json:"category"`
	Priority    Priority  `json:"priority"`
	Status      Status    `json:"status"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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

// func ParseTime(timeStr string) (time.Time, error) {
// 	if timeStr == "" {
// 		return time.Time{}, nil
// 	}
// 	return time.Parse(time.RFC3339, timeStr)
// }

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
