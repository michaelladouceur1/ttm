package handlers

import (
	"strconv"
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/models"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func CloseHandler(cmd *cobra.Command, args []string, store *store.Store) {
	var ids []int64
	for _, arg := range args {
		tmpID, err := strconv.Atoi(arg)
		if err != nil {
			logger.LogError("Error parsing task ID: ", err)
			return
		}

		id, err := fs.GetTaskIDFromListID(int64(tmpID))
		if err != nil {
			logger.LogError("Error getting task ID: ", err)
			return
		}

		ids = append(ids, id)
	}

	for _, id := range ids {
		err := store.UpdateStatus(id, models.StatusClosed)
		if err != nil {
			logger.LogError("Error closing task: ", err)
			return
		}
	}
}
