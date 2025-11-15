package left

import (
	"strconv"
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
		sectionText: strconv.Itoa(ctx.TermWidth) + "x" + strconv.Itoa(ctx.TermHeight),
	}

	return m
}

func (m Model) View() string {
	return style.Width(m.ctx.Dims.Left.Width).Height(m.ctx.Dims.Left.Height).Render(m.sectionText)
}
