package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(0)

type RootModel struct {
	list list.Model
}

func NewRootModel() RootModel {
	items := []list.Item{
		NavItem{title: "Add Task", page: string(AddPage)},
		NavItem{title: "Start Session", page: string(StartPage)},
	}

	m := RootModel{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Terminal Todo Manager"
	m.list.SetShowStatusBar(false)

	return m
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	if model, cmd := defaultUpdate(m, message); model != nil && cmd != nil {
		return model, cmd
	}

	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return tui.switchPage(m.getSelectedPage())
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(message)
	return m, cmd
}

func (m RootModel) View() string {
	return docStyle.Render(m.list.View())
}

func (m RootModel) getSelectedPage() TuiPages {
	index := m.list.GlobalIndex()
	selectedItem := m.list.Items()[index].(NavItem)
	return TuiPages(selectedItem.page)
}
