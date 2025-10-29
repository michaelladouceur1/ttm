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
	var titleDescSearch string
	if len(args) > 0 {
		titleDescSearch = args[0]
	}

	category := models.Category(*listCategoryFlag)
	status := models.Status(*listStatusFlag)
	priority := models.Priority(*listPriorityFlag)

	var err error
	err = category.Validate()
	if err != nil {
		logger.LogError("Error listing tasks: ", err)
		return
	}

	err = status.Validate()
	if err != nil {
		logger.LogError("Error listing tasks: ", err)
		return
	}

	err = priority.Validate()
	if err != nil {
		logger.LogError("Error listing tasks: ", err)
		return
	}

	tasks, err := store.ListTasks(titleDescSearch, category, status, priority)
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
