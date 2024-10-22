package cmd

import (
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:    "completion",
	Hidden: true,
}

func init() {
	RootCmd.AddCommand(completionCmd)
}
