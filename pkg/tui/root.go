package tui

import (
	"ttm/pkg/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(0)

var (
	listTitleStyle     = lipgloss.NewStyle().Bold(true).Foreground(styles.Blue)
	selectedTitleStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(styles.Blue).
				Foreground(styles.DarkBlue).
				Padding(0, 0, 0, 1)
)

type RootModel struct {
	list list.Model
}

func NewRootModel() RootModel {
	items := []list.Item{
		NavItem{title: "Add Task", page: string(AddPage)},
		NavItem{title: "Start Session", page: string(StartPage)},
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.Styles.SelectedTitle = selectedTitleStyle

	m := RootModel{list: list.New(items, delegate, 0, 0)}

	m.list.Title = "Terminal Todo Manager"
	m.list.SetShowStatusBar(false)
	m.list.Styles.Title = listTitleStyle

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
