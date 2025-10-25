package cmd

import (
	"ttm/pkg/config"
	"ttm/pkg/handlers"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func NewAddCmd(cfg *config.Config, st *store.Store) *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new task",
		Args:  cobra.MinimumNArgs(1),
		Run:   func(cmd *cobra.Command, args []string) { handlers.AddHandler(cmd, args, st) },
	}

	addCategoryFlag := &cfg.AddFlags.Category
	addPriorityFlag := &cfg.AddFlags.Priority
	addStatusFlag := &cfg.AddFlags.Status
	addCmd.Flags().StringVarP(addCategoryFlag, "category", "c", *addCategoryFlag, "Default category")
	addCmd.Flags().StringVarP(addPriorityFlag, "priority", "p", *addPriorityFlag, "Default priority")
	addCmd.Flags().StringVarP(addStatusFlag, "status", "s", *addStatusFlag, "Default status")

	return addCmd
}
