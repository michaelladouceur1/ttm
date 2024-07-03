package session

import (
	"fmt"
	"strconv"
	"time"
	"ttm/pkg/fs"
	"ttm/pkg/render"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [task_id]",
	Short: "Start a new session",
	Args:  cobra.MinimumNArgs(1),
	Run:   startHandler,
}

func init() {}

func startHandler(cmd *cobra.Command, args []string) {
	if fs.SessionFileExists() {
		fmt.Println("Session already started. Please end the current session first.")
		return
	}

	taskListIdString := args[0]
	taskListId, err := strconv.Atoi(taskListIdString)
	if err != nil {
		fmt.Println(err)
		return
	}

	taskId, err := fs.GetTaskIDFromListID(int64(taskListId))
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = fs.CreateSessionFile(taskId, time.Now())
	if err != nil {
		fmt.Println(err)
		return
	}

	render.RenderSessionStart(taskId)
}
