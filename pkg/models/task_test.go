package models

import "testing"

func TestTaskValidateRejectsInvalidCategory(t *testing.T) {
	task := Task{
		Category: "invalid",
		Priority: PriorityHigh,
		Status:   StatusOpen,
	}

	if err := task.Validate(); err == nil {
		t.Fatal("Validate() returned nil for an invalid category")
	}
}

func TestSortTasksByIDAndPopulateListIDs(t *testing.T) {
	tasks := []Task{{ID: 3}, {ID: 1}, {ID: 2}}

	SortTasksByID(tasks)
	PopulateListIDs(tasks)

	for i, wantID := range []int64{1, 2, 3} {
		if tasks[i].ID != wantID {
			t.Fatalf("tasks[%d].ID = %d, want %d", i, tasks[i].ID, wantID)
		}
		if tasks[i].ListID != int64(i+1) {
			t.Fatalf("tasks[%d].ListID = %d, want %d", i, tasks[i].ListID, i+1)
		}
	}
}
