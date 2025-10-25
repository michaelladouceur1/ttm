package cmd

import (
	"ttm/pkg/handlers"

	"github.com/spf13/cobra"
)

func NewUpdateCmd() *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update a task",
		Args:  cobra.MinimumNArgs(1),
		Run:   handlers.UpdateHandler,
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
