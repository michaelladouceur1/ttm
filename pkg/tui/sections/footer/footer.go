package footer

import (
	"ttm/pkg/styles"
	"ttm/pkg/tui/context"

	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	ctx        *context.TUIContext
	footerText string
}

var (
	style = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		Foreground(styles.Blue).
		Width(100) // set width to match main view
)

func NewModel(ctx *context.TUIContext) Model {
	m := Model{
		ctx:        ctx,
		footerText: "Press q to quit | Press h for help",
	}

	return m
}

func (m Model) View() string {
	return style.Width(m.ctx.Width).Render(m.footerText)
}
