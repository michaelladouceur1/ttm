/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"ttm/pkg/models"
	"ttm/pkg/render"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run:   listHandler,
}

var listCategoryFlag = &ttmConfig.ListFlags.Category
var listPriorityFlag = &ttmConfig.ListFlags.Priority
var listStatusFlag = &ttmConfig.ListFlags.Status

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(listCategoryFlag, "category", "c", *listCategoryFlag, "Filter tasks by category")
	listCmd.Flags().StringVarP(listPriorityFlag, "priority", "p", *listPriorityFlag, "Filter tasks by priority")
	listCmd.Flags().StringVarP(listStatusFlag, "status", "s", *listStatusFlag, "Filter tasks by status")
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
