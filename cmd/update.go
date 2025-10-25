package cmd

import (
	"ttm/pkg/handlers"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func NewUpdateCmd(store *store.Store) *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update a task",
		Args:  cobra.MinimumNArgs(1),
		Run:   func(cmd *cobra.Command, args []string) { handlers.UpdateHandler(cmd, args, store) },
	}

	updateCmd.Flags().IntP("id", "i", 0, "Task ID")
	updateCmd.Flags().StringP("title", "t", "", "Update title")
	updateCmd.Flags().StringP("description", "d", "", "Update description")
	updateCmd.Flags().StringP("category", "c", "", "Update category")
	updateCmd.Flags().StringP("priority", "p", "", "Update priority")
	updateCmd.Flags().StringP("status", "s", "", "Update status")
	updateCmd.Flags().StringP("openedAt", "a", "", "Update opened time")
	updateCmd.Flags().StringP("closedAt", "b", "", "Update closed time")

	return updateCmd
}
