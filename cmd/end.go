package cmd

import (
	"ttm/pkg/handlers"

	"github.com/spf13/cobra"
)

func NewEndCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "end",
		Short: "End a session",
		Run:   handlers.EndHandler,
	}
}
