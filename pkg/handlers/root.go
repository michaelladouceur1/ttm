package handlers

import (
	"fmt"
	"ttm/pkg/config"
	"ttm/pkg/store"
	"ttm/pkg/tui"

	"github.com/michaelladouceur1/gonfig"
	"github.com/spf13/cobra"
)

func RootHandler(cmd *cobra.Command, args []string, config *gonfig.Gonfig[config.Config], store *store.Store) {
	tui := tui.NewTUI(config, store)
	err := tui.Run()
	if err != nil {
		fmt.Println("Error running TUI: ", err)
	}
}
