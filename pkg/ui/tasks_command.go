package ui

import (
	"strings"
	"ttm/pkg/config"
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/models"
	"ttm/pkg/store"
	"ttm/pkg/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tasksClosedMsg struct {
	content string
}

type filterOption struct {
	label string
	value string
}

type filterGroup struct {
	name     string
	options  []filterOption
	selected map[string]bool
}

var statusFilterOptions = []filterOption{
	{label: "Open", value: string(models.StatusOpen)},
	{label: "Standby", value: string(models.StatusStandby)},
	{label: "Closed", value: string(models.StatusClosed)},
}

var priorityFilterOptions = []filterOption{
	{label: "Low", value: string(models.PriorityLow)},
	{label: "Medium", value: string(models.PriorityMedium)},
	{label: "High", value: string(models.PriorityHigh)},
}

type tasksModel struct {
	cfg     *config.Config
	store   *store.Store
	groups  []filterGroup
	cursor  int
	width   int
	content string
}

func newTasksModel(cfg *config.Config, st *store.Store, width int) tasksModel {
	m := tasksModel{
		cfg:   cfg,
		store: st,
		width: width,
		groups: []filterGroup{
			newFilterGroup("Status", statusFilterOptions, cfg.ListFlags.Status...),
			newFilterGroup("Priority", priorityFilterOptions, cfg.ListFlags.Priority...),
		},
	}
	m.listTasks()
	return m
}

func newFilterGroup(name string, options []filterOption, defaultValues ...string) filterGroup {
	selected := make(map[string]bool)
	for _, value := range defaultValues {
		if value != "" {
			selected[value] = true
		}
	}
	return filterGroup{name: name, options: options, selected: selected}
}

func (m tasksModel) Update(msg tea.Msg) (childModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case tea.KeyMsg:
		key := msg
		switch key.String() {
		case "up":
			total := m.totalOptions()
			m.cursor = (m.cursor - 1 + total) % total
			return m, nil
		case "down":
			total := m.totalOptions()
			m.cursor = (m.cursor + 1) % total
			return m, nil
		case " ", "enter":
			m.toggleCurrent()
			return m, nil
		case "esc":
			return m, func() tea.Msg { return tasksClosedMsg{content: m.content} }
		}
	}

	return m, nil
}

func (m tasksModel) InputView() string {
	return "↑/↓ move   space/enter toggle filter   esc close"
}

func (m tasksModel) View() string {
	var content strings.Builder

	const panelGap = "     "
	content.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.renderFilters(),
		panelGap,
		m.content,
	))
	return content.String()
}

func (m tasksModel) renderFilters() string {
	var body strings.Builder
	cursor := 0
	width := 0
	for i, group := range m.groups {
		label := group.name + "\n"
		body.WriteString(label)
		if len(label) > width {
			width = len(label)
		}
		for _, option := range group.options {
			prefix := "  "
			if cursor == m.cursor {
				prefix = styles.SelectorPrefix
			}
			checkbox := "[ ]"
			if group.selected[option.value] {
				checkbox = "[x]"
			}
			optionText := prefix + checkbox + " " + option.label
			body.WriteString(optionText)
			body.WriteByte('\n')
			cursor++
			if len(optionText) > width {
				width = len(optionText)
			}
		}
		if i < len(m.groups)-1 {
			body.WriteByte('\n')
		}
	}

	padding := 1
	filters := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Main).
		Width(max(1, width+(2*padding))).
		Padding(0, padding).
		Render(body.String())

	return filters
}

func (m *tasksModel) totalOptions() int {
	total := 0
	for _, group := range m.groups {
		total += len(group.options)
	}
	return total
}

func (m *tasksModel) optionAt(index int) (*filterGroup, filterOption) {
	for i := range m.groups {
		group := &m.groups[i]
		if index < len(group.options) {
			return group, group.options[index]
		}
		index -= len(group.options)
	}
	return nil, filterOption{}
}

func (m *tasksModel) toggleCurrent() {
	group, option := m.optionAt(m.cursor)
	if group == nil {
		return
	}
	group.selected[option.value] = !group.selected[option.value]
	m.listTasks()
}

func (m *tasksModel) listTasks() {
	statuses := m.selectedStatuses(m.groups[m.groupIndex("Status")])
	priorities := m.selectedPriorities(m.groups[m.groupIndex("Priority")])

	tasks, err := m.store.ListTasks(
		statuses,
		priorities,
	)
	if err != nil {
		m.setContent("Error listing tasks: " + err.Error())
		return
	}

	if err := fs.UpdateIDMapFile(tasks); err != nil {
		m.setContent("Error updating ID map file: " + err.Error())
		return
	}

	m.setContent(logger.RenderTasks(tasks))
}

func (m *tasksModel) setContent(content string) {
	m.content = content
}

func (m *tasksModel) groupIndex(name string) int {
	for i, group := range m.groups {
		if group.name == name {
			return i
		}
	}
	return -1
}

func (m *tasksModel) selectedStatuses(group filterGroup) []models.Status {
	values := make([]models.Status, 0, len(group.selected))
	for value, selected := range group.selected {
		if selected {
			values = append(values, models.Status(value))
		}
	}
	return values
}

func (m *tasksModel) selectedPriorities(group filterGroup) []models.Priority {
	values := make([]models.Priority, 0, len(group.selected))
	for value, selected := range group.selected {
		if selected {
			values = append(values, models.Priority(value))
		}
	}
	return values
}
