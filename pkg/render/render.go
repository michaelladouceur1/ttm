package render

import (
	"fmt"
	"time"
	"ttm/pkg/models"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Renderer interface {
	RenderTasks(tasks []models.Task)
	RenderAddTask(task models.Task)
}

func RenderAddTask(task models.Task) {
	fmt.Printf("Task added\n")
}

func RenderUpdateTask() {
	fmt.Printf("Task updated\n")
}

func RenderTasks(tasks []models.Task) {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Temp ID", "Title", "Description", "Category", "Priority", "Status", "Duration", "Created At"})
	for _, task := range tasks {
		t.AppendRow(table.Row{task.ListID, task.Title, task.Description, task.Category, task.Priority, task.Status, formatTimeToDuration(task.Duration), task.CreatedAt.Format("2006-01-02 15:04:05")})
	}
	t.SetStyle(table.StyleColoredDark)
	fmt.Println(t.Render())
}

func RenderSessionStart(taskId string) {
	fmt.Printf("Session started for task: %s \n", taskId)
}

func RenderSessionEnd(session models.SessionFile) {
	RenderSessionInfo(session)
}

func RenderSessionInfo(session models.SessionFile) {
	fmt.Printf("Current session for task: %s \n", session.ID)
	fmt.Printf("Start time: %s \n", session.StartTime)
	fmt.Printf("Duration: %s \n", time.Since(session.StartTime))
}

func RenderSessionCancel() {
	fmt.Printf("Session cancelled\n")
}

func formatTimeToDuration(t time.Time) string {
	return fmt.Sprintf("%02dh%02dm%02ds", int(t.Hour()), int(t.Minute()), int(t.Second()))
}
