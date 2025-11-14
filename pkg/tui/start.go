package tui

import tea "github.com/charmbracelet/bubbletea"

type StartModel struct{}

func NewStartModel() StartModel {
	return StartModel{}
}

func (m StartModel) Init() tea.Cmd {
	return nil
}

func (m StartModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	if model, cmd := defaultUpdate(m, message); model != nil && cmd != nil {
		return model, cmd
	}

	return m, nil
}

func (m StartModel) View() string {
	return "Welcome to Terminal Todo Manager!\nPress any key to continue..."
}
