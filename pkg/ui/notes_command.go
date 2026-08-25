package ui

import (
	"strings"
	"ttm/pkg/store"
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
		body.WriteString("- ")
		body.WriteString(note.Content)
		body.WriteByte('\n')
	}
	return strings.TrimRight(body.String(), "\n")
}
