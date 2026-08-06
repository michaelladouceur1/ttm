package ui

import (
	"fmt"
	"time"
	"ttm/pkg/logger"
	"ttm/pkg/models"
	"ttm/pkg/styles"

	"github.com/charmbracelet/lipgloss"
)

func (m *model) addTask(args []string) {
	if len(args) == 0 {
		m.content = "Error: /add requires a title.\n\nUsage: /add \"title\" [\"description\"]"
		return
	}
	if len(args) > 2 {
		m.content = "Error: /add accepts a title and an optional description."
		return
	}

	task := models.Task{
		Title:    args[0],
		Category: models.Category(m.cfg.AddFlags.Category),
		Priority: models.Priority(m.cfg.AddFlags.Priority),
		Status:   models.Status(m.cfg.AddFlags.Status),
		OpenedAt: time.Now(),
	}
	if len(args) == 2 {
		task.Description = args[1]
	}
	if err := task.Validate(); err != nil {
		m.content = "Error adding task: " + err.Error()
		return
	}
	if err := m.store.InsertTask(task); err != nil {
		m.content = "Error adding task: " + err.Error()
		return
	}
	// m.content = fmt.Sprintf("Added task: %s\n\nUse /list to view tasks.", task.Title)
	m.content = lipgloss.NewStyle().Bold(true).Foreground(styles.Highlight1).Render(fmt.Sprintf("Added task: %s\n\nUse /list to view tasks.", task.Title))
}

func (m *model) listTasks(args []string) {
	if len(args) > 1 {
		m.content = "Error: /list accepts at most one search query."
		return
	}

	query := ""
	if len(args) == 1 {
		query = args[0]
	}
	tasks, err := m.store.ListTasks(
		query,
		models.Category(m.cfg.ListFlags.Category),
		models.Status(m.cfg.ListFlags.Status),
		models.Priority(m.cfg.ListFlags.Priority),
	)
	if err != nil {
		m.content = "Error listing tasks: " + err.Error()
		return
	}
	if len(tasks) == 0 {
		m.content = "No tasks found."
		return
	}

	m.content = logger.RenderTasks(tasks)
}
