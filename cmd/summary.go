package cmd

import (
	"ttm/pkg/handlers"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func NewSummaryCmd(store *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "summary",
		Short: "Summarize tasks",
		Run:   func(cmd *cobra.Command, args []string) { handlers.SummaryHandler(cmd, args, store) },
	}
}
