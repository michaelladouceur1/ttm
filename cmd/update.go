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

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a task",
	Run:   updateHandler,
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().IntP("id", "i", 0, "Task ID")
	updateCmd.Flags().StringP("title", "t", "", "Update title")
	updateCmd.Flags().StringP("description", "d", "", "Update description")
	updateCmd.Flags().StringP("category", "c", "", "Update category")
	updateCmd.Flags().StringP("priority", "p", "", "Update priority")
	updateCmd.Flags().StringP("status", "s", "", "Update status")
	updateCmd.Flags().StringP("start", "a", "", "Update start time")
	updateCmd.Flags().StringP("end", "b", "", "Update end time")
}

func updateHandler(cmd *cobra.Command, args []string) {
	idFlag, _ := cmd.Flags().GetInt("id")
	titleFlag, _ := cmd.Flags().GetString("title")
	descriptionFlag, _ := cmd.Flags().GetString("description")
	categoryFlag, _ := cmd.Flags().GetString("category")
	priorityFlag, _ := cmd.Flags().GetString("priority")
	statusFlag, _ := cmd.Flags().GetString("status")
	startTimeFlag, _ := cmd.Flags().GetString("start")
	endTimeFlag, _ := cmd.Flags().GetString("end")

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

	if titleFlag == "" && descriptionFlag == "" && categoryFlag == "" && priorityFlag == "" && statusFlag == "" && startTimeFlag == "" && endTimeFlag == "" {
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

	if startTimeFlag != "" {
		err = taskStore.UpdateStartTime(idFlag, startTimeFlag)
		if err != nil {
			fmt.Println("Error updating start time: ", err)
			return
		}
	}

	if endTimeFlag != "" {
		err = taskStore.UpdateEndTime(idFlag, endTimeFlag)
		if err != nil {
			fmt.Println("Error updating end time: ", err)
			return
		}
	}

	render.RenderUpdateTask()
}
