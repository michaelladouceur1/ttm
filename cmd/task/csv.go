package task

import (
	"ttm/pkg/handlers"

	"github.com/spf13/cobra"
)

var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "Export tasks to CSV",
	Run:   handlers.CSVHandler,
}

func init() {
	csvCmd.Flags().StringP("category", "c", "", "Filter tasks by category")
	csvCmd.Flags().StringP("status", "s", "", "Filter tasks by status")
	csvCmd.Flags().StringP("priority", "p", "", "Filter tasks by priority")
}
