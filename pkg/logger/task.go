package logger

import (
	"fmt"
	"strconv"
	"time"
	"ttm/pkg/models"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/lipgloss/tree"
)

const (
	TempIdColumn = iota
	IdColumn
	TitleColumn
	DescriptionColumn
	CategoryColumn
	PriorityColumn
	StatusColumn
	DurationColumn
	CreatedAtColumn
)

func LogAddTask(task models.Task) {
	fmt.Println(createTaskSummaryTree(task, "Task Added"))
}

func LogUpdateTask(task models.Task) {
	fmt.Println(createTaskSummaryTree(task, "Task Updated"))
}

func LogTasks(tasks []models.Task) {
	t := createTable()
	t.Headers("Temp ID", "ID", "Title", "Description", "Category", "Priority", "Status", "Duration", "Created At")
	for _, task := range tasks {
		t.Row(createRowStrings(task)...)
	}
	fmt.Println(t)
}

func LogCloseTasks(tasks []models.Task) {
	data := []SummaryTreeItem{}
	for i, task := range tasks {
		data = append(data, SummaryTreeItem{
			Key:   fmt.Sprintf("Task %d", i+1),
			Value: task.Title,
		})
	}
	fmt.Println(createSummaryTree(data, "Tasks Closed"))
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

func createTaskSummaryTree(task models.Task, title string) *tree.Tree {
	data := []SummaryTreeItem{
		{"Title", task.Title},
		{"Description", task.Description},
		{"Category", string(task.Category)},
		{"Priority", string(task.Priority)},
		{"Status", string(task.Status)},
	}
	return createSummaryTree(data, title)
}

func createTable() *table.Table {
	return table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(blue)).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return tableHeaderStyle
			}

			var style lipgloss.Style
			switch col {
			case TempIdColumn:
				style = tempIdStyle
			case IdColumn:
				style = idStyle
			case TitleColumn:
				style = titleStyle
			case DescriptionColumn:
				style = descriptionStyle
			case CategoryColumn:
				style = categoryStyle
			case PriorityColumn:
				style = priorityStyle
			case StatusColumn:
				style = statusStyle
			case DurationColumn:
				style = cellStyle
			case CreatedAtColumn:
				style = createdAtStyle
			default:
				style = cellStyle
			}

			switch row % 2 {
			case 0: // even
				style = style.Foreground(lightGray)
			default: // odd
				style = style.Foreground(gray)
			}

			return style
		})
}

func createRowStrings(task models.Task) []string {
	return []string{
		toIntString(task.ListID),
		toIntString(task.ID),
		task.Title,
		task.Description,
		string(task.Category),
		string(task.Priority),
		string(task.Status),
		toDuration(task.Duration),
		task.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func toIntString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func toDuration(t time.Time) string {
	return fmt.Sprintf("%02dh%02dm%02ds", int(t.Hour()), int(t.Minute()), int(t.Second()))
}
