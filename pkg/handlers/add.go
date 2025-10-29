package handlers

import (
	"time"
	"ttm/pkg/logger"
	"ttm/pkg/models"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func AddHandler(cmd *cobra.Command, args []string, store *store.Store) {
	addCategoryFlag := cmd.Flags().Lookup("category").Value.String()
	addPriorityFlag := cmd.Flags().Lookup("priority").Value.String()
	addStatusFlag := cmd.Flags().Lookup("status").Value.String()

	var title, description string
	title = args[0]

	if len(args) > 1 {
		description = args[1]
	}

	task := models.Task{
		Title:       title,
		Description: description,
		Category:    models.Category(addCategoryFlag),
		Priority:    models.Priority(addPriorityFlag),
		Status:      models.Status(addStatusFlag),
		OpenedAt:    time.Now(),
	}

	err := task.Validate()
	if err != nil {
		logger.LogError("Error adding task: ", err)
		return
	}

	err = store.InsertTask(task)
	if err != nil {
		logger.LogError("Error adding task: ", err)
		return
	}

	logger.LogAddTask(task)
}
