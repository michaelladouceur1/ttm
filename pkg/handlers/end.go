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

	sf, err := fs.RemoveSessionFile()
	if err != nil {
		logger.LogError("Error ending session: ", err)
		return
	}

	store.AddSession(models.Session{
		TaskId:    int64(sf.ID),
		StartTime: sf.StartTime,
		EndTime:   time.Now(),
	})

	task, err := store.GetTaskByID(sf.ID)
	if err != nil {
		logger.LogError("Error ending session: ", err)
		return
	}

	logger.LogSessionEnd(sf, task)
}
