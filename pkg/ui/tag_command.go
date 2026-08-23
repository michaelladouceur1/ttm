package ui

import (
	"sort"
	"strconv"
	"strings"
	"ttm/pkg/fs"
	"ttm/pkg/store"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type tagClosedMsg struct {
	content string
}

type tagModel struct {
	input       textinput.Model
	store       *store.Store
	listID      string
	content     string
	tags        []string
	suggestions []string
	selected    int
}

func newTagModel(st *store.Store, width int, listID string) tagModel {
	input := textinput.New()
	input.Prompt = "> "
	input.Placeholder = "Enter tags (comma-separated)..."
	input.Width = max(1, width-6)
	input.Focus()

	m := tagModel{
		input:   input,
		store:   st,
		listID:  listID,
		content: "Add tags to task " + listID + "\n\nTags",
	}
	m.loadTags()
	return m
}

func (m tagModel) Update(msg tea.Msg) (childModel, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
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
			if m.saveTags() {
				return m, func() tea.Msg { return tagClosedMsg{content: m.content} }
			}
			return m, nil
		case "esc":
			return m, func() tea.Msg { return tagClosedMsg{content: "Tagging cancelled."} }
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.updateSuggestions()
	return m, cmd
}

func (m tagModel) InputView() string {
	return m.input.View()
}

func (m tagModel) View() string {
	var body strings.Builder
	body.WriteString(m.content)
	if len(m.suggestions) == 0 {
		return body.String()
	}

	body.WriteString("\n\nMatching tags\n")
	for i, tag := range m.suggestions {
		prefix := "  "
		if i == m.selected {
			prefix = "> "
		}
		body.WriteString(prefix + tag + "\n")
	}
	body.WriteString("\nTab completes the selected tag.")
	return body.String()
}

func (m *tagModel) loadTags() {
	tagCounts, err := m.store.ListTagCounts()
	if err != nil {
		m.content = "Unable to load existing tags: " + err.Error()
		return
	}

	m.tags = make([]string, 0, len(tagCounts))
	for _, tag := range tagCounts {
		m.tags = append(m.tags, tag.Tag)
	}
}

func (m *tagModel) updateSuggestions() {
	m.suggestions = matchingTags(m.tags, m.input.Value())
	if m.selected >= len(m.suggestions) {
		m.selected = 0
	}
}

func (m *tagModel) completeSuggestion() {
	m.input.SetValue(completedTagValue(m.input.Value(), m.suggestions[m.selected]))
	m.input.CursorEnd()
	m.updateSuggestions()
}

func (m *tagModel) saveTags() bool {
	tags := parseTags(m.input.Value())
	if len(tags) == 0 {
		m.content = "Enter at least one tag."
		return false
	}

	listID, err := strconv.Atoi(m.listID)
	if err != nil {
		m.content = "Invalid task ID: " + m.listID
		return false
	}
	taskID, err := fs.GetTaskIDFromTempID(int64(listID))
	if err != nil {
		m.content = "Task not found: " + m.listID
		return false
	}
	if err := m.store.InsertTags(taskID, tags); err != nil {
		m.content = "Failed to add tags: " + err.Error()
		return false
	}

	m.content = "Added tags to task " + m.listID + ": " + strings.Join(tags, ", ")
	return true
}

func currentTag(value string) string {
	parts := strings.Split(value, ",")
	return strings.TrimSpace(parts[len(parts)-1])
}

func matchingTags(tags []string, value string) []string {
	current := currentTag(value)
	if current == "" {
		return nil
	}

	selectedTags := parseTags(value)
	var suggestions []string
	for _, tag := range tags {
		if strings.HasPrefix(strings.ToLower(tag), strings.ToLower(current)) && !containsTag(selectedTags, tag, current) {
			suggestions = append(suggestions, tag)
		}
	}
	sort.Slice(suggestions, func(i, j int) bool {
		return strings.ToLower(suggestions[i]) < strings.ToLower(suggestions[j])
	})
	return suggestions
}

func completedTagValue(value, suggestion string) string {
	return strings.TrimSuffix(value, currentTag(value)) + suggestion
}

func containsTag(tags []string, candidate, current string) bool {
	for _, tag := range tags {
		if strings.EqualFold(tag, candidate) && !strings.EqualFold(tag, current) {
			return true
		}
	}
	return false
}
