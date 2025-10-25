package cmd

import (
	"ttm/pkg/handlers"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func NewEndCmd(store *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "end",
		Short: "End a session",
		Run:   func(cmd *cobra.Command, args []string) { handlers.EndHandler(cmd, args, store) },
	}
}
