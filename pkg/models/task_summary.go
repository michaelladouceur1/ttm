package models

import "time"

type TaskSummary struct {
	TasksByDay []TasksByDay
}

type TasksByDay struct {
	Day   time.Time
	Tasks []Task
}
