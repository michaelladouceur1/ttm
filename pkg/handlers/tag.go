package handlers

import (
	"strconv"
	"strings"
	"ttm/pkg/logger"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func TagHandler(cmd *cobra.Command, args []string, store *store.Store) {
	if len(args) < 2 {
		cmd.PrintErrln("Error: Task ID and at least one tag are required")
		return
	}

	taskIDArg, err := strconv.Atoi(args[0])
	if err != nil {
		cmd.PrintErrln("Error: Invalid task ID")
		return
	}
	tagsArg := args[1]

	taskID := int64(taskIDArg)
	tags := strings.Split(tagsArg, ",")

	if err := store.InsertTags(taskID, tags); err != nil {
		cmd.PrintErrln("Error:", err)
	}

	updatedTask, err := store.GetTaskByID(taskID)
	if err != nil {
		cmd.PrintErrln("Error:", err)
		return
	}

	logger.LogUpdateTask(updatedTask)
}
