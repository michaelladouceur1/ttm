package footer

import (
	"ttm/pkg/tui/context"

	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	ctx        *context.TUIContext
	footerText string
}

var (
	style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		AlignVertical(lipgloss.Center).
		Padding(0, 1)
)

func NewModel(ctx *context.TUIContext) Model {
	m := Model{
		ctx:        ctx,
		footerText: "Press q to quit | Press h for help",
	}

	return m
}

func (m Model) View() string {
	return style.Width(m.ctx.Dims.Footer.Width).Height(m.ctx.Dims.Footer.Height).Render(m.footerText)
}
