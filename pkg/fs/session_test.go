package fs

import (
	"errors"
	"os"
	"testing"
	"time"
)

func TestSessionFileLifecycle(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	startTime := time.Date(2026, 8, 3, 12, 0, 0, 0, time.UTC)
	if err := CreateSessionFile(42, startTime); err != nil {
		t.Fatalf("CreateSessionFile() error = %v", err)
	}

	session, err := ReadSessionFile()
	if err != nil {
		t.Fatalf("ReadSessionFile() error = %v", err)
	}
	if session.TaskID != 42 {
		t.Fatalf("session.TaskID = %d, want 42", session.TaskID)
	}
	if !session.StartTime.Equal(startTime) {
		t.Fatalf("session.StartTime = %v, want %v", session.StartTime, startTime)
	}

	if err := RemoveSessionFile(); err != nil {
		t.Fatalf("RemoveSessionFile() error = %v", err)
	}

	_, err = ReadSessionFile()
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("ReadSessionFile() error = %v, want an os.ErrNotExist error", err)
	}
}
