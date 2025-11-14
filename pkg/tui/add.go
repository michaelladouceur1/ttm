package tui

import (
	"strings"
	"ttm/pkg/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(styles.Blue).Padding(0, 0, 0, 2)
	focusedStyle = lipgloss.NewStyle().Foreground(styles.Blue).Padding(0, 0, 0, 1)
	noStyle      = lipgloss.NewStyle().Padding(0, 0, 0, 1)
	cursorStyle  = focusedStyle
)

type AddModel struct {
	focusIndex int
	inputs     []textinput.Model
}

func NewAddModel() AddModel {
	m := AddModel{inputs: make([]textinput.Model, 2)}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle

		switch i {
		case 0:
			t.Placeholder = "Title"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.Width = 30
			t.CharLimit = 30
		case 1:
			t.Placeholder = "Description"
			t.Blur()
			t.PromptStyle = noStyle
			t.TextStyle = noStyle
			t.Width = 128
			t.CharLimit = 128

		}

		m.inputs[i] = t
	}

	return m
}

func (m AddModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m AddModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	if model, cmd := defaultUpdate(m, message); model != nil && cmd != nil {
		return model, cmd
	}

	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			var validInputs = true
			for i := range m.inputs {
				if m.inputs[i].Value() == "" {
					validInputs = false
				}
			}
			if validInputs {
				// Submit the form
			}
			fallthrough
		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex >= len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) - 1
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(message)

	return m, cmd
}

func (m AddModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m AddModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Add a New Task"))
	b.WriteRune('\n')
	b.WriteRune('\n')

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		b.WriteRune('\n')
	}

	// b.WriteString(helpStyle.Render("cursor mode is "))
	// b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	// b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}
