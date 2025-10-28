/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"ttm/pkg/config"
	"ttm/pkg/store"
	"ttm/pkg/store/db"

	"github.com/spf13/cobra"
)

// var taskStore = store.NewStore(db.NewDBStore())

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ttm",
	Short: "Terminal Todo Manager",
}

func init() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println("Error initializing config: ", err)
		os.Exit(1)
	}

	store := store.NewStore(db.NewDBStore("sqlite"))
	store.Init()

	RootCmd.AddCommand(NewAddCmd(cfg.Config, store))
	RootCmd.AddCommand(NewCancelCmd())
	RootCmd.AddCommand(NewCloseCmd(store))
	RootCmd.AddCommand(NewEndCmd(store))
	RootCmd.AddCommand(NewInfoCmd(store))
	RootCmd.AddCommand(NewListCmd(cfg.Config, store))
	RootCmd.AddCommand(NewStartCmd(store))
	RootCmd.AddCommand(NewUpdateCmd(store))
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
