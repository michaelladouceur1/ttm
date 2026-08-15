package ui

import (
	"fmt"
	"strings"
	"ttm/pkg/models"
	"ttm/pkg/store"

	tea "github.com/charmbracelet/bubbletea"
)

type tagsClosedMsg struct {
	content string
}

type tagsModel struct {
	content string
}

func newTagsModel(st *store.Store) tagsModel {
	m := tagsModel{}
	m.listTags(st)
	return m
}

func (m tagsModel) Update(msg tea.Msg) (childModel, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "esc" {
		return m, func() tea.Msg { return tagsClosedMsg{content: m.content} }
	}
	return m, nil
}

func (m tagsModel) InputView() string {
	return "Press Esc to return"
}

func (m tagsModel) View() string {
	return m.content
}

func (m *tagsModel) listTags(st *store.Store) {
	tags, err := st.ListTagCounts()
	if err != nil {
		m.content = "Error listing tags: " + err.Error()
		return
	}
	if len(tags) == 0 {
		m.content = "No tags found."
		return
	}

	m.content = renderTagCounts(tags)
}

func renderTagCounts(tags []models.TagCount) string {
	var content strings.Builder
	content.WriteString("Tags\n")
	for _, tag := range tags {
		fmt.Fprintf(&content, "\n(%d) %s", tag.Count, tag.Tag)
	}
	return content.String()
}
