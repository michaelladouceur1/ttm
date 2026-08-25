package cmd

import (
	"strings"
	"ttm/pkg/config"
	"ttm/pkg/handlers"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func NewListCmd(cfg *config.Config, store *store.Store) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		Run:   func(cmd *cobra.Command, args []string) { handlers.ListHandler(cmd, args, cfg, store) },
	}

	listPriorityFlag := strings.Join(cfg.ListFlags.Priority, ",")
	listStatusFlag := strings.Join(cfg.ListFlags.Status, ",")

	listCmd.Flags().StringVarP(&listPriorityFlag, "priority", "p", listPriorityFlag, "Filter tasks by priority")
	listCmd.Flags().StringVarP(&listStatusFlag, "status", "s", listStatusFlag, "Filter tasks by status")

	return listCmd
}
