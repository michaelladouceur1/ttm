package handlers

import (
	"fmt"
	"strconv"
	"time"
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func StartHandler(cmd *cobra.Command, args []string, store *store.Store) {
	if fs.SessionFileExists() {
		fmt.Println("Session already started. Please end the current session first.")
		return
	}

	taskListIdString := args[0]
	taskListId, err := strconv.Atoi(taskListIdString)
	if err != nil {
		fmt.Println(err)
		return
	}

	taskId, err := fs.GetTaskIDFromListID(int64(taskListId))
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = fs.CreateSessionFile(taskId, time.Now())
	if err != nil {
		fmt.Println(err)
		return
	}

	task, err := store.GetTaskByID(taskId)
	if err != nil {
		fmt.Println(err)
		return
	}

	logger.LogSessionStart(task)
}
