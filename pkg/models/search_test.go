package models

import (
	"reflect"
	"testing"
)

func TestParseTaskSearchPlainText(t *testing.T) {
	got, err := ParseTaskSearch("release notes")
	if err != nil {
		t.Fatalf("ParseTaskSearch() error = %v", err)
	}
	want := TaskSearch{General: "release notes"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseTaskSearch() = %#v, want %#v", got, want)
	}
}

func TestParseTaskSearchFilters(t *testing.T) {
	got, err := ParseTaskSearch("$tags:work, urgent $title:Task 1 $status:open")
	if err != nil {
		t.Fatalf("ParseTaskSearch() error = %v", err)
	}
	want := TaskSearch{Tags: []string{"work", "urgent"}, Title: "Task 1", Status: "open"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseTaskSearch() = %#v, want %#v", got, want)
	}
}

func TestParseTaskSearchRejectsInvalidFilters(t *testing.T) {
	for _, input := range []string{"", "$title:", "$unknown:value"} {
		if _, err := ParseTaskSearch(input); err == nil {
			t.Errorf("ParseTaskSearch(%q) returned nil error", input)
		}
	}
}
