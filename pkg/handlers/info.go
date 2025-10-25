package handlers

import (
	"fmt"
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func InfoHandler(cmd *cobra.Command, args []string, store *store.Store) {
	if !fs.SessionFileExists() {
		fmt.Println("No session found. Please start a session first.")
		return
	}

	sf, err := fs.ReadSessionFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	task, err := store.GetTaskByID(sf.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	logger.LogSessionInfo(sf, task)
}
