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

var (
	// Colors
	blue      = lipgloss.Color("#2cd5fb")
	darkBlue  = lipgloss.Color("#219cb8ff")
	gray      = lipgloss.Color("#6b6b6bff")
	lightGray = lipgloss.Color("#e4e4e4ff")

	// General styles
	headerStyle   = lipgloss.NewStyle().Bold(true).Foreground(blue)
	textStyle     = lipgloss.NewStyle().Foreground(darkBlue)
	treeConnStyle = lipgloss.NewStyle().Foreground(gray)

	// Table styles
	// Header style
	tableHeaderStyle = lipgloss.NewStyle().Foreground(blue).Bold(true).Align(lipgloss.Center)

	// Cell styles
	cellStyle        = lipgloss.NewStyle().Padding(0, 1).Width(14)
	tempIdStyle      = cellStyle.Width(9)
	idStyle          = cellStyle.Width(5)
	titleStyle       = cellStyle.Width(20)
	descriptionStyle = cellStyle.Width(30)
	categoryStyle    = cellStyle.Width(10)
	priorityStyle    = cellStyle.Width(10)
	statusStyle      = cellStyle.Width(8)
	createdAtStyle   = cellStyle.Width(21)

	// Row styles
	oddRowStyle  = cellStyle.Foreground(gray)
	evenRowStyle = cellStyle.Foreground(lightGray)
)

func LogAddTask(task models.Task) {
	fmt.Println(createSummaryTree(task, "Task Added"))
}

func LogUpdateTask(task models.Task) {
	fmt.Println(createSummaryTree(task, "Task Updated"))
}

func LogTasks(tasks []models.Task) {
	t := createTable()
	t.Headers("Temp ID", "ID", "Title", "Description", "Category", "Priority", "Status", "Duration", "Created At")
	for _, task := range tasks {
		t.Row(createRowStrings(task)...)
	}
	fmt.Println(t)
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

func createSummaryTree(task models.Task, title string) *tree.Tree {
	return tree.Root("⚙ "+title).
		Child(
			"Title"+fmt.Sprintf(treeConnStyle.Render(" ──────── "))+fmt.Sprintf(textStyle.Render(task.Title)),
			"Description"+fmt.Sprintf(treeConnStyle.Render(" ── "))+fmt.Sprintf(textStyle.Render(task.Description)),
			"Category"+fmt.Sprintf(treeConnStyle.Render(" ───── "))+fmt.Sprintf(textStyle.Render(string(task.Category))),
			"Priority"+fmt.Sprintf(treeConnStyle.Render(" ───── "))+fmt.Sprintf(textStyle.Render(string(task.Priority))),
			"Status"+fmt.Sprintf(treeConnStyle.Render(" ─────── "))+fmt.Sprintf(textStyle.Render(string(task.Status))),
		).
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(treeConnStyle).
		RootStyle(headerStyle).
		ItemStyle(textStyle)
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
