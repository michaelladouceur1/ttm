package cmd

import (
	"fmt"
	"os"
	"ttm/cmd/handlers"
	c "ttm/pkg/config"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task",
	Args:  cobra.MinimumNArgs(1),
	Run:   handlers.AddHandler,
}

func init() {
	config, err := c.Load()
	if err != nil {
		fmt.Println("Error loading config: ", err)
		os.Exit(1)
	}

	addCategoryFlag := &config.AddFlags.Category
	addPriorityFlag := &config.AddFlags.Priority
	addStatusFlag := &config.AddFlags.Status

	addCmd.Flags().StringVarP(addCategoryFlag, "category", "c", *addCategoryFlag, "Default category")
	addCmd.Flags().StringVarP(addPriorityFlag, "priority", "p", *addPriorityFlag, "Default priority")
	addCmd.Flags().StringVarP(addStatusFlag, "status", "s", *addStatusFlag, "Default status")
}
