package logger

import (
	"fmt"
	"ttm/pkg/styles"

	"github.com/charmbracelet/lipgloss"
)

const (
	Separator       = "─"
	SeparatorMargin = 2
)

var (
	// General styles
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(styles.Blue)
	textStyle   = lipgloss.NewStyle().Foreground(styles.DarkBlue)
	errorStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#d34545ff"))

	treeConnStyle = lipgloss.NewStyle().Foreground(styles.Gray)

	// Table styles
	// Header style
	tableHeaderStyle = lipgloss.NewStyle().Foreground(styles.Blue).Bold(true).Align(lipgloss.Center)

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
	oddRowStyle  = cellStyle.Foreground(styles.Gray)
	evenRowStyle = cellStyle.Foreground(styles.LightGray)
)

func LogMessage(strs ...string) {
	fmt.Println(headerStyle.Render(strs...))
}

func LogError(strs ...any) {
	fmt.Println(errorStyle.Render(fmt.Sprint(strs...)))
}
