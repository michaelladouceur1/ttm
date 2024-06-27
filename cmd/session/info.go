package session

import (
	"fmt"
	"ttm/pkg/fs"
	"ttm/pkg/render"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get session info",
	Run:   infoHandler,
}

func init() {}

func infoHandler(cmd *cobra.Command, args []string) {
	if !fs.SessionFileExists() {
		fmt.Println("No session found. Please start a session first.")
		return
	}

	sf, err := fs.ReadSessionFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	render.RenderSessionInfo(sf)
}
