package logger

import (
	"fmt"
	"strconv"
	"strings"
	"ttm/pkg/models"
	"ttm/pkg/styles"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/lipgloss/tree"
)

type ClassicStyle struct {
	TermWidth  int
	TermHeight int
}

func (ClassicStyle) RenderMessage(message string) string {
	return lipgloss.NewStyle().Bold(true).Foreground(styles.Main).Render(message)
}

func (ClassicStyle) RenderError(message string) string {
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#d34545ff")).Render(message)
}

func (ClassicStyle) RenderSummary(title string, items []SummaryItem) string {
	longestKey := 0
	for _, item := range items {
		if len(item.Key) > longestKey {
			longestKey = len(item.Key)
		}
	}

	connectorStyle := lipgloss.NewStyle().Foreground(styles.Highlight2)
	textStyle := lipgloss.NewStyle().Foreground(styles.Highlight1)
	children := make([]any, 0, len(items))
	for _, item := range items {
		padding := longestKey - len(item.Key) + 2
		value := item.Key + connectorStyle.Render(" "+strings.Repeat("─", padding)) + " " + textStyle.Render(item.Value)
		children = append(children, value)
	}

	return tree.Root("⚙ " + title).
		Child(children...).
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(connectorStyle).
		RootStyle(lipgloss.NewStyle().Bold(true).Foreground(styles.Main)).
		ItemStyle(textStyle).
		String()
}

func (cs ClassicStyle) RenderTasks(tasks []models.Task) string {
	cellStyle := lipgloss.NewStyle().Padding(0, 1).Width(14)
	columnStyles := []lipgloss.Style{
		cellStyle.Width(5),
		cellStyle.Width(20),
		cellStyle.Width(30),
		cellStyle.Width(10),
		cellStyle.Width(10),
		cellStyle.Width(8),
		cellStyle.Width(21),
		cellStyle,
		cellStyle.Width(21),
	}
	headerStyle := lipgloss.NewStyle().Foreground(styles.Main).Bold(true).Align(lipgloss.Center)
	table := table.New().
		Width(max(1, cs.TermWidth-5)).
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(styles.Main)).
		Headers("ID", "Title", "Description", "Category", "Priority", "Status", "Tags", "Duration", "Created At").
		StyleFunc(func(row, column int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			style := columnStyles[column]
			if row%2 == 0 {
				return style.Foreground(styles.Highlight3)
			}
			return style.Foreground(styles.Highlight2)
		})
	for _, task := range tasks {
		table.Row(
			strconv.FormatInt(task.ListID, 10),
			task.Title,
			task.Description,
			string(task.Category),
			string(task.Priority),
			string(task.Status),
			strings.Join(task.Tags, ","),
			fmt.Sprintf("%02dh%02dm%02ds", int(task.Duration.Hour()), int(task.Duration.Minute()), int(task.Duration.Second())),
			task.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	}
	return table.String()
}
