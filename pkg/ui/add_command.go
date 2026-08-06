package ui

import (
	"fmt"
	"strings"
	"time"
	"ttm/pkg/config"
	"ttm/pkg/models"
	"ttm/pkg/store"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type addStep int

const (
	addStepTitle addStep = iota
	addStepDescription
	addStepPriority
	addStepTags
)

var priorities = []models.Priority{
	models.PriorityLow,
	models.PriorityMedium,
	models.PriorityHigh,
}

type addCompleteMsg struct {
	title string
}

type addCancelledMsg struct{}

type addModel struct {
	input    textinput.Model
	cfg      *config.Config
	store    *store.Store
	step     addStep
	draft    models.Task
	priority int
	content  string
}

func newAddModel(cfg *config.Config, st *store.Store, inputWidth int) addModel {
	input := textinput.New()
	input.Prompt = "> "
	input.Placeholder = "Enter title..."
	input.Width = inputWidth
	input.Focus()

	return addModel{
		input:    input,
		cfg:      cfg,
		store:    st,
		step:     addStepTitle,
		priority: priorityIndex(models.Priority(cfg.AddFlags.Priority)),
		content:  "Create task\n\nTitle",
	}
}

func (m addModel) Update(msg tea.Msg) (childModel, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up":
			if m.step == addStepPriority {
				m.priority = (m.priority - 1 + len(priorities)) % len(priorities)
				m.input.SetValue(string(priorities[m.priority]))
				return m, nil
			}
		case "down":
			if m.step == addStepPriority {
				m.priority = (m.priority + 1) % len(priorities)
				m.input.SetValue(string(priorities[m.priority]))
				return m, nil
			}
		case "enter":
			return m.submitAddField()
		case "esc":
			return m, func() tea.Msg { return addCancelledMsg{} }
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m addModel) InputView() string {
	return m.input.View()
}

func (m addModel) View() string {
	var body strings.Builder
	if m.step == addStepPriority {
		body.WriteString("Priority\n")
		for i, priority := range priorities {
			prefix := "  "
			if i == m.priority {
				prefix = "> "
			}
			fmt.Fprintf(&body, "%s%s\n", prefix, priority)
		}
		body.WriteString("\n")
	}
	body.WriteString(m.content)
	return body.String()
}

func (m addModel) submitAddField() (childModel, tea.Cmd) {
	value := strings.TrimSpace(m.input.Value())

	switch m.step {
	case addStepTitle:
		if value == "" {
			m.content = "A title is required."
			return m, nil
		}
		m.draft.Title = value
		m.nextAddStep(addStepDescription, "Enter description... (optional)", "Description")
	case addStepDescription:
		m.draft.Description = value
		m.nextAddStep(addStepPriority, "Choose priority: low, medium, high", "Priority")
	case addStepPriority:
		priority := models.Priority(strings.ToLower(value))
		if value == "" {
			priority = priorities[m.priority]
		}
		if priority != models.PriorityLow && priority != models.PriorityMedium && priority != models.PriorityHigh {
			m.content = "Choose a priority: low, medium, or high."
			return m, nil
		}
		m.draft.Priority = priority
		m.nextAddStep(addStepTags, "Enter tags (comma-separated)...", "Tags")
	case addStepTags:
		m.draft.Tags = parseTags(value)
		return m.saveTask()
	}
	return m, nil
}

func (m *addModel) nextAddStep(step addStep, placeholder, field string) {
	m.step = step
	m.input.SetValue("")
	m.input.Placeholder = placeholder
	m.content = "Create task\n\n" + field
}

func (m addModel) saveTask() (childModel, tea.Cmd) {
	m.draft.Category = models.Category(m.cfg.AddFlags.Category)
	m.draft.Status = models.Status(m.cfg.AddFlags.Status)
	m.draft.OpenedAt = time.Now()

	if err := m.draft.Validate(); err != nil {
		m.content = "Error adding task: " + err.Error()
		return m, nil
	}
	if err := m.store.InsertTask(m.draft); err != nil {
		m.content = "Error adding task: " + err.Error()
		return m, nil
	}
	return m, func() tea.Msg { return addCompleteMsg{title: m.draft.Title} }
}

func parseTags(value string) []string {
	var tags []string
	for _, tag := range strings.Split(value, ",") {
		if tag = strings.TrimSpace(tag); tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags
}

func priorityIndex(priority models.Priority) int {
	for i, option := range priorities {
		if priority == option {
			return i
		}
	}
	return 2
}
