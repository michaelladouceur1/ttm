package googledocs

import (
	"testing"
	"ttm/pkg/models"

	docs "google.golang.org/api/docs/v1"
)

func TestNextIDs(t *testing.T) {
	if got := nextTaskID([]models.Task{{ID: 2}, {ID: 5}}); got != 6 {
		t.Fatalf("nextTaskID() = %d, want 6", got)
	}
	if got := nextSessionID([]models.Session{{ID: 3}, {ID: 4}}); got != 5 {
		t.Fatalf("nextSessionID() = %d, want 5", got)
	}
}

func TestDocumentText(t *testing.T) {
	document := &docs.Document{
		Body: &docs.Body{Content: []*docs.StructuralElement{
			{Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{
				{TextRun: &docs.TextRun{Content: "{\"tasks\":"}},
				{TextRun: &docs.TextRun{Content: "[]}"}},
			}}},
		}},
	}

	if got := documentText(document); got != "{\"tasks\":[]}" {
		t.Fatalf("documentText() = %q", got)
	}
}
