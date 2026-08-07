package ui

import (
	"strings"
	"ttm/pkg/config"
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/models"
	"ttm/pkg/store"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type listClosedMsg struct {
	content string
}

type listKeyMap struct {
	search key.Binding
	scroll key.Binding
	close  key.Binding
	help   key.Binding
}

func newListKeyMap() listKeyMap {
	return listKeyMap{
		search: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "search")),
		scroll: key.NewBinding(key.WithKeys("up", "down", "pgup", "pgdown", "home", "end"), key.WithHelp("up/down", "scroll")),
		close:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "close")),
		help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	}
}

func (k listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.search, k.scroll, k.close, k.help}
}

func (k listKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.search, k.scroll, k.close, k.help}}
}

type listModel struct {
	input    textinput.Model
	cfg      *config.Config
	store    *store.Store
	content  string
	viewport viewport.Model
	help     help.Model
	keys     listKeyMap
	width    int
	height   int
}

func newListModel(cfg *config.Config, st *store.Store, width, height int, args []string) listModel {
	input := textinput.New()
	input.Prompt = "> "
	input.Placeholder = "Enter search query..."
	input.Focus()

	keys := newListKeyMap()
	help := help.New()
	m := listModel{
		input: input,
		cfg:   cfg,
		store: st,
		help:  help,
		keys:  keys,
	}
	m.resize(width, height)
	if len(args) == 1 {
		m.input.SetValue(args[0])
	}
	m.listTasks(args)
	return m
}

func (m listModel) Update(msg tea.Msg) (childModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.help) {
			m.help.ShowAll = !m.help.ShowAll
			m.resize(m.width, m.height)
			return m, nil
		}
		switch msg.String() {
		case "enter":
			query := strings.TrimSpace(m.input.Value())
			if query == "" {
				m.listTasks(nil)
			} else {
				m.listTasks([]string{query})
			}
			return m, nil
		case "esc":
			return m, func() tea.Msg { return listClosedMsg{content: m.content} }
		case "up", "down", "pgup", "pgdown", "home", "end":
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m listModel) InputView() string {
	return m.input.View()
}

func (m listModel) View() string {
	return m.viewport.View() + "\n" + m.help.View(m.keys)
}

func (m *listModel) listTasks(args []string) {
	if len(args) > 1 {
		m.setContent("Error: /list accepts at most one search query.")
		return
	}

	query := ""
	if len(args) == 1 {
		query = args[0]
	}
	tasks, err := m.store.ListTasks(
		query,
		models.Category(m.cfg.ListFlags.Category),
		models.Status(m.cfg.ListFlags.Status),
		models.Priority(m.cfg.ListFlags.Priority),
	)
	if err != nil {
		m.setContent("Error listing tasks: " + err.Error())
		return
	}
	if len(tasks) == 0 {
		m.setContent("No tasks found.")
		return
	}

	err = fs.UpdateIDMapFile(tasks)
	if err != nil {
		m.setContent("Error updating ID map file: " + err.Error())
		return
	}

	m.setContent(logger.RenderTasks(tasks))
}

func (m *listModel) setContent(content string) {
	m.content = content
	m.viewport.SetContent(content)
}

func (m *listModel) resize(width, height int) {
	m.width = width
	m.height = height
	m.input.Width = max(1, width-6)
	m.help.Width = max(1, width-4)
	m.viewport.Width = max(1, width-4)
	m.viewport.Height = max(1, height-7-lipgloss.Height(m.help.View(m.keys)))
}
