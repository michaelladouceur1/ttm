package session

import (
	"fmt"
	"ttm/pkg/fs"
	"ttm/pkg/render"

	"github.com/spf13/cobra"
)

var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel a session",
	Run:   cancelHandler,
}

func init() {}

func cancelHandler(cmd *cobra.Command, args []string) {
	if !fs.SessionFileExists() {
		fmt.Println("No session found. Please start a session first.")
		return
	}

	_, err := fs.RemoveSessionFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	render.RenderSessionCancel()
}
