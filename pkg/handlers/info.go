package handlers

import (
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func InfoHandler(cmd *cobra.Command, args []string, store *store.Store) {
	if !fs.SessionFileExists() {
		logger.LogMessage("No session found. Please start a session first.")
		return
	}

	sf, err := fs.ReadSessionFile()
	if err != nil {
		logger.LogError("Error getting session info: ", err)
		return
	}

	task, err := store.GetTaskByID(sf.ID)
	if err != nil {
		logger.LogError("Error getting session info: ", err)
		return
	}

	logger.LogSessionInfo(sf, task)
}
