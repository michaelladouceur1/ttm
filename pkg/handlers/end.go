package handlers

import (
	"time"
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/models"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func EndHandler(cmd *cobra.Command, args []string, store *store.Store) {
	if !fs.SessionFileExists() {
		logger.LogError("No session found. Please start a session first.")
		return
	}

	sf, err := fs.ReadSessionFile()
	if err != nil {
		logger.LogError("Error ending session: ", err)
		return
	}

	if err := store.AddSession(models.Session{
		TaskId:    sf.TaskID,
		StartTime: sf.StartTime,
		EndTime:   time.Now(),
	}); err != nil {
		logger.LogError("Error ending session: ", err)
		return
	}

	if err := fs.RemoveSessionFile(); err != nil {
		logger.LogError("Error ending session: ", err)
		return
	}

	task, err := store.GetTaskByID(sf.TaskID)
	if err != nil {
		logger.LogError("Error ending session: ", err)
		return
	}

	logger.LogSessionEnd(sf, task)
}
