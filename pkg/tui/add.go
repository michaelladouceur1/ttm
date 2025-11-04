package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type AddModel struct {
	str string
}

func NewAddModel() AddModel {
	return AddModel{str: ""}
}

func (m AddModel) Init() tea.Cmd {
	return nil
}

func (m AddModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	if model, cmd := defaultUpdate(m, message); model != nil && cmd != nil {
		return model, cmd
	}

	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "f":
			m.str += "f"
		}
	}

	return m, nil
}

func (m AddModel) View() string {
	return "Add Task View: " + m.str
}
