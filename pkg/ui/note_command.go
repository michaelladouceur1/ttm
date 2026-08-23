package ui

import (
	"strconv"
	"ttm/pkg/config"
	"ttm/pkg/fs"
	"ttm/pkg/store"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type noteClosedMsg struct {
	content string
}

type noteModel struct {
	input   textinput.Model
	cfg     *config.Config
	store   *store.Store
	listID  string
	taskID  int64
	content string
}

func newNoteModel(cfg *config.Config, st *store.Store, width int, listID string) noteModel {
	input := textinput.New()
	input.Prompt = "> "
	input.Placeholder = "Enter note content..."
	input.Width = max(1, width-6)
	input.Focus()

	m := noteModel{
		input:   input,
		cfg:     cfg,
		store:   st,
		listID:  listID,
		content: "Task notes for " + listID,
	}

	// m.loadTaskDetails()
	return m
}

func newActiveNoteModel(cfg *config.Config, st *store.Store, width int, taskID int64) noteModel {
	m := newNoteModel(cfg, st, width, strconv.FormatInt(taskID, 10))
	m.taskID = taskID
	return m
}

func (m noteModel) Update(msg tea.Msg) (childModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			content := m.input.Value()
			m.saveNote(content)
			return m, func() tea.Msg { return noteClosedMsg{content: m.content} }
		case "esc":
			return m, func() tea.Msg { return noteClosedMsg{content: m.content} }
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m noteModel) InputView() string {
	return m.input.View()
}

func (m noteModel) View() string {
	return m.content
}

func (m *noteModel) saveNote(content string) {
	taskID := m.taskID
	if taskID == 0 {
		id, err := strconv.Atoi(m.listID)
		if err != nil {
			m.content = "Invalid task ID: " + m.listID
			return
		}

		taskID, err = fs.GetTaskIDFromTempID(int64(id))
		if err != nil {
			m.content = "Task not found: " + m.listID
			return
		}
	}

	note, err := m.store.InsertNote(taskID, content)

	if err != nil {
		m.content = "Failed to save note: " + err.Error()
		return
	}

	m.content = "Note saved successfully to task " + m.listID + ": " + note.Content
}
