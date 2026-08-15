package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"ttm/pkg/config"
	"ttm/pkg/store"
	"ttm/pkg/styles"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type command struct {
	name        string
	description string
}

var commands = []command{
	{name: "add", description: "Add a task with a guided form"},
	{name: "cancel", description: "Discard the active session"},
	{name: "tasks", description: "List open tasks, optionally matching a query"},
	{name: "tags", description: "List all tags and their task counts"},
	{name: "start", description: "Start a session for a task"},
	{name: "end", description: "End and save the active session"},
	{name: "detail", description: "Show details for a task, including notes and sessions"},
	{name: "note", description: "Add a task note"},
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
	viewport    viewport.Model
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

	m := model{
		input:    input,
		cfg:      cfg,
		store:    st,
		width:    80,
		height:   24,
		content:  "Type / to see available commands.\n\nUse Ctrl+C to exit.",
		viewport: viewport.New(76, 18),
	}
	m.setViewportContent()
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case addCompleteMsg:
		m.active = nil
		m.content = fmt.Sprintf("Added task: %s\n\nUse /tasks to view tasks.", msg.title)
		m.setViewportContent()
		return m, nil
	case addCancelledMsg:
		m.active = nil
		m.content = "Task creation cancelled."
		m.setViewportContent()
		return m, nil
	case tasksClosedMsg:
		m.active = nil
		m.content = msg.content
		m.setViewportContent()
		return m, nil
	case detailClosedMsg:
		m.active = nil
		m.content = msg.content
		m.setViewportContent()
		return m, nil
	case noteClosedMsg:
		m.active = nil
		m.content = msg.content
		m.setViewportContent()
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = max(1, msg.Width-6)
		m.viewport.Width = max(1, msg.Width-4)
		m.viewport.Height = max(1, msg.Height-6)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	if m.active != nil {
		var cmd tea.Cmd
		m.active, cmd = m.active.Update(msg)
		m.setViewportContent()
		if isViewportNavigation(msg) {
			m.viewport, _ = m.viewport.Update(msg)
		}
		return m, cmd
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "up":
			if len(m.suggestions) > 0 {
				m.selected = (m.selected - 1 + len(m.suggestions)) % len(m.suggestions)
				m.setViewportContent()
				return m, nil
			}
		case "down":
			if len(m.suggestions) > 0 {
				m.selected = (m.selected + 1) % len(m.suggestions)
				m.setViewportContent()
				return m, nil
			}
		case "tab":
			if len(m.suggestions) > 0 {
				m.completeSuggestion()
				m.setViewportContent()
				return m, nil
			}
		case "enter":
			if m.input.Value() == "/" && len(m.suggestions) > 0 {
				m.completeSuggestion()
				return m, nil
			}
			m.execute()
			m.setViewportContent()
			return m, nil
		case "esc":
			m.input.SetValue("")
			m.suggestions = nil
			m.selected = 0
			m.setViewportContent()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.updateSuggestions()
	m.setViewportContent()
	if isViewportNavigation(msg) {
		m.viewport, _ = m.viewport.Update(msg)
	}
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

	orderedCommands := make([]command, len(commands))
	copy(orderedCommands, commands)
	sort.Slice(orderedCommands, func(i, j int) bool {
		return orderedCommands[i].name < orderedCommands[j].name
	})

	for _, command := range orderedCommands {
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
	case "tasks":
		m.active = newTasksModel(m.cfg, m.store, m.width, args[1:])
	case "tags":
		if len(args) != 1 {
			m.content = "Usage: /tags"
			break
		}
		m.content = listTags(m.store)
	case "start":
		if len(args) != 2 {
			m.content = "Usage: /start <task_id>"
			break
		}
		m.content = startSession(m.store, args[1], time.Now())
	case "end":
		if len(args) != 1 {
			m.content = "Usage: /end"
			break
		}
		m.content = endSession(m.store, time.Now())
	case "cancel":
		if len(args) != 1 {
			m.content = "Usage: /cancel"
			break
		}
		m.content = cancelSession()
	case "detail":
		if len(args) != 2 {
			m.content = "Usage: /detail <task_id>"
			break
		}
		m.active = newDetailModel(m.cfg, m.store, args[1])
	case "note":
		if len(args) != 2 {
			m.content = "Usage: /note <task_id>"
			break
		}
		m.active = newNoteModel(m.cfg, m.store, m.input.Width, args[1])
	default:
		m.content = fmt.Sprintf("Unknown command: /%s\n\nType / to see available commands.", args[0])
	}
	m.input.SetValue("")
	m.suggestions = nil
	m.selected = 0
}

func (m model) View() string {
	inputView := m.input.View()
	if m.active != nil {
		inputView = m.active.InputView()
	}

	input := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Main).
		Width(max(1, m.width-2)).
		Render(inputView)

	content := lipgloss.NewStyle().
		Width(max(1, m.width-2)).
		Padding(1, 1).
		Render(m.viewport.View())

	return input + "\n" + content
}

func (m *model) setViewportContent() {
	var body strings.Builder
	if m.active == nil && len(m.suggestions) > 0 {
		longestSuggestion := 0
		for _, command := range m.suggestions {
			if len(command.name) > longestSuggestion {
				longestSuggestion = len(command.name)
			}
		}
		body.WriteString("Commands\n")
		for i, command := range m.suggestions {
			prefix := "  "
			if i == m.selected {
				prefix = "> "
			}
			fmt.Fprintf(&body, "%s/%-*s   %s\n", prefix, longestSuggestion, command.name, command.description)
		}
		body.WriteString("\n")
	}
	if m.active != nil {
		body.WriteString(m.active.View())
	} else {
		body.WriteString(m.content)
	}
	m.viewport.SetContent(body.String())
}

func isViewportNavigation(msg tea.Msg) bool {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return false
	}
	switch key.String() {
	case "up", "down", "pgup", "pgdown", "home", "end":
		return true
	default:
		return false
	}
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
