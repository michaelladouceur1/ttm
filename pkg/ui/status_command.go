package ui

import (
	"fmt"
	"strconv"
	"strings"
	"ttm/pkg/fs"
	"ttm/pkg/models"
	"ttm/pkg/store"
)

func openTasks(st *store.Store, taskListIDs []string) string {
	return updateTasksStatus(st, taskListIDs, models.StatusOpen, "Tasks Opened")
}

func standbyTasks(st *store.Store, taskListIDs []string) string {
	return updateTasksStatus(st, taskListIDs, models.StatusStandby, "Tasks Put on Standby")
}

func updateTasksStatus(st *store.Store, taskListIDs []string, status models.Status, heading string) string {
	ids := make([]int64, 0, len(taskListIDs))
	for _, taskListID := range taskListIDs {
		listID, err := strconv.ParseInt(taskListID, 10, 64)
		if err != nil {
			return "Error updating task status: " + err.Error()
		}

		id, err := fs.GetTaskIDFromTempID(listID)
		if err != nil {
			return "Error updating task status: " + err.Error()
		}

		ids = append(ids, id)
	}

	for _, id := range ids {
		if err := st.UpdateStatus(id, status); err != nil {
			return "Error updating task status: " + err.Error()
		}
	}

	var body strings.Builder
	body.WriteString(heading + "\n")
	for i, id := range ids {
		task, err := st.GetTaskByID(id)
		if err != nil {
			return "Error retrieving updated task: " + err.Error()
		}
		fmt.Fprintf(&body, "Task %d: %s\n", i+1, task.Title)
	}

	return strings.TrimRight(body.String(), "\n")
}
