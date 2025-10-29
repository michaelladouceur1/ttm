package logger

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	Separator       = "─"
	SeparatorMargin = 2
)

var (
	// Colors
	blue      = lipgloss.Color("#2cd5fb")
	darkBlue  = lipgloss.Color("#219cb8ff")
	gray      = lipgloss.Color("#6b6b6bff")
	lightGray = lipgloss.Color("#e4e4e4ff")

	// General styles
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(blue)
	textStyle   = lipgloss.NewStyle().Foreground(darkBlue)
	errorStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#d34545ff"))

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

func LogMessage(strs ...string) {
	fmt.Println(headerStyle.Render(strs...))
}

func LogError(strs ...any) {
	fmt.Println(errorStyle.Render(fmt.Sprint(strs...)))
}
