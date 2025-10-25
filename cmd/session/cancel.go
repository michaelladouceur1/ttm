package session

import (
	"ttm/pkg/handlers"

	"github.com/spf13/cobra"
)

var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel a session",
	Run:   handlers.CancelHandler,
}

func init() {}
