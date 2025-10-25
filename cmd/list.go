package cmd

import (
	"ttm/pkg/config"
	"ttm/pkg/handlers"

	"github.com/spf13/cobra"
)

func NewListCmd(cfg *config.Config) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		Run:   func(cmd *cobra.Command, args []string) { handlers.ListHandler(cmd, args, cfg) },
	}

	listCategoryFlag := &cfg.ListFlags.Category
	listPriorityFlag := &cfg.ListFlags.Priority
	listStatusFlag := &cfg.ListFlags.Status

	listCmd.Flags().StringVarP(listCategoryFlag, "category", "c", *listCategoryFlag, "Filter tasks by category")
	listCmd.Flags().StringVarP(listPriorityFlag, "priority", "p", *listPriorityFlag, "Filter tasks by priority")
	listCmd.Flags().StringVarP(listStatusFlag, "status", "s", *listStatusFlag, "Filter tasks by status")

	return listCmd
}
