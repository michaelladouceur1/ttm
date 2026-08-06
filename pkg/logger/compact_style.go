package logger

import (
	"fmt"
	"strings"
	"ttm/pkg/models"
)

// CompactStyle is a color-free, line-oriented theme suited to narrow terminals.
type CompactStyle struct{}

func (CompactStyle) RenderMessage(message string) string {
	return message
}

func (CompactStyle) RenderError(message string) string {
	return "error: " + message
}

func (CompactStyle) RenderSummary(title string, items []SummaryItem) string {
	var output strings.Builder
	output.WriteString(title)
	for _, item := range items {
		fmt.Fprintf(&output, "\n%s: %s", item.Key, item.Value)
	}
	return output.String()
}

func (CompactStyle) RenderTasks(tasks []models.Task) string {
	var output strings.Builder
	output.WriteString("ID\tTitle\tStatus\tTags")
	for _, task := range tasks {
		fmt.Fprintf(&output, "\n%d\t%s\t%s\t%s", task.ListID, task.Title, task.Status, strings.Join(task.Tags, ","))
	}
	return output.String()
}
