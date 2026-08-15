package ui

import (
	"fmt"
	"strconv"
	"time"
	"ttm/pkg/fs"
	"ttm/pkg/store"
)

func startSession(st *store.Store, taskListID string, startedAt time.Time) string {
	if fs.SessionFileExists() {
		return "Session already started. Please end the current session first."
	}

	listID, err := strconv.ParseInt(taskListID, 10, 64)
	if err != nil {
		return "Error starting session: " + err.Error()
	}

	taskID, err := fs.GetTaskIDFromTempID(listID)
	if err != nil {
		return "Error starting session: " + err.Error()
	}

	task, err := st.GetTaskByID(taskID)
	if err != nil {
		return "Error starting session: " + err.Error()
	}

	if err := fs.CreateSessionFile(taskID, startedAt); err != nil {
		return "Error starting session: " + err.Error()
	}

	return fmt.Sprintf(
		"Session started.\n\nTask: %s\nStarted: %s",
		task.Title,
		startedAt.Round(time.Second).Format("2006-01-02 15:04:05"),
	)
}
