/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package task

import (
	"fmt"
	"time"
	"ttm/cmd"
	"ttm/pkg/config"
	"ttm/pkg/fs"
	"ttm/pkg/models"
	"ttm/pkg/render"
	"ttm/pkg/store"
	"ttm/pkg/store/db"

	"github.com/spf13/cobra"
)

var ttmConfig = config.NewConfig()
var taskStore = store.NewStore(db.NewDBStore())

// taskCmd represents the task command
var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks",
	Args:  cobra.MinimumNArgs(1),
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task",
	Args:  cobra.MinimumNArgs(1),
	Run:   addHandler,
}

var listTaskCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run:   listHandler,
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a task",
	Run:   updateHandler,
}

var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "Export tasks to CSV",
	Run:   csvHandler,
}

var addCategoryFlag = &ttmConfig.AddFlags.Category
var addPriorityFlag = &ttmConfig.AddFlags.Priority
var addStatusFlag = &ttmConfig.AddFlags.Status

var listCategoryFlag = &ttmConfig.ListFlags.Category
var listPriorityFlag = &ttmConfig.ListFlags.Priority
var listStatusFlag = &ttmConfig.ListFlags.Status

func init() {
	ttmConfig.Load()
	taskStore.Init()

	cmd.RootCmd.AddCommand(taskCmd)

	taskCmd.AddCommand(addCmd)
	taskCmd.AddCommand(listTaskCmd)
	taskCmd.AddCommand(updateCmd)
	taskCmd.AddCommand(csvCmd)

	addCmd.Flags().StringVarP(addCategoryFlag, "category", "c", *addCategoryFlag, "Default category")
	addCmd.Flags().StringVarP(addPriorityFlag, "priority", "p", *addPriorityFlag, "Default priority")
	addCmd.Flags().StringVarP(addStatusFlag, "status", "s", *addStatusFlag, "Default status")

	listTaskCmd.Flags().StringVarP(listCategoryFlag, "category", "c", *listCategoryFlag, "Filter tasks by category")
	listTaskCmd.Flags().StringVarP(listPriorityFlag, "priority", "p", *listPriorityFlag, "Filter tasks by priority")
	listTaskCmd.Flags().StringVarP(listStatusFlag, "status", "s", *listStatusFlag, "Filter tasks by status")

	updateCmd.Flags().IntP("id", "i", 0, "Task ID")
	updateCmd.Flags().StringP("title", "t", "", "Update title")
	updateCmd.Flags().StringP("description", "d", "", "Update description")
	updateCmd.Flags().StringP("category", "c", "", "Update category")
	updateCmd.Flags().StringP("priority", "p", "", "Update priority")
	updateCmd.Flags().StringP("status", "s", "", "Update status")
	updateCmd.Flags().StringP("openedAt", "a", "", "Update opened time")
	updateCmd.Flags().StringP("closedAt", "b", "", "Update closed time")

	csvCmd.Flags().StringP("category", "c", "", "Filter tasks by category")
	csvCmd.Flags().StringP("status", "s", "", "Filter tasks by status")
	csvCmd.Flags().StringP("priority", "p", "", "Filter tasks by priority")
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

func listHandler(cmd *cobra.Command, args []string) {
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

	tasks, err := taskStore.ListTasks(titleDescSearch, category, status, priority)
	if err != nil {
		fmt.Println("Error listing tasks: ", err)
		return
	}

	render.RenderTasks(tasks)
}

func updateHandler(cmd *cobra.Command, args []string) {
	idFlag, _ := cmd.Flags().GetInt("id")
	titleFlag, _ := cmd.Flags().GetString("title")
	descriptionFlag, _ := cmd.Flags().GetString("description")
	categoryFlag, _ := cmd.Flags().GetString("category")
	priorityFlag, _ := cmd.Flags().GetString("priority")
	statusFlag, _ := cmd.Flags().GetString("status")
	openedAtFlag, _ := cmd.Flags().GetString("openedAt")
	closedAtFlag, _ := cmd.Flags().GetString("closedAt")

	if idFlag == 0 {
		fmt.Println("Please provide a task ID to update")
		return
	}

	category := models.Category(categoryFlag)
	priority := models.Priority(priorityFlag)
	status := models.Status(statusFlag)

	var err error
	err = category.Validate()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = priority.Validate()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = status.Validate()
	if err != nil {
		fmt.Println(err)
		return
	}

	if titleFlag == "" && descriptionFlag == "" && categoryFlag == "" && priorityFlag == "" && statusFlag == "" && openedAtFlag == "" && closedAtFlag == "" {
		fmt.Println("Please provide at least one field to update")
		return
	}

	if titleFlag != "" {
		err = taskStore.UpdateTitle(idFlag, titleFlag)
		if err != nil {
			fmt.Println("Error updating title: ", err)
			return
		}
	}

	if descriptionFlag != "" {
		err = taskStore.UpdateDescription(idFlag, descriptionFlag)
		if err != nil {
			fmt.Println("Error updating description: ", err)
			return
		}
	}

	if categoryFlag != "" {
		err = taskStore.UpdateCategory(idFlag, category)
		if err != nil {
			fmt.Println("Error updating category: ", err)
			return
		}
	}

	if priorityFlag != "" {
		err = taskStore.UpdatePriority(idFlag, priority)
		if err != nil {
			fmt.Println("Error updating priority: ", err)
			return
		}
	}

	if statusFlag != "" {
		err = taskStore.UpdateStatus(idFlag, status)
		if err != nil {
			fmt.Println("Error updating status: ", err)
			return
		}
	}

	if openedAtFlag != "" {
		openedTime, err := time.Parse(time.RFC3339, openedAtFlag)
		if err != nil {
			fmt.Println("Error parsing start time: ", err)
			return
		}
		err = taskStore.UpdateOpenedAt(idFlag, openedTime)
		if err != nil {
			fmt.Println("Error updating start time: ", err)
			return
		}
	}

	if closedAtFlag != "" {
		closedTime, err := time.Parse(time.RFC3339, closedAtFlag)
		if err != nil {
			fmt.Println("Error parsing end time: ", err)
			return
		}
		err = taskStore.UpdateClosedAt(idFlag, closedTime)
		if err != nil {
			fmt.Println("Error updating end time: ", err)
			return
		}
	}

	render.RenderUpdateTask()
}

func csvHandler(cmd *cobra.Command, args []string) {
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

	tasks, err := taskStore.ListTasks(titleDescSearch, category, status, priority)
	if err != nil {
		fmt.Println("Error listing tasks: ", err)
		return
	}

	err = fs.TasksToCSV(tasks)
	if err != nil {
		fmt.Println("Error exporting tasks to CSV: ", err)
		return
	}
}
