package handlers

import (
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/models"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func SearchHandler(_ *cobra.Command, args []string, store *store.Store) {
	search, err := models.ParseTaskSearch(args[0])
	if err != nil {
		logger.LogError("Error parsing search query: ", err)
		return
	}

	tasks, err := store.SearchTasks(search)
	if err != nil {
		logger.LogError("Error searching tasks: ", err)
		return
	}

	if err := fs.UpdateIDMapFile(tasks); err != nil {
		logger.LogError("Error searching tasks: ", err)
		return
	}

	logger.LogTasks(tasks)
}
