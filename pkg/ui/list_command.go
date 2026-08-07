package ui

import (
	"strings"
	"ttm/pkg/config"
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/models"
	"ttm/pkg/store"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type listClosedMsg struct {
	content string
}

type listModel struct {
	input   textinput.Model
	cfg     *config.Config
	store   *store.Store
	content string
}

func newListModel(cfg *config.Config, st *store.Store, inputWidth int, args []string) listModel {
	input := textinput.New()
	input.Prompt = "> "
	input.Placeholder = "Enter search query..."
	input.Width = inputWidth
	input.Focus()

	m := listModel{
		input: input,
		cfg:   cfg,
		store: st,
	}
	if len(args) == 1 {
		m.input.SetValue(args[0])
	}
	m.listTasks(args)
	return m
}

func (m listModel) Update(msg tea.Msg) (childModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.input.Width = max(1, msg.Width-6)
		return m, nil
	case tea.KeyMsg:
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
	return m.content
}

func (m *listModel) listTasks(args []string) {
	if len(args) > 1 {
		m.content = "Error: /list accepts at most one search query."
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
		m.content = "Error listing tasks: " + err.Error()
		return
	}
	if len(tasks) == 0 {
		m.content = "No tasks found."
		return
	}

	err = fs.UpdateIDMapFile(tasks)
	if err != nil {
		m.content = "Error updating ID map file: " + err.Error()
		return
	}

	m.content = logger.RenderTasks(tasks)
}
