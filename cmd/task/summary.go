/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package task

import (
	"ttm/cmd/handlers"

	"github.com/spf13/cobra"
)

// summaryCmd represents the summary command
var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Summarize tasks",
	Run:   handlers.SummaryHandler,
}

func init() {
	summaryCmd.Flags().IntP("days", "d", 7, "Number of days to summarize")
}
