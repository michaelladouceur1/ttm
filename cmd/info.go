package cmd

import (
	"ttm/pkg/handlers"

	"github.com/spf13/cobra"
)

func NewInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Get session info",
		Run:   handlers.InfoHandler,
	}
}
