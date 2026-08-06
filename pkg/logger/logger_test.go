package logger

import (
	"strings"
	"testing"
	"ttm/pkg/models"
)

func TestNewStyle(t *testing.T) {
	for _, name := range []string{ClassicStyleName, CompactStyleName, "COMPACT"} {
		if _, err := NewStyle(name); err != nil {
			t.Errorf("NewStyle(%q) error = %v", name, err)
		}
	}
	if _, err := NewStyle("unknown"); err == nil {
		t.Error("NewStyle(unknown) returned nil error")
	}
}

func TestCompactStyleRendersTasks(t *testing.T) {
	rendered := CompactStyle{}.RenderTasks([]models.Task{{
		ID:     7,
		Title:  "Write tests",
		Status: models.StatusOpen,
		Tags:   []string{"work"},
	}})

	for _, value := range []string{"ID\tTitle\tStatus\tTags", "7\tWrite tests\topen\twork"} {
		if !strings.Contains(rendered, value) {
			t.Errorf("RenderTasks() = %q, missing %q", rendered, value)
		}
	}
}
