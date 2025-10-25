package session

import (
	"ttm/pkg/handlers"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get session info",
	Run:   handlers.InfoHandler,
}

func init() {}
