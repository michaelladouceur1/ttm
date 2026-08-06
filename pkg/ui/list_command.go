package ui

import (
	"ttm/pkg/logger"
	"ttm/pkg/models"
)

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
