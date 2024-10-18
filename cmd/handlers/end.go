package handlers

import (
	"fmt"
	"time"
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/models"

	"github.com/spf13/cobra"
)

func EndHandler(cmd *cobra.Command, args []string) {
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

	logger.LogSessionEnd(sf)
}
