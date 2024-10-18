package logger

import (
	"fmt"
	"time"
	"ttm/pkg/models"

	"github.com/jedib0t/go-pretty/v6/table"
)

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

func formatTimeToDuration(t time.Time) string {
	return fmt.Sprintf("%02dh%02dm%02ds", int(t.Hour()), int(t.Minute()), int(t.Second()))
}
