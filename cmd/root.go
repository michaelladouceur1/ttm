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

var ttmConfig = config.NewConfig()
var taskStore = store.NewStore(db.NewDBStore())

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ttm",
	Short: "Terminal Todo Manager",
}

func Execute() {
	ttmConfig.Init()
	ttmConfig.Load()

	taskStore.Init()

	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
