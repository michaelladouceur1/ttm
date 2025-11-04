package handlers

import (
	"fmt"
	"ttm/pkg/tui"

	"github.com/spf13/cobra"
)

func RootHandler(cmd *cobra.Command, args []string) {
	tui := tui.NewTUI()
	err := tui.Run()
	if err != nil {
		fmt.Println("Error running TUI: ", err)
	}
}
