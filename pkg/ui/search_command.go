package ui

import (
	"strings"
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/models"
	"ttm/pkg/store"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type searchClosedMsg struct {
	content string
}

type searchModel struct {
	input   textinput.Model
	store   *store.Store
	content string
}

func newSearchModel(st *store.Store, width int) searchModel {
	input := textinput.New()
	input.Prompt = "> "
	input.Placeholder = "Search tasks..."
	input.Width = max(1, width-6)
	input.Focus()

	return searchModel{
		input:   input,
		store:   st,
		content: "Type a search query.",
	}
}

func (m searchModel) Update(msg tea.Msg) (childModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.input.Width = max(1, msg.Width-6)
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "esc" {
			return m, func() tea.Msg { return searchClosedMsg{content: m.content} }
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	query := strings.TrimSpace(m.input.Value())
	if query == "" {
		m.content = "Type a search query."
	} else {
		m.content = searchTasks(m.store, query)
	}
	return m, cmd
}

func (m searchModel) InputView() string {
	return m.input.View()
}

func (m searchModel) View() string {
	return m.content
}

func searchTasks(st *store.Store, query string) string {
	search, err := models.ParseTaskSearch(query)
	if err != nil {
		return "Error parsing search query: " + err.Error()
	}

	tasks, err := st.SearchTasks(search)
	if err != nil {
		return "Error searching tasks: " + err.Error()
	}

	if err := fs.UpdateIDMapFile(tasks); err != nil {
		return "Error searching tasks: " + err.Error()
	}

	return logger.RenderTasks(tasks)
}
