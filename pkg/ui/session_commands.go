package ui

import (
	"fmt"
	"strconv"
	"time"
	"ttm/pkg/fs"
	"ttm/pkg/models"
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

func cancelSession() string {
	if !fs.SessionFileExists() {
		return "No session found. Please start a session first."
	}

	if err := fs.RemoveSessionFile(); err != nil {
		return "Error cancelling session: " + err.Error()
	}

	return "Session cancelled."
}
