package handlers

import (
	"fmt"
	"time"
	"ttm/pkg/logger"
	"ttm/pkg/models"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func SummaryHandler(cmd *cobra.Command, args []string, store *store.Store) {
	days, _ := cmd.Flags().GetInt("days")

	taskSummary := models.TaskSummary{}
	for day := days - 1; day >= 0; day-- {
		currentDay := time.Now().AddDate(0, 0, -day)
		beginningOfDay := time.Date(currentDay.Year(), currentDay.Month(), currentDay.Day(), 0, 0, 0, 0, currentDay.Location())
		endOfDay := time.Date(currentDay.Year(), currentDay.Month(), currentDay.Day(), 23, 59, 59, 0, currentDay.Location())

		sessions, err := store.GetSessionsByTimeRange(beginningOfDay, endOfDay)
		if err != nil {
			fmt.Println(err)
			return
		}

		var tasks []models.Task
		for _, session := range sessions {
			task, err := store.GetTaskByID(session.TaskId)
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

	logger.LogTaskSummary(taskSummary)
}

func taskExistsInList(task models.Task, tasks []models.Task) bool {
	for _, t := range tasks {
		if t.ID == task.ID {
			return true
		}
	}
	return false
}
