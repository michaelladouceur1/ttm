package handlers

import (
	"fmt"
	"ttm/pkg/config"
	"ttm/pkg/tui"

	"github.com/michaelladouceur1/gonfig"
	"github.com/spf13/cobra"
)

func RootHandler(cmd *cobra.Command, args []string, config *gonfig.Gonfig[config.Config]) {
	tui := tui.NewTUI(config)
	err := tui.Run()
	if err != nil {
		fmt.Println("Error running TUI: ", err)
	}
}
