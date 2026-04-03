package cmd

import (
	"ttm/pkg/config"
	"ttm/pkg/handlers"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func NewSummaryCmd(cfg *config.Config, store *store.Store) *cobra.Command {
	summaryCmd := &cobra.Command{
		Use:   "summary",
		Short: "Summarize tasks",
		Run:   func(cmd *cobra.Command, args []string) { handlers.SummaryHandler(cmd, args, store) },
	}

	daysToDisplay := &cfg.DaysToDisplay
	summaryCmd.Flags().IntVarP(daysToDisplay, "days", "d", *daysToDisplay, "Number of days to display in the summary")

	return summaryCmd
}
