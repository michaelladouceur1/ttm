package handlers

import (
	"fmt"
	"strconv"
	"time"
	"ttm/pkg/fs"
	"ttm/pkg/logger"

	"github.com/spf13/cobra"
)

func StartHandler(cmd *cobra.Command, args []string) {
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

	logger.LogSessionStart(taskId)
}
