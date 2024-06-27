/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"
	"ttm/pkg/models"
	"ttm/pkg/render"

	"github.com/spf13/cobra"
)

// listTaskCmd represents the list command
var listTaskCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run:   listHandler,
}

var listCategoryFlag = &ttmConfig.ListFlags.Category
var listPriorityFlag = &ttmConfig.ListFlags.Priority
var listStatusFlag = &ttmConfig.ListFlags.Status

func init() {
	rootCmd.AddCommand(listTaskCmd)

	listTaskCmd.Flags().StringVarP(listCategoryFlag, "category", "c", *listCategoryFlag, "Filter tasks by category")
	listTaskCmd.Flags().StringVarP(listPriorityFlag, "priority", "p", *listPriorityFlag, "Filter tasks by priority")
	listTaskCmd.Flags().StringVarP(listStatusFlag, "status", "s", *listStatusFlag, "Filter tasks by status")
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

	tasks = getTasksDuration(tasks)

	render.RenderTasks(tasks)
}

func getTasksDuration(tasks []models.Task) []models.Task {
	taskChannel := make(chan models.Task)
	for _, task := range tasks {
		go getTaskDuration(task, taskChannel)
	}

	var tasksWithDuration []models.Task
	for range tasks {
		taskWithDuration := <-taskChannel
		tasksWithDuration = append(tasksWithDuration, taskWithDuration)
	}

	return tasksWithDuration
}

func getTaskDuration(task models.Task, taskChannel chan models.Task) error {
	sessions, err := taskStore.GetSessionByTaskID(int(task.ID))
	if err != nil {
		return err
	}

	var totalDuration time.Time
	for _, session := range sessions {
		sessionDuration := session.EndTime.Sub(session.StartTime)
		totalDuration = totalDuration.Add(sessionDuration)
	}

	task.Duration = totalDuration

	taskChannel <- task

	return nil
}
