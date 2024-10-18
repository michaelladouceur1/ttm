package cmd

import (
	"ttm/cmd/handlers"

	"github.com/spf13/cobra"
)

var closeCmd = &cobra.Command{
	Use:   "close",
	Short: "Close a task",
	Args:  cobra.MinimumNArgs(1),
	Run:   handlers.CloseHandler,
}

func init() {}
