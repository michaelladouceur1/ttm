package ui

import (
	"strings"
	"ttm/pkg/store"
	"ttm/pkg/styles"
)

func listNotes(st *store.Store, taskID int64) string {
	notes, err := st.GetNotesByTaskID(taskID)
	if err != nil {
		return "Error listing task notes: " + err.Error()
	}
	if len(notes) == 0 {
		return "No notes found."
	}

	var body strings.Builder
	for _, note := range notes {
		body.WriteString(styles.Bullet)
		body.WriteString(" ")
		body.WriteString(note.Content)
		body.WriteByte('\n')
	}
	return strings.TrimRight(body.String(), "\n")
}
