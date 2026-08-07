package ui

import (
	"strconv"
	"ttm/pkg/config"
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/store"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type detailClosedMsg struct {
	content string
}

type detailKeyMap struct {
	close key.Binding
	help  key.Binding
}

func newDetailKeyMap() detailKeyMap {
	return detailKeyMap{
		close: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "close")),
		help:  key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	}
}

func (k detailKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.close, k.help}
}

func (k detailKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.close, k.help}}
}

type detailModel struct {
	input   textinput.Model
	cfg     *config.Config
	store   *store.Store
	listID  string
	content string
	help    help.Model
	keys    detailKeyMap
}

func newDetailModel(cfg *config.Config, st *store.Store, listID string, inputWidth int) detailModel {
	input := textinput.New()
	input.Prompt = "> "
	input.Width = inputWidth

	keys := newDetailKeyMap()
	help := help.New()
	help.Width = max(1, inputWidth)
	m := detailModel{
		input:   input,
		cfg:     cfg,
		store:   st,
		listID:  listID,
		content: "Task details for " + listID,
		help:    help,
		keys:    keys,
	}

	m.loadTaskDetails()
	return m
}

func (m detailModel) Update(msg tea.Msg) (childModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.input.Width = max(1, msg.Width-6)
		m.help.Width = max(1, msg.Width-4)
		return m, nil
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.help) {
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return detailClosedMsg{content: m.content} }
		}
	}
	return m, nil
}

func (m detailModel) InputView() string {
	return m.input.View()
}

func (m detailModel) View() string {
	return m.content + "\n" + m.help.View(m.keys)
}

func (m *detailModel) loadTaskDetails() {
	id, err := strconv.Atoi(m.listID)
	if err != nil {
		m.content = "Invalid task ID: " + m.listID
		return
	}

	taskID, err := fs.GetTaskIDFromTempID(int64(id))
	if err != nil {
		m.content = "Task not found: " + m.listID
		return
	}

	task, err := m.store.GetTaskByID(taskID)
	if err != nil {
		m.content = "Task not found: " + m.listID
		return
	}

	m.content = logger.RenderTaskDetails(task)
}
