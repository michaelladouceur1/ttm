package ui

import (
	"fmt"
	"strconv"
	"strings"
	"ttm/pkg/fs"
	"ttm/pkg/models"
	"ttm/pkg/store"
)

func closeTasks(st *store.Store, taskListIDs []string) string {
	var ids []int64
	for _, taskListID := range taskListIDs {
		listID, err := strconv.ParseInt(taskListID, 10, 64)
		if err != nil {
			return "Error closing task: " + err.Error()
		}

		id, err := fs.GetTaskIDFromTempID(listID)
		if err != nil {
			return "Error closing task: " + err.Error()
		}

		ids = append(ids, id)
	}

	for _, id := range ids {
		if err := st.UpdateStatus(id, models.StatusClosed); err != nil {
			return "Error closing task: " + err.Error()
		}
	}

	var titles []string
	for _, id := range ids {
		task, err := st.GetTaskByID(id)
		if err != nil {
			return "Error retrieving closed task: " + err.Error()
		}
		titles = append(titles, task.Title)
	}

	var body strings.Builder
	body.WriteString("Tasks Closed\n")
	for i, title := range titles {
		fmt.Fprintf(&body, "Task %d: %s\n", i+1, title)
	}

	return strings.TrimRight(body.String(), "\n")
}
