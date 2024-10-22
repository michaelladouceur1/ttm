/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"ttm/pkg/config"
	"ttm/pkg/store"
	"ttm/pkg/store/db"

	"github.com/spf13/cobra"
)

var taskStore = store.NewStore(db.NewDBStore())

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ttm",
	Short: "Terminal Todo Manager",
}

func init() {
	RootCmd.AddCommand(addCmd)
	RootCmd.AddCommand(cancelCmd)
	RootCmd.AddCommand(closeCmd)
	RootCmd.AddCommand(endCmd)
	RootCmd.AddCommand(infoCmd)
	RootCmd.AddCommand(listTaskCmd)
	RootCmd.AddCommand(startCmd)
	RootCmd.AddCommand(updateCmd)
}

func Execute() {
	err := config.Init()
	if err != nil {
		os.Exit(1)
	}

	taskStore.Init()

	err = RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
