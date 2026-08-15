package ui

import (
	"strconv"
	"strings"
	"ttm/pkg/config"
	"ttm/pkg/fs"
	"ttm/pkg/models"
	"ttm/pkg/store"
	"ttm/pkg/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type detailClosedMsg struct {
	content string
}

type detailModel struct {
	input   textinput.Model
	cfg     *config.Config
	store   *store.Store
	listID  string
	content string
}

func newDetailModel(cfg *config.Config, st *store.Store, listID string) detailModel {
	input := textinput.New()
	input.Prompt = "> "

	m := detailModel{
		input:   input,
		cfg:     cfg,
		store:   st,
		listID:  listID,
		content: "Task details for " + listID,
	}

	m.loadTaskDetails()
	return m
}

func (m detailModel) Update(msg tea.Msg) (childModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return detailClosedMsg{content: m.content} }
		}
	}

	return m, func() tea.Msg { return detailClosedMsg{content: m.content} }
}

func (m detailModel) InputView() string {
	return m.input.View()
}

func (m detailModel) View() string {
	return m.content
}

func (m *detailModel) loadTaskDetails() {
	id, err := strconv.Atoi(m.listID)
	if err != nil {
		m.content = "Invalid task ID: " + m.listID
		return
	}

	taskID, err := fs.GetTaskIDFromTempID(int64(id))
	if err != nil {
		m.content = "Task not found: " + m.listID
		return
	}

	task, err := m.store.GetTaskByID(taskID)
	if err != nil {
		m.content = "Task not found: " + m.listID
		return
	}

	var content strings.Builder
	mainTable, width := m.renderMainTable(task)
	content.WriteString(mainTable)
	content.WriteString("\n")
	content.WriteString(m.renderDescription(task, width))
	content.WriteString("\n")

	notes, err := m.store.GetNotesByTaskID(taskID)
	if err != nil {
		m.content = "Failed to load task notes: " + err.Error()
		return
	}

	content.WriteString(m.renderNotes(notes, width))

	m.content = content.String()
}

func (m *detailModel) renderMainTable(task models.Task) (string, int) {
	padding := 5

	dID := strconv.FormatInt(task.ListID, 10)
	dTitle := task.Title
	dPriority := string(task.Priority)
	dStatus := string(task.Status)
	dDuration := task.Duration.Format("15:04")
	dCreatedAt := task.CreatedAt.Format("2006-01-02 15:04")
	dClosedAt := task.ClosedAt.Format("2006-01-02 15:04")

	cellStyle := lipgloss.NewStyle().Foreground(styles.Highlight1)
	columns := []lipgloss.Style{
		cellStyle.Width(max(2+padding, len(dID)+padding)),         // ID
		cellStyle.Width(max(35, len(dTitle)+padding)),             // Title
		cellStyle.Width(max(8+padding, len(dPriority)+padding)),   // Priority
		cellStyle.Width(max(6+padding, len(dStatus)+padding)),     // Status
		cellStyle.Width(max(8+padding, len(dDuration)+padding)),   // Duration
		cellStyle.Width(max(10+padding, len(dCreatedAt)+padding)), // Created At
		cellStyle.Width(max(9+padding, len(dClosedAt)+padding)),   // Closed At
	}

	width := 0
	for _, col := range columns {
		width += col.GetWidth()
	}

	table := table.New().
		Border(lipgloss.HiddenBorder()).
		StyleFunc(func(row, column int) lipgloss.Style {
			if row == 0 {
				return lipgloss.NewStyle().Bold(true).Foreground(styles.Main)
			}
			return columns[column]
		}).
		Row(
			"ID",
			"Title",
			"Priority",
			"Status",
			"Duration",
			"Created At",
			"Closed At",
		).
		Row(
			dID,
			dTitle,
			dPriority,
			dStatus,
			dDuration,
			dCreatedAt,
			dClosedAt,
		)

	return table.String(), width
}

func (m *detailModel) renderDescription(task models.Task, width int) string {
	cellStyle := lipgloss.NewStyle().Foreground(styles.Highlight1)
	columns := []lipgloss.Style{
		cellStyle.PaddingRight(5), // Description
		cellStyle.Width(20),       // Tags
	}

	table := table.New().
		Width(width).
		Border(lipgloss.HiddenBorder()).
		StyleFunc(func(row, column int) lipgloss.Style {
			if row == 0 {
				return lipgloss.NewStyle().Bold(true).Foreground(styles.Main)
			}
			return columns[column]
		}).
		Row("Description", "Tags").
		Row(task.Description, strings.Join(task.Tags, ", "))

	return table.String()
}

func (m *detailModel) renderNotes(notes []models.Note, width int) string {
	heading := lipgloss.NewStyle().Bold(true).Foreground(styles.Main).Render("Notes")
	if len(notes) == 0 {
		return heading
	}

	cellStyle := lipgloss.NewStyle().Foreground(styles.Highlight1)
	columns := []lipgloss.Style{
		cellStyle.Width(21), // Created At
		cellStyle,           // Note
	}

	table := table.New().
		Width(width).
		Border(lipgloss.HiddenBorder()).
		StyleFunc(func(row, column int) lipgloss.Style {
			if row == 0 {
				return lipgloss.NewStyle().Bold(true).Foreground(styles.Main)
			}
			return columns[column]
		}).
		Row("Created At", "Note")

	for _, note := range notes {
		table.Row(note.CreatedAt.Format("2006-01-02 15:04"), note.Content)
	}

	return heading + "\n" + table.String()
}
