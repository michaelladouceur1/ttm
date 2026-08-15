package ui

import (
	"fmt"
	"time"
	"ttm/pkg/fs"
	"ttm/pkg/models"
	"ttm/pkg/store"
)

func endSession(st *store.Store, endedAt time.Time) string {
	if !fs.SessionFileExists() {
		return "No session found. Please start a session first."
	}

	session, err := fs.ReadSessionFile()
	if err != nil {
		return "Error ending session: " + err.Error()
	}

	if err := st.AddSession(models.Session{
		TaskId:    session.TaskID,
		StartTime: session.StartTime,
		EndTime:   endedAt,
	}); err != nil {
		return "Error ending session: " + err.Error()
	}

	if err := fs.RemoveSessionFile(); err != nil {
		return "Error ending session: " + err.Error()
	}

	task, err := st.GetTaskByID(session.TaskID)
	if err != nil {
		return "Error ending session: " + err.Error()
	}

	return fmt.Sprintf(
		"Session ended.\n\nTask: %s\nStarted: %s\nDuration: %s",
		task.Title,
		session.StartTime.Round(time.Second).Format("2006-01-02 15:04:05"),
		endedAt.Sub(session.StartTime).Round(time.Second),
	)
}
