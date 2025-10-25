package handlers

import (
	"fmt"
	"time"
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/models"
	"ttm/pkg/store"

	"github.com/spf13/cobra"
)

func EndHandler(cmd *cobra.Command, args []string, store *store.Store) {
	if !fs.SessionFileExists() {
		fmt.Println("No session found. Please start a session first.")
		return
	}

	sf, err := fs.RemoveSessionFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	store.AddSession(models.Session{
		TaskId:    int64(sf.ID),
		StartTime: sf.StartTime,
		EndTime:   time.Now(),
	})

	task, err := store.GetTaskByID(sf.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	logger.LogSessionEnd(sf, task)
}
