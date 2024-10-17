package cmd

import (
	"fmt"
	"strconv"
	"ttm/pkg/fs"
	"ttm/pkg/models"

	"github.com/spf13/cobra"
)

var closeCmd = &cobra.Command{
	Use:   "close",
	Short: "Close a task",
	Args:  cobra.MinimumNArgs(1),
	Run:   closeHandler,
}

func init() {}

func closeHandler(cmd *cobra.Command, args []string) {
	var ids []int64
	for _, arg := range args {
		tmpID, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Println("Error parsing task ID: ", err)
			return
		}

		id, err := fs.GetTaskIDFromListID(int64(tmpID))
		if err != nil {
			fmt.Println("Error getting task ID: ", err)
			return
		}

		ids = append(ids, id)
	}

	for _, id := range ids {
		err := taskStore.UpdateStatus(id, models.StatusClosed)
		if err != nil {
			fmt.Println("Error closing task: ", err)
			return
		}
	}
}
