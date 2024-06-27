/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package session

import (
	"ttm/cmd"
	"ttm/pkg/store"
	"ttm/pkg/store/db"

	"github.com/spf13/cobra"
)

var taskStore = store.NewStore(db.NewDBStore())

// sessionCmd represents the session command
var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Manage sessions",
	Args:  cobra.MinimumNArgs(1),
}

// var listCmd = &cobra.Command{
// 	Use:   "list",
// 	Short: "List all sessions",
// 	Run:   listHandler,
// }

func init() {
	taskStore.Init()

	cmd.RootCmd.AddCommand(sessionCmd)

	sessionCmd.AddCommand(startCmd)
	sessionCmd.AddCommand(endCmd)
	sessionCmd.AddCommand(cancelCmd)
	sessionCmd.AddCommand(infoCmd)
	// sessionCmd.AddCommand(listCmd)

	// listCmd.Flags().StringP("id", "i", "", "Filter by task ID")
	// listCmd.Flags().StringP("days", "d", "", "Previous days to show")
}
