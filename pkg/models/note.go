package models

import "time"

type Note struct {
	ID        int64     `json:"id"`
	TaskID    int64     `json:"task_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
