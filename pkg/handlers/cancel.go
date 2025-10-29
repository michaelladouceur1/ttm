package handlers

import (
	"ttm/pkg/fs"
	"ttm/pkg/logger"

	"github.com/spf13/cobra"
)

func CancelHandler(cmd *cobra.Command, args []string) {
	if !fs.SessionFileExists() {
		logger.LogMessage("No session found. Please start a session first.")
		return
	}

	_, err := fs.RemoveSessionFile()
	if err != nil {
		logger.LogError("Error cancelling session: ", err)
		return
	}

	logger.LogSessionCancel()
}
