package cmd

import (
	"ttm/pkg/handlers"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func NewSearchCmd(store *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "search <query>",
		Short: "Search tasks",
		Long: `Search tasks by a substring across all task fields, or use field filters.

Examples:
  ttm search "release notes"
  ttm search "*tags:work,urgent *title:Task 1"`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			handlers.SearchHandler(cmd, args, store)
		},
	}
}
