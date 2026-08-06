/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"ttm/pkg/config"
	"ttm/pkg/logger"
	"ttm/pkg/store"
	"ttm/pkg/store/googledocs"
	postgres "ttm/pkg/store/postgres"
	sqlite "ttm/pkg/store/sqlite"

	"github.com/spf13/cobra"
)

// var taskStore = store.NewStore(db.NewDBStore())

// rootCmd represents the base command when called without any subcommands
var RootCmd *cobra.Command

func init() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println("Error initializing config: ", err)
		os.Exit(1)
	}

	if err := logger.Configure(cfg.Config.Logging.Theme); err != nil {
		fmt.Println("Error initializing logger: ", err)
		os.Exit(1)
	}

	var strategy store.StoreStrategy
	switch cfg.Config.Storage.Type {
	case "", "sqlite":
		strategy = sqlite.NewStore()
	case "postgres":
		strategy = postgres.NewStore()
	case "google-docs":
		strategy = googledocs.NewStore(cfg.Config.Storage.GoogleDocs)
	default:
		fmt.Printf("Error initializing storage: unsupported storage type %q\n", cfg.Config.Storage.Type)
		os.Exit(1)
	}

	st := store.NewStore(strategy)
	if err := st.Init(); err != nil {
		fmt.Println("Error initializing storage: ", err)
		os.Exit(1)
	}

	RootCmd = &cobra.Command{
		Use:   "ttm",
		Short: "Terminal Todo Manager",
		// Run:   func(cmd *cobra.Command, args []string) { handlers.RootHandler(cmd, args, cfg) },
	}

	RootCmd.AddCommand(NewAddCmd(cfg.Config, st))
	RootCmd.AddCommand(NewCancelCmd())
	RootCmd.AddCommand(NewCloseCmd(st))
	RootCmd.AddCommand(NewEndCmd(st))
	RootCmd.AddCommand(NewInfoCmd(st))
	RootCmd.AddCommand(NewListCmd(cfg.Config, st))
	RootCmd.AddCommand(NewSearchCmd(st))
	RootCmd.AddCommand(NewStartCmd(st))
	RootCmd.AddCommand(NewUpdateCmd(st))
	RootCmd.AddCommand(NewSummaryCmd(cfg.Config, st))
	RootCmd.AddCommand(NewTagCmd(st))
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
