package header

import (
	"ttm/pkg/styles"
	"ttm/pkg/tui/context"

	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	ctx        *context.TUIContext
	headerText string
}

var (
	style = lipgloss.NewStyle().
		Background(styles.DarkBlue).
		Padding(0, 1)
)

func NewModel(ctx *context.TUIContext) Model {
	m := Model{
		ctx:        ctx,
		headerText: "Welcome to Terminal Todo Manager",
	}

	return m
}

func (m Model) View() string {
	return style.Width(m.ctx.Dims.Header.Width).Height(m.ctx.Dims.Header.Height).Render(m.headerText)
}
