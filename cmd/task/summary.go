/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package task

import (
	"fmt"
	"time"
	"ttm/pkg/models"
	"ttm/pkg/render"

	"github.com/spf13/cobra"
)

// summaryCmd represents the summary command
var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Summarize tasks",
	Run:   summaryHandler,
}

func init() {
	taskCmd.AddCommand(summaryCmd)

	summaryCmd.Flags().IntP("days", "d", 7, "Number of days to summarize")
}

func summaryHandler(cmd *cobra.Command, args []string) {
	days, _ := cmd.Flags().GetInt("days")

	taskSummary := models.TaskSummary{}
	for day := days - 1; day >= 0; day-- {
		currentDay := time.Now().AddDate(0, 0, -day)
		beginningOfDay := time.Date(currentDay.Year(), currentDay.Month(), currentDay.Day(), 0, 0, 0, 0, currentDay.Location())
		endOfDay := time.Date(currentDay.Year(), currentDay.Month(), currentDay.Day(), 23, 59, 59, 0, currentDay.Location())

		sessions, err := taskStore.GetSessionsByTimeRange(beginningOfDay, endOfDay)
		if err != nil {
			fmt.Println(err)
			return
		}

		var tasks []models.Task
		for _, session := range sessions {
			task, err := taskStore.GetTaskByID(session.TaskId)
			if err != nil {
				fmt.Println(err)
				return
			}
			if taskExistsInList(task, tasks) {
				continue
			}
			tasks = append(tasks, task)
		}

		taskSummary.AddDay(models.TaskSummaryDay{
			Day:   beginningOfDay,
			Tasks: tasks,
		})
	}

	render.RenderTaskSummary(taskSummary)
}

func taskExistsInList(task models.Task, tasks []models.Task) bool {
	for _, t := range tasks {
		if t.ID == task.ID {
			return true
		}
	}
	return false
}
