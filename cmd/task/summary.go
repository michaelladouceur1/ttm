/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package task

import (
	"fmt"
	"time"

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
}

func summaryHandler(cmd *cobra.Command, args []string) {
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -7)

	sessions, err := taskStore.GetSessionsByTimeRange(startTime, endTime)
	if err != nil {
		fmt.Println("Error getting sessions by time range: ", err)
		return
	}

	var dayTaskIds []int64
	for _, session := range sessions {
		fmt.Println("Session: ", session.ID, session.TaskId)
		if contains[int64](dayTaskIds, session.TaskId) {
			continue
		}
		dayTaskIds = append(dayTaskIds, session.TaskId)

		task, err := taskStore.GetTaskByID(session.TaskId)
		if err != nil {
			fmt.Println("Error getting task by ID: ", err)
			return
		}

		fmt.Println("Task: ", task.ID, task.Title)
	}
}

func contains[T comparable](slice []T, element T) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}
