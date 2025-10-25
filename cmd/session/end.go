package session

import (
	"ttm/pkg/handlers"

	"github.com/spf13/cobra"
)

var endCmd = &cobra.Command{
	Use:   "end",
	Short: "End a session",
	Run:   handlers.EndHandler,
}

func init() {}
