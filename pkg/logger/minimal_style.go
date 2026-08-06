package logger

import (
	"strconv"
	"strings"
	"ttm/pkg/models"
	"ttm/pkg/styles"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/lipgloss/tree"
)

const (
	IDWidth           = 5
	TitleWidth        = 20
	CategoryWidth     = 12
	PriorityWidth     = 12
	StatusWidth       = 12
	TagsWidth         = 40
	DurationWidth     = 15
	CreatedAtWidth    = 21
	DefaultTableWidth = IDWidth + TitleWidth + CategoryWidth + PriorityWidth + StatusWidth + TagsWidth + DurationWidth + CreatedAtWidth
)

type MinimalStyle struct {
	TermWidth  int
	TermHeight int
}

func (MinimalStyle) RenderMessage(message string) string {
	return lipgloss.NewStyle().Bold(true).Foreground(styles.Main).Render(message)
}

func (MinimalStyle) RenderError(message string) string {
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#d34545ff")).Render(message)
}

func (MinimalStyle) RenderSummary(title string, items []SummaryItem) string {
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

func (ms MinimalStyle) RenderTasks(tasks []models.Task) string {
	cellStyle := lipgloss.NewStyle()
	columnStyles := []lipgloss.Style{
		cellStyle.Width(IDWidth),        // ID
		cellStyle.Width(TitleWidth),     // Title
		cellStyle,                       // Description
		cellStyle.Width(CategoryWidth),  // Category
		cellStyle.Width(PriorityWidth),  // Priority
		cellStyle.Width(StatusWidth),    // Status
		cellStyle.Width(TagsWidth),      // Tags
		cellStyle.Width(DurationWidth),  // Duration
		cellStyle.Width(CreatedAtWidth), // Created At
	}

	headerStyle := lipgloss.NewStyle().Foreground(styles.Main).Bold(true)

	descriptionWidth := max(1, ms.TermWidth-DefaultTableWidth) // Total width minus other column widths and padding

	table := table.New().
		Width(max(1, ms.TermWidth-5)).
		// Border(lipgloss.NormalBorder()).BorderRow(true).BorderColumn(false).BorderLeft(false).BorderRight(false).
		Border(lipgloss.HiddenBorder()).
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
		trimmedTitle := task.Title
		if len(trimmedTitle) > TitleWidth-5 {
			trimmedTitle = trimmedTitle[:TitleWidth-5] + "..."
		}

		trimmedDescription := task.Description
		if len(trimmedDescription) > descriptionWidth-20 {
			trimmedDescription = trimmedDescription[:descriptionWidth-20] + "..."
		}

		trimmedTags := strings.Join(task.Tags, ",")
		if len(trimmedTags) > TagsWidth-5 {
			trimmedTags = trimmedTags[:TagsWidth-5] + "..."
		}

		table.Row(
			strconv.FormatInt(task.ListID, 10),
			trimmedTitle,
			trimmedDescription,
			string(task.Category),
			string(task.Priority),
			string(task.Status),
			trimmedTags,
			task.Duration.Format("15:04:05"),
			task.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	}
	return table.String()
}
