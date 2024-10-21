package cmd

import (
	"fmt"
	"os"
	"ttm/cmd/handlers"
	c "ttm/pkg/config"

	"github.com/spf13/cobra"
)

var listTaskCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run:   handlers.ListHandler,
}

func init() {
	config, err := c.Load()
	if err != nil {
		fmt.Println("Error loading config: ", err)
		os.Exit(1)
	}

	listCategoryFlag := &config.ListFlags.Category
	listPriorityFlag := &config.ListFlags.Priority
	listStatusFlag := &config.ListFlags.Status

	listTaskCmd.Flags().StringVarP(listCategoryFlag, "category", "c", *listCategoryFlag, "Filter tasks by category")
	listTaskCmd.Flags().StringVarP(listPriorityFlag, "priority", "p", *listPriorityFlag, "Filter tasks by priority")
	listTaskCmd.Flags().StringVarP(listStatusFlag, "status", "s", *listStatusFlag, "Filter tasks by status")
}
