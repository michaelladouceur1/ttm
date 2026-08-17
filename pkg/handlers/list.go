package handlers

import (
	"ttm/pkg/config"
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/models"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func ListHandler(cmd *cobra.Command, args []string, cfg *config.Config, store *store.Store) {
	listCategoryFlag := &cfg.ListFlags.Category
	listPriorityFlag := &cfg.ListFlags.Priority
	listStatusFlag := &cfg.ListFlags.Status

	categories := []models.Category{}
	for _, cat := range *listCategoryFlag {
		category := models.Category(cat)
		if err := category.Validate(); err != nil {
			logger.LogError("Error listing tasks: ", err)
			return
		}
		categories = append(categories, category)
	}

	priorities := []models.Priority{}
	for _, prio := range *listPriorityFlag {
		priority := models.Priority(prio)
		if err := priority.Validate(); err != nil {
			logger.LogError("Error listing tasks: ", err)
			return
		}
		priorities = append(priorities, priority)
	}

	statuses := []models.Status{}
	for _, stat := range *listStatusFlag {
		println("status: ", stat)
		status := models.Status(stat)
		if err := status.Validate(); err != nil {
			logger.LogError("Error listing tasks: ", err)
			return
		}
		statuses = append(statuses, status)
	}

	tasks, err := store.ListTasks(categories, statuses, priorities)
	if err != nil {
		logger.LogError("Error listing tasks: ", err)
		return
	}

	err = fs.UpdateIDMapFile(tasks)
	if err != nil {
		logger.LogError("Error listing tasks: ", err)
		return
	}

	logger.LogTasks(tasks)
}
