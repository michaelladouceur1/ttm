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
		Border(lipgloss.NormalBorder(), false, true).
		Foreground(styles.Blue)
)

func NewModel(ctx *context.TUIContext) Model {
	m := Model{
		ctx:         ctx,
		sectionText: "Middle Section",
	}

	return m
}

func (m Model) View() string {
	return style.Width(m.ctx.Dims.Middle.Width).Height(m.ctx.Dims.Middle.Height).Render(m.sectionText)
}
