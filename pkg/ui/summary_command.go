package ui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"ttm/pkg/models"
	"ttm/pkg/store"
	"ttm/pkg/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const defaultSummaryDays = 7

type summaryClosedMsg struct {
	content string
}

type summarySession struct {
	session models.Session
	task    models.Task
}

type summaryTask struct {
	task     models.Task
	duration time.Duration
}

type summaryModel struct {
	store   *store.Store
	width   int
	days    int
	now     time.Time
	content string
}

func newSummaryModel(st *store.Store, width, days int) summaryModel {
	m := summaryModel{
		store: st,
		width: width,
		days:  days,
		now:   time.Now(),
	}
	m.loadSummary()
	return m
}

func summaryDays(args []string) (int, error) {
	if len(args) == 0 {
		return defaultSummaryDays, nil
	}
	if len(args) != 1 {
		return 0, fmt.Errorf("usage: /summary [TIME_PERIOD_IN_DAYS]")
	}

	days, err := strconv.Atoi(args[0])
	if err != nil || days < 1 {
		return 0, fmt.Errorf("TIME_PERIOD_IN_DAYS must be a positive whole number")
	}
	return days, nil
}

func (m summaryModel) Update(msg tea.Msg) (childModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.loadSummary()
	case tea.KeyMsg:
		if msg.String() == "esc" {
			return m, func() tea.Msg { return summaryClosedMsg{content: m.content} }
		}
	}
	return m, nil
}

func (m summaryModel) InputView() string {
	return "esc close"
}

func (m summaryModel) View() string {
	return m.content
}

func (m *summaryModel) loadSummary() {
	start := startOfDay(m.now).AddDate(0, 0, -(m.days - 1))
	end := startOfDay(m.now).AddDate(0, 0, 1).Add(-time.Nanosecond)
	sessions, err := m.store.GetSessionsByTimeRange(start, end)
	if err != nil {
		m.content = "Error loading summary: " + err.Error()
		return
	}

	taskByID := make(map[int64]models.Task)
	summarySessions := make([]summarySession, 0, len(sessions))
	taskDurations := make(map[int64]time.Duration)
	for _, session := range sessions {
		task, ok := taskByID[session.TaskId]
		if !ok {
			task, err = m.store.GetTaskByID(session.TaskId)
			if err != nil {
				m.content = "Error loading task details: " + err.Error()
				return
			}
			taskByID[session.TaskId] = task
		}
		summarySessions = append(summarySessions, summarySession{session: session, task: task})
		taskDurations[task.ID] += session.EndTime.Sub(session.StartTime)
	}

	tasks := make([]summaryTask, 0, len(taskDurations))
	for taskID, duration := range taskDurations {
		tasks = append(tasks, summaryTask{task: taskByID[taskID], duration: duration})
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].task.Title < tasks[j].task.Title
	})

	m.content = renderSummary(summarySessions, tasks, m.width)
}

func renderSummary(sessions []summarySession, tasks []summaryTask, width int) string {
	const panelGap = "    "
	panelWidth := max(1, (width-4-len(panelGap))/2)
	panelStyle := lipgloss.NewStyle().Width(panelWidth)
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		panelStyle.Render(renderSummarySessions(sessions, panelWidth)),
		panelGap,
		panelStyle.Render(renderSummaryTasks(tasks, panelWidth)),
	)
}

func renderSummarySessions(sessions []summarySession, width int) string {
	heading := lipgloss.NewStyle().Bold(true).Foreground(styles.Main).Render("Sessions")
	if len(sessions) == 0 {
		return heading + "\n\nNo sessions recorded."
	}

	byDay := make(map[time.Time][]summarySession)
	for _, session := range sessions {
		day := startOfDay(session.session.StartTime)
		byDay[day] = append(byDay[day], session)
	}
	days := make([]time.Time, 0, len(byDay))
	for day := range byDay {
		days = append(days, day)
	}
	sort.Slice(days, func(i, j int) bool { return days[i].After(days[j]) })

	var body strings.Builder
	body.WriteString(heading)
	for _, day := range days {
		daySessions := byDay[day]
		sort.Slice(daySessions, func(i, j int) bool {
			return daySessions[i].session.StartTime.Before(daySessions[j].session.StartTime)
		})
		body.WriteString("\n\n")
		body.WriteString(lipgloss.NewStyle().Bold(true).Foreground(styles.Highlight1).Render(day.Format("Monday, January 2, 2006")))
		body.WriteString("\n")
		body.WriteString(summarySessionsTable(daySessions, width))
	}
	return body.String()
}

func summarySessionsTable(sessions []summarySession, width int) string {
	cellStyle := lipgloss.NewStyle().Foreground(styles.Highlight1)
	columnStyles := []lipgloss.Style{
		cellStyle,
		cellStyle.Width(20),
		cellStyle.Width(20),
	}
	t := table.New().
		Width(width).
		Border(lipgloss.HiddenBorder()).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return lipgloss.NewStyle().Bold(true).Foreground(styles.Main)
			}
			return columnStyles[col]
		}).
		Row("Task", "Started", "Duration")
	for _, session := range sessions {
		t.Row(
			fmt.Sprintf("%d %s", session.task.ID, session.task.Title),
			session.session.StartTime.Format("15:04"),
			formatDuration(session.session.EndTime.Sub(session.session.StartTime)),
		)
	}
	return t.String()
}

func renderSummaryTasks(tasks []summaryTask, width int) string {
	heading := lipgloss.NewStyle().Bold(true).Foreground(styles.Main).Render("Tasks")
	if len(tasks) == 0 {
		return heading + "\n\nNo tasks worked on."
	}

	cellStyle := lipgloss.NewStyle().Foreground(styles.Highlight1)
	columnStyles := []lipgloss.Style{
		cellStyle,
		cellStyle.Width(20),
	}
	t := table.New().
		Width(width).
		Border(lipgloss.HiddenBorder()).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return lipgloss.NewStyle().Bold(true).Foreground(styles.Main)
			}
			return columnStyles[col]
		}).
		Row("Task", "Total")
	for _, task := range tasks {
		t.Row(
			fmt.Sprintf("%d %s", task.task.ID, task.task.Title),
			formatDuration(task.duration),
		)
	}
	return heading + "\n\n" + t.String()
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
