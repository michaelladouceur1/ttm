package cmd

import (
	"ttm/pkg/handlers"

	"github.com/spf13/cobra"
)

func NewCancelCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a session",
		Run:   handlers.CancelHandler,
	}
}
