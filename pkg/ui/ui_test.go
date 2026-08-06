package ui

import (
	"reflect"
	"testing"
	"ttm/pkg/config"
	"ttm/pkg/models"
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
	if len(m.suggestions) != 2 {
		t.Fatalf("suggestions = %d, want 2", len(m.suggestions))
	}

	m.input.SetValue("/li")
	m.updateSuggestions()
	if len(m.suggestions) != 1 || m.suggestions[0].name != "list" {
		t.Fatalf("suggestions = %#v, want list", m.suggestions)
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

	if m.addStep != addStepTitle || m.input.Placeholder != "Enter title..." {
		t.Fatalf("add form did not start at title: step=%d placeholder=%q", m.addStep, m.input.Placeholder)
	}

	m.input.SetValue("Plan release")
	m.submitAddField()
	if m.addStep != addStepDescription || m.draft.Title != "Plan release" {
		t.Fatalf("title was not saved: step=%d title=%q", m.addStep, m.draft.Title)
	}

	m.input.SetValue("Prepare notes")
	m.submitAddField()
	if m.addStep != addStepPriority || m.draft.Description != "Prepare notes" {
		t.Fatalf("description was not saved: step=%d description=%q", m.addStep, m.draft.Description)
	}

	m.input.SetValue("medium")
	m.submitAddField()
	if m.addStep != addStepTags || m.draft.Priority != models.PriorityMedium {
		t.Fatalf("priority was not saved: step=%d priority=%q", m.addStep, m.draft.Priority)
	}
}

func TestParseTags(t *testing.T) {
	got := parseTags("work, urgent, ,planning")
	want := []string{"work", "urgent", "planning"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseTags() = %#v, want %#v", got, want)
	}
}
