package cmd

import (
	"ttm/pkg/handlers"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func NewInfoCmd(store *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Get session info",
		Run:   func(cmd *cobra.Command, args []string) { handlers.InfoHandler(cmd, args, store) },
	}
}
