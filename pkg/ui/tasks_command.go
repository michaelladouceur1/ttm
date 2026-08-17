package ui

import (
	"fmt"
	"strings"
	"ttm/pkg/config"
	"ttm/pkg/fs"
	"ttm/pkg/logger"
	"ttm/pkg/models"
	"ttm/pkg/store"

	tea "github.com/charmbracelet/bubbletea"
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

var categoryFilterOptions = []filterOption{
	{label: "Task", value: string(models.CategoryTask)},
	{label: "Meeting", value: string(models.CategoryMeeting)},
}

var statusFilterOptions = []filterOption{
	{label: "Open", value: string(models.StatusOpen)},
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
	content string
}

func newTasksModel(cfg *config.Config, st *store.Store) tasksModel {
	m := tasksModel{
		cfg:   cfg,
		store: st,
		groups: []filterGroup{
			newFilterGroup("Category", categoryFilterOptions, cfg.ListFlags.Category),
			newFilterGroup("Status", statusFilterOptions, cfg.ListFlags.Status),
			newFilterGroup("Priority", priorityFilterOptions, cfg.ListFlags.Priority),
		},
	}
	m.listTasks()
	return m
}

func newFilterGroup(name string, options []filterOption, defaultValue string) filterGroup {
	selected := make(map[string]bool)
	if defaultValue != "" {
		selected[defaultValue] = true
	}
	return filterGroup{name: name, options: options, selected: selected}
}

func (m tasksModel) Update(msg tea.Msg) (childModel, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
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
	return m.renderFilters() + "\n" + m.content
}

func (m tasksModel) renderFilters() string {
	var body strings.Builder
	cursor := 0
	body.WriteString("Filters\n")
	for _, group := range m.groups {
		fmt.Fprintf(&body, "\n%s\n", group.name)
		for _, option := range group.options {
			prefix := "  "
			if cursor == m.cursor {
				prefix = "> "
			}
			checkbox := "[ ]"
			if group.selected[option.value] {
				checkbox = "[x]"
			}
			fmt.Fprintf(&body, "%s%s %s\n", prefix, checkbox, option.label)
			cursor++
		}
	}
	return body.String()
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

func (m *tasksModel) matchesFilters(task models.Task) bool {
	return groupMatches(m.groups[0], string(task.Category)) &&
		groupMatches(m.groups[1], string(task.Status)) &&
		groupMatches(m.groups[2], string(task.Priority))
}

func groupMatches(group filterGroup, value string) bool {
	if len(group.selected) == 0 {
		return true
	}
	return group.selected[value]
}

func (m *tasksModel) listTasks() {
	tasks, err := m.store.ListTasks("", "", "", "")
	if err != nil {
		m.setContent("Error listing tasks: " + err.Error())
		return
	}

	filtered := make([]models.Task, 0, len(tasks))
	for _, task := range tasks {
		if m.matchesFilters(task) {
			filtered = append(filtered, task)
		}
	}

	if len(filtered) == 0 {
		m.setContent("No tasks found.")
		return
	}

	if err := fs.UpdateIDMapFile(filtered); err != nil {
		m.setContent("Error updating ID map file: " + err.Error())
		return
	}

	m.setContent(logger.RenderTasks(filtered))
}

func (m *tasksModel) setContent(content string) {
	m.content = content
}
