package logger

import (
	"fmt"
	"strings"
	"ttm/pkg/models"
)

func (l *Logger) LogAddTask(task models.Task) {
	fmt.Println(l.style.RenderSummary("Task Added", taskSummaryItems(task)))
}

func (l *Logger) LogUpdateTask(task models.Task) {
	fmt.Println(l.style.RenderSummary("Task Updated", taskSummaryItems(task)))
}

func (l *Logger) LogTasks(tasks []models.Task) {
	fmt.Println(l.style.RenderTasks(tasks))
}

func (l *Logger) RenderTasks(tasks []models.Task) string {
	return l.style.RenderTasks(tasks)
}

func (l *Logger) LogCloseTasks(tasks []models.Task) {
	items := make([]SummaryItem, 0, len(tasks))
	for i, task := range tasks {
		items = append(items, SummaryItem{
			Key:   fmt.Sprintf("Task %d", i+1),
			Value: task.Title,
		})
	}
	fmt.Println(l.style.RenderSummary("Tasks Closed", items))
}

func (l *Logger) LogSessionSummary(taskSummary models.TaskSummary) {
	for _, day := range taskSummary.Days {
		if len(day.Tasks) == 0 {
			continue
		}
		items := make([]SummaryItem, 0, len(day.Tasks))
		for _, task := range day.Tasks {
			items = append(items, SummaryItem{Key: task.Title, Value: task.Description})
		}
		fmt.Println(l.style.RenderSummary(day.Day.Format("2006-01-02"), items))
	}
}

func LogAddTask(task models.Task) {
	defaultLogger.LogAddTask(task)
}

func LogUpdateTask(task models.Task) {
	defaultLogger.LogUpdateTask(task)
}

func LogTasks(tasks []models.Task) {
	defaultLogger.LogTasks(tasks)
}

func RenderTasks(tasks []models.Task) string {
	return defaultLogger.RenderTasks(tasks)
}

func LogCloseTasks(tasks []models.Task) {
	defaultLogger.LogCloseTasks(tasks)
}

func LogSessionSummary(taskSummary models.TaskSummary) {
	defaultLogger.LogSessionSummary(taskSummary)
}

func taskSummaryItems(task models.Task) []SummaryItem {
	return []SummaryItem{
		{Key: "Title", Value: task.Title},
		{Key: "Description", Value: task.Description},
		{Key: "Priority", Value: string(task.Priority)},
		{Key: "Status", Value: string(task.Status)},
		{Key: "Tags", Value: strings.Join(task.Tags, ",")},
	}
}
