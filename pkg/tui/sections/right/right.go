package right

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
		Foreground(styles.Blue)
)

func NewModel(ctx *context.TUIContext) Model {
	m := Model{
		ctx:         ctx,
		sectionText: "Right Section",
	}

	return m
}

func (m Model) View() string {
	return style.Width(m.ctx.Dims.Right.Width).Height(m.ctx.Dims.Right.Height).Render(m.sectionText)
}
