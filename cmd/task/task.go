/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package task

import (
	"ttm/cmd"
	"ttm/pkg/config"
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

func init() {
	ttmConfig.Load()
	taskStore.Init()

	cmd.RootCmd.AddCommand(taskCmd)

	taskCmd.AddCommand(addCmd)
	taskCmd.AddCommand(listTaskCmd)
	taskCmd.AddCommand(updateCmd)
	taskCmd.AddCommand(csvCmd)

}
