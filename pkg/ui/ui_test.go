package ui

import (
	"reflect"
	"strings"
	"testing"
	"time"
	"ttm/pkg/config"
	"ttm/pkg/models"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{input: `/list`, want: []string{"list"}},
		{input: `/add "Plan release" "Prepare release notes"`, want: []string{"add", "Plan release", "Prepare release notes"}},
		{input: `/add 'Plan release'`, want: []string{"add", "Plan release"}},
	}

	for _, test := range tests {
		got, err := parseCommand(test.input)
		if err != nil {
			t.Fatalf("parseCommand(%q) returned an error: %v", test.input, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("parseCommand(%q) = %#v, want %#v", test.input, got, test.want)
		}
	}
}

func TestParseCommandRejectsInvalidInput(t *testing.T) {
	for _, input := range []string{"add task", `/add "unfinished`} {
		if _, err := parseCommand(input); err == nil {
			t.Errorf("parseCommand(%q) returned nil error", input)
		}
	}
}

func TestUpdateSuggestions(t *testing.T) {
	m := newModel(nil, nil)
	m.input.SetValue("/")
	m.updateSuggestions()
	if len(m.suggestions) != len(commands) {
		t.Fatalf("suggestions = %d, want %d", len(m.suggestions), len(commands))
	}

	m.input.SetValue("/tag")
	m.updateSuggestions()
	if len(m.suggestions) != 1 || m.suggestions[0].name != "tags" {
		t.Fatalf("suggestions = %#v, want tags", m.suggestions)
	}
}

func TestAddFormWalksThroughFields(t *testing.T) {
	m := newModel(&config.Config{
		AddFlags: config.ConfigDefaultFlags{
			Category: string(models.CategoryTask),
			Priority: string(models.PriorityHigh),
			Status:   string(models.StatusOpen),
		},
	}, nil)
	m.input.SetValue("/add")
	m.execute()

	add, ok := m.active.(addModel)
	if !ok {
		t.Fatalf("active model = %T, want addModel", m.active)
	}
	if add.step != addStepTitle || add.input.Placeholder != "Enter title..." {
		t.Fatalf("add form did not start at title: step=%d placeholder=%q", add.step, add.input.Placeholder)
	}

	add.input.SetValue("Plan release")
	updated, _ := add.submitAddField()
	add = updated.(addModel)
	if add.step != addStepDescription || add.draft.Title != "Plan release" {
		t.Fatalf("title was not saved: step=%d title=%q", add.step, add.draft.Title)
	}

	add.input.SetValue("Prepare notes")
	updated, _ = add.submitAddField()
	add = updated.(addModel)
	if add.step != addStepPriority || add.draft.Description != "Prepare notes" {
		t.Fatalf("description was not saved: step=%d description=%q", add.step, add.draft.Description)
	}

	add.input.SetValue("medium")
	updated, _ = add.submitAddField()
	add = updated.(addModel)
	if add.step != addStepTags || add.draft.Priority != models.PriorityMedium {
		t.Fatalf("priority was not saved: step=%d priority=%q", add.step, add.draft.Priority)
	}
}

func TestParseTags(t *testing.T) {
	got := parseTags("work, urgent, ,planning")
	want := []string{"work", "urgent", "planning"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseTags() = %#v, want %#v", got, want)
	}
}

func TestTasksCommandCreatesChildModel(t *testing.T) {
	m := newModel(nil, nil)
	m.input.SetValue("/tasks one two")
	m.execute()

	list, ok := m.active.(tasksModel)
	if !ok {
		t.Fatalf("active model = %T, want tasksModel", m.active)
	}
	if list.content != "Error: /list accepts at most one search query." {
		t.Errorf("tasks content = %q", list.content)
	}
}

func TestRenderTagCounts(t *testing.T) {
	rendered := renderTagCounts([]models.TagCount{
		{Tag: "planning", Count: 2},
		{Tag: "work", Count: 1},
	})
	if rendered != "Tags\n\nplanning (2)\nwork (1)" {
		t.Errorf("renderTagCounts() = %q", rendered)
	}
}

func TestTagsModelClosesToRoot(t *testing.T) {
	m := tagsModel{content: "Tags\n\nwork (1)"}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	msg := cmd()

	if _, ok := msg.(tagsClosedMsg); !ok {
		t.Fatalf("close message = %T, want tagsClosedMsg", msg)
	}
}

func TestListModelClosesToRoot(t *testing.T) {
	m := newTasksModel(nil, nil, 80, []string{"one", "two"})
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	msg := cmd()

	if _, ok := msg.(tasksClosedMsg); !ok {
		t.Fatalf("close message = %T, want tasksClosedMsg", msg)
	}
}

func TestRootViewportScrollsCommandContent(t *testing.T) {
	m := newModel(nil, nil)
	m.viewport = viewport.New(20, 1)
	m.active = detailModel{content: "first\nsecond"}
	m.setViewportContent()

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(model)
	if m.viewport.YOffset != 1 {
		t.Errorf("viewport offset = %d, want 1", m.viewport.YOffset)
	}
}

func TestRenderNotesShowsCreationTimeAndContent(t *testing.T) {
	notes := []models.Note{
		{
			Content:   "Follow up with the design team.",
			CreatedAt: time.Date(2026, time.August, 14, 20, 1, 0, 0, time.UTC),
		},
	}

	rendered := (&detailModel{}).renderNotes(notes, 80)
	for _, want := range []string{"Notes", "Created At", "2026-08-14 20:01", "Follow up with the design team."} {
		if !strings.Contains(rendered, want) {
			t.Errorf("renderNotes() = %q, want it to contain %q", rendered, want)
		}
	}
}
