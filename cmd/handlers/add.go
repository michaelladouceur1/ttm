package handlers

import (
	"fmt"
	"time"
	"ttm/pkg/models"
	"ttm/pkg/render"

	"github.com/spf13/cobra"
)

func AddHandler(cmd *cobra.Command, args []string) {
	addCategoryFlag := &config.AddFlags.Category
	addPriorityFlag := &config.AddFlags.Priority
	addStatusFlag := &config.AddFlags.Status

	var title, description string
	title = args[0]

	if len(args) > 1 {
		description = args[1]
	}

	task := models.Task{
		Title:       title,
		Description: description,
		Category:    models.Category(*addCategoryFlag),
		Priority:    models.Priority(*addPriorityFlag),
		Status:      models.Status(*addStatusFlag),
		OpenedAt:    time.Now(),
	}

	err := task.Validate()
	if err != nil {
		fmt.Println("Error adding task: ", err)
		return
	}

	err = store.InsertTask(task)
	if err != nil {
		fmt.Println("Error adding task: ", err)
		return
	}

	render.RenderAddTask(task)
}
