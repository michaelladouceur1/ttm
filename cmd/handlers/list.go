package handlers

import (
	"fmt"
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/models"

	"github.com/spf13/cobra"
)

func ListHandler(cmd *cobra.Command, args []string) {
	listCategoryFlag := &config.ListFlags.Category
	listPriorityFlag := &config.ListFlags.Priority
	listStatusFlag := &config.ListFlags.Status

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
		fmt.Println(err)
		return
	}

	err = status.Validate()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = priority.Validate()
	if err != nil {
		fmt.Println(err)
		return
	}

	tasks, err := store.ListTasks(titleDescSearch, category, status, priority)
	if err != nil {
		fmt.Println("Error listing tasks: ", err)
		return
	}

	err = fs.UpdateIDMapFile(tasks)
	if err != nil {
		fmt.Println("Error updating ID map file: ", err)
		return
	}

	logger.LogTasks(tasks)
}
