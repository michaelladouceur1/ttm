package cmd

import (
	"ttm/pkg/handlers"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func NewCloseCmd(st *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "close",
		Short: "Close a task",
		Args:  cobra.MinimumNArgs(1),
		Run:   func(cmd *cobra.Command, args []string) { handlers.CloseHandler(cmd, args, st) },
	}
}
