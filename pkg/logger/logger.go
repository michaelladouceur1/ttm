package logger

import (
	"fmt"
	"strconv"
	"time"
	"ttm/pkg/models"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Renderer interface {
	LogTasks(tasks []models.Task)
	LogAddTask(task models.Task)
}

func LogAddTask(task models.Task) {
	fmt.Printf("Task added\n")
}

func LogUpdateTask() {
	fmt.Printf("Task updated\n")
}

func LogTasks(tasks []models.Task) {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Temp ID", "ID", "Title", "Description", "Category", "Priority", "Status", "Duration", "Created At"})
	for _, task := range tasks {
		t.AppendRow(table.Row{task.ListID, task.ID, task.Title, task.Description, task.Category, task.Priority, task.Status, formatTimeToDuration(task.Duration), task.CreatedAt.Format("2006-01-02 15:04:05")})
	}
	t.SetStyle(table.StyleColoredDark)
	fmt.Println(t.Render())
}

func LogTaskSummary(taskSummary models.TaskSummary) {
	for _, day := range taskSummary.Days {
		if len(day.Tasks) == 0 {
			continue
		}
		fmt.Printf("%s: %s\n", day.Day.Weekday().String(), day.Day.Format("2006-01-02"))
		for _, task := range day.Tasks {
			task.CalculateDuration()
			fmt.Printf("   • %s: %s\n", task.Title, task.Description)
		}
		fmt.Println()
	}
}

func LogSessionStart(taskId int64) {
	taskIdString := strconv.Itoa(int(taskId))
	fmt.Printf("Session started for task: %s \n", taskIdString)
}

func LogSessionEnd(session models.SessionFile) {
	LogSessionInfo(session)
}

func LogSessionInfo(session models.SessionFile) {
	taskIdString := strconv.Itoa(int(session.ID))
	fmt.Printf("Current session for task: %s \n", taskIdString)
	fmt.Printf("Start time: %s \n", session.StartTime)
	fmt.Printf("Duration: %s \n", time.Since(session.StartTime))
}

func LogSessionCancel() {
	fmt.Printf("Session cancelled\n")
}

func formatTimeToDuration(t time.Time) string {
	return fmt.Sprintf("%02dh%02dm%02ds", int(t.Hour()), int(t.Minute()), int(t.Second()))
}
