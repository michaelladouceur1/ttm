package sessions

import (
	"strconv"
	"time"
	"ttm/pkg/tui/context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	ctx        *context.TUIContext
	HeaderText string
}

var (
	style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		Padding(0, 1)
)

func NewModel(ctx *context.TUIContext) Model {
	m := Model{
		ctx:        ctx,
		HeaderText: "sessions",
	}

	// end := time.Now()
	// start := end.AddDate(0, 0, -5)

	// sessions, err := m.ctx.Store.GetSessionsByTimeRange(start, end)
	// if err != nil {
	// 	m.HeaderText = "Error loading sessions"
	// } else {
	// 	var sessionString string
	// 	for _, session := range sessions {
	// 		sessionString += strconv.Itoa(int(session.TaskId)) + " "
	// 	}
	// 	m.HeaderText = "Sessions Loaded: " + sessionString
	// }

	return m
}

func (m Model) Init() tea.Cmd {
	end := time.Now()
	start := end.AddDate(0, 0, -5)

	sessions, err := m.ctx.Store.GetSessionsByTimeRange(start, end)
	if err != nil {
		m.HeaderText = "Error loading sessions"
	} else {
		var sessionString string
		for _, session := range sessions {
			sessionString += strconv.Itoa(int(session.TaskId)) + " "
		}
		m.HeaderText = "Sessions Loaded: " + sessionString
	}

	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return style.Width(m.ctx.Dims.Sessions.Width).Height(m.ctx.Dims.Sessions.Height).Render(m.HeaderText)
}
