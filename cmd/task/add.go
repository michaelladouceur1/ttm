package task

import (
	"fmt"
	"time"
	"ttm/pkg/models"
	"ttm/pkg/render"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task",
	Args:  cobra.MinimumNArgs(1),
	Run:   addHandler,
}

var addCategoryFlag = &ttmConfig.AddFlags.Category
var addPriorityFlag = &ttmConfig.AddFlags.Priority
var addStatusFlag = &ttmConfig.AddFlags.Status

func init() {

	addCmd.Flags().StringVarP(addCategoryFlag, "category", "c", *addCategoryFlag, "Default category")
	addCmd.Flags().StringVarP(addPriorityFlag, "priority", "p", *addPriorityFlag, "Default priority")
	addCmd.Flags().StringVarP(addStatusFlag, "status", "s", *addStatusFlag, "Default status")
}

func addHandler(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Please provide a title for the task")
		return
	}

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

	err = taskStore.InsertTask(task)
	if err != nil {
		fmt.Println("Error adding task: ", err)
		return
	}

	render.RenderAddTask(task)
}
