package ui

import (
	"reflect"
	"testing"
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
