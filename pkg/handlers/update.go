package handlers

import (
	"strconv"
	"time"
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/models"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func UpdateHandler(cmd *cobra.Command, args []string, store *store.Store) {
	idArg, err := strconv.Atoi(args[0])
	if err != nil {
		logger.LogError("Error parsing task ID: ", err)
		return
	}

	titleFlag, _ := cmd.Flags().GetString("title")
	descriptionFlag, _ := cmd.Flags().GetString("description")
	categoryFlag, _ := cmd.Flags().GetString("category")
	priorityFlag, _ := cmd.Flags().GetString("priority")
	statusFlag, _ := cmd.Flags().GetString("status")
	openedAtFlag, _ := cmd.Flags().GetString("openedAt")
	closedAtFlag, _ := cmd.Flags().GetString("closedAt")

	id, err := fs.GetTaskIDFromListID(int64(idArg))
	if err != nil {
		logger.LogError("Error getting task ID: ", err)
		return
	}

	category := models.Category(categoryFlag)
	priority := models.Priority(priorityFlag)
	status := models.Status(statusFlag)

	err = category.Validate()
	if err != nil {
		logger.LogError(err)
		return
	}

	err = priority.Validate()
	if err != nil {
		logger.LogError(err)
		return
	}

	err = status.Validate()
	if err != nil {
		logger.LogError(err)
		return
	}

	if titleFlag == "" && descriptionFlag == "" && categoryFlag == "" && priorityFlag == "" && statusFlag == "" && openedAtFlag == "" && closedAtFlag == "" {
		logger.LogError("Please provide at least one field to update")
		return
	}

	if titleFlag != "" {
		err = store.UpdateTitle(id, titleFlag)
		if err != nil {
			logger.LogError("Error updating title: ", err)
			return
		}
	}

	if descriptionFlag != "" {
		err = store.UpdateDescription(id, descriptionFlag)
		if err != nil {
			logger.LogError("Error updating description: ", err)
			return
		}
	}

	if categoryFlag != "" {
		err = store.UpdateCategory(id, category)
		if err != nil {
			logger.LogError("Error updating category: ", err)
			return
		}
	}

	if priorityFlag != "" {
		err = store.UpdatePriority(id, priority)
		if err != nil {
			logger.LogError("Error updating priority: ", err)
			return
		}
	}

	if statusFlag != "" {
		err = store.UpdateStatus(id, status)
		if err != nil {
			logger.LogError("Error updating status: ", err)
			return
		}
	}

	if openedAtFlag != "" {
		openedTime, err := time.Parse(time.RFC3339, openedAtFlag)
		if err != nil {
			logger.LogError("Error parsing start time: ", err)
			return
		}
		err = store.UpdateOpenedAt(id, openedTime)
		if err != nil {
			logger.LogError("Error updating start time: ", err)
			return
		}
	}

	if closedAtFlag != "" {
		closedTime, err := time.Parse(time.RFC3339, closedAtFlag)
		if err != nil {
			logger.LogError("Error parsing end time: ", err)
			return
		}
		err = store.UpdateClosedAt(id, closedTime)
		if err != nil {
			logger.LogError("Error updating end time: ", err)
			return
		}
	}

	updatedTask, err := store.GetTaskByID(id)
	if err != nil {
		logger.LogError("Error retrieving updated task: ", err)
		return
	}

	logger.LogUpdateTask(updatedTask)
}
