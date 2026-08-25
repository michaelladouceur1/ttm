package handlers

import (
	"strings"

	"ttm/pkg/config"
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/models"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func ListHandler(cmd *cobra.Command, args []string, cfg *config.Config, store *store.Store) {
	listPriorityFlag := cmd.Flags().Lookup("priority").Value.String()
	listStatusFlag := cmd.Flags().Lookup("status").Value.String()

	priorities := []models.Priority{}
	if listPriorityFlag != "" {
		for prio := range strings.SplitSeq(listPriorityFlag, ",") {
			priority := models.Priority(prio)
			if err := priority.Validate(); err != nil {
				logger.LogError("Error listing tasks: ", err)
				return
			}
			priorities = append(priorities, priority)
		}
	}

	statuses := []models.Status{}
	if listStatusFlag != "" {
		for stat := range strings.SplitSeq(listStatusFlag, ",") {
			status := models.Status(stat)
			if err := status.Validate(); err != nil {
				logger.LogError("Error listing tasks: ", err)
				return
			}
			statuses = append(statuses, status)
		}
	}

	tasks, err := store.ListTasks(statuses, priorities)
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
