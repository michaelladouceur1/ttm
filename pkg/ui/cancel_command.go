package ui

import (
	"ttm/pkg/fs"
)

func cancelSession() string {
	if !fs.SessionFileExists() {
		return "No session found. Please start a session first."
	}

	if err := fs.RemoveSessionFile(); err != nil {
		return "Error cancelling session: " + err.Error()
	}

	return "Session cancelled."
}
