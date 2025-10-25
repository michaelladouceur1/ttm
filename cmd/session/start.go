package session

import (
	"ttm/pkg/handlers"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [task_id]",
	Short: "Start a new session",
	Args:  cobra.MinimumNArgs(1),
	Run:   handlers.StartHandler,
}

func init() {}
