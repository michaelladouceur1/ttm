package ui

import (
	"fmt"
	"strings"
	"time"
	"ttm/pkg/models"
)

func (m *model) startAddForm() {
	m.addStep = addStepTitle
	m.draft = models.Task{}
	m.priority = priorityIndex(models.Priority(m.cfg.AddFlags.Priority))
	m.input.SetValue("")
	m.input.Placeholder = "Enter title..."
	m.content = "Create task\n\nTitle"
}

func (m *model) submitAddField() {
	value := strings.TrimSpace(m.input.Value())

	switch m.addStep {
	case addStepTitle:
		if value == "" {
			m.content = "A title is required."
			return
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
			return
		}
		m.draft.Priority = priority
		m.nextAddStep(addStepTags, "Enter tags (comma-separated)...", "Tags")
	case addStepTags:
		m.draft.Tags = parseTags(value)
		m.saveTask()
	}
}

func (m *model) nextAddStep(step addStep, placeholder, field string) {
	m.addStep = step
	m.input.SetValue("")
	m.input.Placeholder = placeholder
	m.content = "Create task\n\n" + field
}

func (m *model) saveTask() {
	m.draft.Category = models.Category(m.cfg.AddFlags.Category)
	m.draft.Status = models.Status(m.cfg.AddFlags.Status)
	m.draft.OpenedAt = time.Now()

	if err := m.draft.Validate(); err != nil {
		m.content = "Error adding task: " + err.Error()
		return
	}
	if err := m.store.InsertTask(m.draft); err != nil {
		m.content = "Error adding task: " + err.Error()
		return
	}
	m.resetAddForm(fmt.Sprintf("Added task: %s\n\nUse /list to view tasks.", m.draft.Title))
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

func (m *model) resetAddForm(content string) {
	m.addStep = addStepNone
	m.draft = models.Task{}
	m.input.SetValue("")
	m.input.Placeholder = "Type / for commands"
	m.suggestions = nil
	m.selected = 0
	m.content = content
}
