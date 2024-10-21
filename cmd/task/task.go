/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package task

import (
	"ttm/cmd"

	"github.com/spf13/cobra"
)

// taskCmd represents the task command
var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks",
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	cmd.RootCmd.AddCommand(taskCmd)

	taskCmd.AddCommand(csvCmd)
	taskCmd.AddCommand(summaryCmd)

}
