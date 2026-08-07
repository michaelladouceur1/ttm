package ui

import (
	"strconv"
	"ttm/pkg/config"
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/store"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type detailClosedMsg struct {
	content string
}

type detailModel struct {
	input   textinput.Model
	cfg     *config.Config
	store   *store.Store
	listID  string
	content string
}

func newDetailModel(cfg *config.Config, st *store.Store, listID string) detailModel {
	input := textinput.New()
	input.Prompt = "> "

	m := detailModel{
		input:   input,
		cfg:     cfg,
		store:   st,
		listID:  listID,
		content: "Task details for " + listID,
	}

	m.loadTaskDetails()
	return m
}

func (m detailModel) Update(msg tea.Msg) (childModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return detailClosedMsg{content: m.content} }
		}
	}

	return m, func() tea.Msg { return detailClosedMsg{content: m.content} } // Return a close message immediately to reset back to main command view
}

func (m detailModel) InputView() string {
	return m.input.View()
}

func (m detailModel) View() string {
	return m.content
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
