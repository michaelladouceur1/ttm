package models

import "time"

type Session struct {
	ID        int64     `json:"id"`
	TaskId    int64     `json:"task_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type SessionFile struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"start_time"`
}
