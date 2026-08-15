package ui

import (
	"fmt"
	"strings"
	"ttm/pkg/models"
	"ttm/pkg/store"
)

func listTags(st *store.Store) string {
	tags, err := st.ListTagCounts()
	if err != nil {
		return "Error listing tags: " + err.Error()
	}
	if len(tags) == 0 {
		return "No tags found."
	}

	return renderTagCounts(tags)
}

func renderTagCounts(tags []models.TagCount) string {
	var content strings.Builder
	content.WriteString("Tags\n")
	for _, tag := range tags {
		fmt.Fprintf(&content, "\n(%d) %s", tag.Count, tag.Tag)
	}
	return content.String()
}
