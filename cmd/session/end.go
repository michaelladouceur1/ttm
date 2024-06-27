package session

import (
	"fmt"
	"strconv"
	"time"
	"ttm/pkg/fs"
	"ttm/pkg/models"
	"ttm/pkg/render"

	"github.com/spf13/cobra"
)

var endCmd = &cobra.Command{
	Use:   "end",
	Short: "End a session",
	Run:   endHandler,
}

func init() {}

func endHandler(cmd *cobra.Command, args []string) {
	if !fs.SessionFileExists() {
		fmt.Println("No session found. Please start a session first.")
		return
	}

	sf, err := fs.RemoveSessionFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	taskId, err := strconv.Atoi(sf.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	taskStore.AddSession(models.Session{
		TaskId:    int64(taskId),
		StartTime: sf.StartTime,
		EndTime:   time.Now(),
	})

	render.RenderSessionEnd(sf)
}
