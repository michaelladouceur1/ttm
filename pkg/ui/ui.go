package ui

import (
	"fmt"
	"strings"
	"ttm/pkg/config"
	"ttm/pkg/store"
	"ttm/pkg/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type command struct {
	name        string
	description string
}

var commands = []command{
	{name: "add", description: "Add a task with a guided form"},
	{name: "list", description: "List open tasks, optionally matching a query"},
	{name: "detail", description: "Show details for a task, including notes and sessions"},
}

// Run starts the interactive terminal UI.
func Run(cfg *config.Config, st *store.Store) error {
	program := tea.NewProgram(newModel(cfg, st), tea.WithAltScreen())
	_, err := program.Run()
	return err
}

type model struct {
	input       textinput.Model
	cfg         *config.Config
	store       *store.Store
	width       int
	height      int
	selected    int
	suggestions []command
	content     string
	active      childModel
}

type childModel interface {
	Update(tea.Msg) (childModel, tea.Cmd)
	InputView() string
	View() string
}

func newModel(cfg *config.Config, st *store.Store) model {
	input := textinput.New()
	input.Prompt = "> "
	input.Placeholder = "Type / for commands"
	input.Focus()

	return model{
		input:   input,
		cfg:     cfg,
		store:   st,
		width:   80,
		height:  24,
		content: "Type / to see available commands.\n\nUse Ctrl+C to exit.",
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case addCompleteMsg:
		m.active = nil
		m.content = fmt.Sprintf("Added task: %s\n\nUse /list to view tasks.", msg.title)
		return m, nil
	case addCancelledMsg:
		m.active = nil
		m.content = "Task creation cancelled."
		return m, nil
	case listClosedMsg:
		m.active = nil
		m.content = msg.content
		return m, nil
	case detailClosedMsg:
		m.active = nil
		m.content = msg.content
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = max(1, msg.Width-6)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	if m.active != nil {
		var cmd tea.Cmd
		m.active, cmd = m.active.Update(msg)
		return m, cmd
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "up":
			if len(m.suggestions) > 0 {
				m.selected = (m.selected - 1 + len(m.suggestions)) % len(m.suggestions)
				return m, nil
			}
		case "down":
			if len(m.suggestions) > 0 {
				m.selected = (m.selected + 1) % len(m.suggestions)
				return m, nil
			}
		case "tab":
			if len(m.suggestions) > 0 {
				m.completeSuggestion()
				return m, nil
			}
		case "enter":
			if m.input.Value() == "/" && len(m.suggestions) > 0 {
				m.completeSuggestion()
				return m, nil
			}
			m.execute()
			return m, nil
		case "esc":
			m.input.SetValue("")
			m.suggestions = nil
			m.selected = 0
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.updateSuggestions()
	return m, cmd
}

func (m *model) updateSuggestions() {
	value := strings.TrimPrefix(m.input.Value(), "/")
	if !strings.HasPrefix(m.input.Value(), "/") || strings.ContainsAny(value, " \t") {
		m.suggestions = nil
		m.selected = 0
		return
	}

	m.suggestions = m.suggestions[:0]
	for _, command := range commands {
		if strings.HasPrefix(command.name, strings.ToLower(value)) {
			m.suggestions = append(m.suggestions, command)
		}
	}
	if m.selected >= len(m.suggestions) {
		m.selected = 0
	}
}

func (m *model) completeSuggestion() {
	m.input.SetValue("/" + m.suggestions[m.selected].name + " ")
	m.input.CursorEnd()
	m.suggestions = nil
	m.selected = 0
}

func (m *model) execute() {
	args, err := parseCommand(m.input.Value())
	if err != nil {
		m.content = "Error: " + err.Error()
		return
	}
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "add":
		if len(args) != 1 {
			m.content = "Usage: /add\n\nTasks are created one field at a time."
			break
		}
		m.active = newAddModel(m.cfg, m.store, m.input.Width)
	case "list":
		m.active = newListModel(m.cfg, m.store, m.input.Width, args[1:])
	case "detail":
		if len(args) != 2 {
			m.content = "Usage: /detail <task_id>"
			break
		}
		m.active = newDetailModel(m.cfg, m.store, args[1])
	default:
		m.content = fmt.Sprintf("Unknown command: /%s\n\nType / to see available commands.", args[0])
	}
	m.input.SetValue("")
	m.suggestions = nil
	m.selected = 0
}

func (m model) View() string {
	inputView := m.input.View()
	bodyContent := m.content
	if m.active != nil {
		inputView = m.active.InputView()
		bodyContent = m.active.View()
	}

	input := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Main).
		Width(max(1, m.width-2)).
		Render(inputView)

	var body strings.Builder
	if m.active == nil && len(m.suggestions) > 0 {
		body.WriteString("Commands\n")
		for i, command := range m.suggestions {
			prefix := "  "
			if i == m.selected {
				prefix = "> "
			}
			fmt.Fprintf(&body, "%s/%-5s %s\n", prefix, command.name, command.description)
		}
		body.WriteString("\n")
	}
	body.WriteString(bodyContent)

	content := lipgloss.NewStyle().
		Width(max(1, m.width-2)).
		Height(max(1, m.height-4)).
		Padding(1, 1).
		Render(body.String())

	return input + "\n" + content
}

func parseCommand(input string) ([]string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, nil
	}
	if !strings.HasPrefix(input, "/") {
		return nil, fmt.Errorf("commands must start with /")
	}

	var args []string
	var current strings.Builder
	var quote rune
	escaped := false
	for _, char := range input[1:] {
		switch {
		case escaped:
			current.WriteRune(char)
			escaped = false
		case char == '\\' && quote != 0:
			escaped = true
		case quote != 0:
			if char == quote {
				quote = 0
			} else {
				current.WriteRune(char)
			}
		case char == '"' || char == '\'':
			quote = char
		case char == ' ' || char == '\t':
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(char)
		}
	}
	if escaped || quote != 0 {
		return nil, fmt.Errorf("unterminated quoted argument")
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args, nil
}
