package handlers

import (
	"fmt"
	"ttm/pkg/fs"
	"ttm/pkg/logger"

	"github.com/spf13/cobra"
)

func InfoHandler(cmd *cobra.Command, args []string) {
	if !fs.SessionFileExists() {
		fmt.Println("No session found. Please start a session first.")
		return
	}

	sf, err := fs.ReadSessionFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	logger.LogSessionInfo(sf)
}
