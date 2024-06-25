/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"ttm/pkg/fs"
	"ttm/pkg/models"

	"github.com/spf13/cobra"
)

// csvCmd represents the csv command
var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "Export tasks to CSV",
	Run:   csvHandler,
}

func init() {
	rootCmd.AddCommand(csvCmd)

	csvCmd.Flags().StringP("category", "c", "", "Filter tasks by category")
	csvCmd.Flags().StringP("status", "s", "", "Filter tasks by status")
	csvCmd.Flags().StringP("priority", "p", "", "Filter tasks by priority")
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
