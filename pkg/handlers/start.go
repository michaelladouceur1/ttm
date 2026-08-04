package handlers

import (
	"strconv"
	"time"
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func StartHandler(cmd *cobra.Command, args []string, store *store.Store) {
	if fs.SessionFileExists() {
		logger.LogError("Session already started. Please end the current session first.")
		return
	}

	taskListIdString := args[0]
	taskListId, err := strconv.Atoi(taskListIdString)
	if err != nil {
		logger.LogError("Error starting session: ", err)
		return
	}

	taskID, err := fs.GetTaskIDFromTempID(int64(taskListId))
	if err != nil {
		logger.LogError("Error starting session: ", err)
		return
	}

	task, err := store.GetTaskByID(taskID)
	if err != nil {
		logger.LogError("Error starting session: ", err)
		return
	}

	start := time.Now()
	if err := fs.CreateSessionFile(taskID, start); err != nil {
		logger.LogError("Error starting session: ", err)
		return
	}

	logger.LogSessionStart(task, start)
}
