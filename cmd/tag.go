package cmd

import (
	"ttm/pkg/handlers"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func NewTagCmd(store *store.Store) *cobra.Command {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Tag a task",
		Args:  cobra.MinimumNArgs(2),
		Run:   func(cmd *cobra.Command, args []string) { handlers.TagHandler(cmd, args, store) },
	}

	return tagCmd
}
