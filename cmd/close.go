package cmd

import (
	"ttm/pkg/handlers"

	"github.com/spf13/cobra"
)

func NewCloseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "close",
		Short: "Close a task",
		Args:  cobra.MinimumNArgs(1),
		Run:   handlers.CloseHandler,
	}
}
