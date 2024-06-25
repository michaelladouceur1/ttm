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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ttm",
	Short: "Terminal Todo Manager",
}

var ttmConfig = config.NewConfig()
var taskStore = store.NewStore(db.NewDBStore())

func Execute() {
	ttmConfig.Init()
	ttmConfig.Load()

	taskStore.Init()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
