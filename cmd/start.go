package cmd

import (
	"ttm/pkg/handlers"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func NewStartCmd(store *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "start [task_id]",
		Short: "Start a new session",
		Args:  cobra.MinimumNArgs(1),
		Run:   func(cmd *cobra.Command, args []string) { handlers.StartHandler(cmd, args, store) },
	}
}
