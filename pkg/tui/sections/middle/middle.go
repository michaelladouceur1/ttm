package middle

import (
	"ttm/pkg/styles"
	"ttm/pkg/tui/context"

	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	ctx         *context.TUIContext
	sectionText string
}

var (
	style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		Foreground(styles.Blue).
		Width(100) // set width to match main view
)

func NewModel(ctx *context.TUIContext) Model {
	m := Model{
		ctx:         ctx,
		sectionText: "Middle Section",
	}

	return m
}

func (m Model) View() string {
	return style.Width(m.ctx.MiddleDims.Width).Height(m.ctx.MiddleDims.Height).Render(m.sectionText)
}
