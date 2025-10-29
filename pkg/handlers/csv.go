package handlers

import (
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/models"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func CSVHandler(cmd *cobra.Command, args []string, store *store.Store) {
	var titleDescSearch string
	if len(args) > 0 {
		titleDescSearch = args[0]
	}

	categoryFlag := cmd.Flags().Lookup("category").Value.String()
	statusFlag := cmd.Flags().Lookup("status").Value.String()
	priorityFlag := cmd.Flags().Lookup("priority").Value.String()

	category := models.Category(categoryFlag)
	status := models.Status(statusFlag)
	priority := models.Priority(priorityFlag)

	var err error
	err = category.Validate()
	if err != nil {
		logger.LogError("Error exporting to CSV: ", err)
		return
	}

	err = status.Validate()
	if err != nil {
		logger.LogError("Error exporting to CSV: ", err)
		return
	}

	err = priority.Validate()
	if err != nil {
		logger.LogError("Error exporting to CSV: ", err)
		return
	}

	tasks, err := store.ListTasks(titleDescSearch, category, status, priority)
	if err != nil {
		logger.LogError("Error exporting to CSV: ", err)
		return
	}

	err = fs.TasksToCSV(tasks)
	if err != nil {
		logger.LogError("Error exporting to CSV: ", err)
		return
	}
}
