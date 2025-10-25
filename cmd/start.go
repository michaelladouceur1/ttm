package cmd

import (
	"ttm/pkg/handlers"

	"github.com/spf13/cobra"
)

func NewStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start [task_id]",
		Short: "Start a new session",
		Args:  cobra.MinimumNArgs(1),
		Run:   handlers.StartHandler,
	}
}
