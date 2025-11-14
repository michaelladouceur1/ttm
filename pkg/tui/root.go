package tui

import (
	"strings"
	"ttm/pkg/styles"
	"ttm/pkg/tui/context"
	"ttm/pkg/tui/sections/footer"
	"ttm/pkg/tui/sections/left"
	"ttm/pkg/tui/sections/middle"
	"ttm/pkg/tui/sections/right"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(0).Border(lipgloss.NormalBorder(), true, true, true, true).Padding(0)

var (
	listTitleStyle     = lipgloss.NewStyle().Bold(true).Foreground(styles.Blue)
	selectedTitleStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(styles.Blue).
				Foreground(styles.DarkBlue).
				Padding(0, 0, 0, 1)
)

type RootModel struct {
	list          list.Model
	ctx           *context.TUIContext
	leftSection   left.Model
	middleSection middle.Model
	rightSection  right.Model
	footer        footer.Model
}

func NewRootModel(ctx *context.TUIContext) RootModel {
	items := []list.Item{
		NavItem{title: "Add Task", page: string(AddPage)},
		NavItem{title: "Start Session", page: string(StartPage)},
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.Styles.SelectedTitle = selectedTitleStyle

	m := RootModel{
		list:          list.New(items, delegate, 0, 0),
		ctx:           ctx,
		footer:        footer.NewModel(ctx),
		leftSection:   left.NewModel(ctx),
		middleSection: middle.NewModel(ctx),
		rightSection:  right.NewModel(ctx),
	}

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
		leftDims, middleDims, rightDims, footerDims := calculateSectionDims(msg.Width, msg.Height)
		m.ctx.TermWidth = msg.Width
		m.ctx.TermHeight = msg.Height
		m.ctx.LeftDims = leftDims
		m.ctx.MiddleDims = middleDims
		m.ctx.RightDims = rightDims
		m.ctx.FooterDims = footerDims
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(message)
	return m, cmd
}

func (m RootModel) View() string {
	s := strings.Builder{}
	// listView := docStyle.Render(m.list.View())
	// textView := lipgloss.NewStyle().Render("Welcome to Terminal Todo Manager!\n\nUse the arrow keys to navigate and press Enter to select an option.")

	// s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, listView, textView))
	// s.WriteString("\n")
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Bottom, m.leftSection.View(), m.middleSection.View(), m.rightSection.View()))
	s.WriteString("\n")
	s.WriteString(m.footer.View())
	return s.String()
}

func (m RootModel) getSelectedPage() TuiPages {
	index := m.list.GlobalIndex()
	selectedItem := m.list.Items()[index].(NavItem)
	return TuiPages(selectedItem.page)
}
