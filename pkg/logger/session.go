package logger

import (
	"fmt"
	"time"
	"ttm/pkg/models"

	"github.com/charmbracelet/lipgloss/tree"
)

func LogSessionStart(task models.Task, start time.Time) {
	fmt.Println(createSessionStartSummaryTree(task, start, "Session Started"))
}

func LogSessionEnd(session models.SessionFile, task models.Task) {
	fmt.Println(createSessionSummaryTree(session, task, "Session Ended"))
}

func LogSessionInfo(session models.SessionFile, task models.Task) {
	fmt.Println(createSessionSummaryTree(session, task, "Session Info"))
}

func LogSessionCancel() {
	LogMessage("Session cancelled.")
}

func createSessionStartSummaryTree(task models.Task, start time.Time, title string) *tree.Tree {
	data := []SummaryTreeItem{
		{"Task Title", task.Title},
		{"Task Description", task.Description},
		{"Start Time", start.Round(time.Second).Format("2006-01-02 15:04:05")},
	}
	return createSummaryTree(data, title)
}

func createSessionSummaryTree(session models.SessionFile, task models.Task, title string) *tree.Tree {
	data := []SummaryTreeItem{
		{"Task Title", task.Title},
		{"Start Time", session.StartTime.Round(time.Second).Format("2006-01-02 15:04:05")},
		{"Duration", time.Since(session.StartTime).Round(time.Second).String()},
	}
	return createSummaryTree(data, title)
}
