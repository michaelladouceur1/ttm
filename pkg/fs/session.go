package fs

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"ttm/pkg/models"
	"ttm/pkg/paths"
)

func CreateSessionFile(taskID int64, startTime time.Time) error {
	mutex.Lock()
	defer mutex.Unlock()

	data, err := json.Marshal(models.SessionFile{TaskID: taskID, StartTime: startTime})
	if err != nil {
		return fmt.Errorf("encode session file: %w", err)
	}

	if err := os.MkdirAll(paths.GetTTMDirectory(), 0o755); err != nil {
		return fmt.Errorf("create TTM directory: %w", err)
	}
	if err := os.WriteFile(paths.GetSessionPath(), data, 0o600); err != nil {
		return fmt.Errorf("write session file: %w", err)
	}
	return nil
}

func RemoveSessionFile() error {
	mutex.Lock()
	defer mutex.Unlock()

	if err := os.Remove(paths.GetSessionPath()); err != nil {
		return fmt.Errorf("remove session file: %w", err)
	}
	return nil
}

func ReadSessionFile() (models.SessionFile, error) {
	mutex.Lock()
	defer mutex.Unlock()

	data, err := os.ReadFile(paths.GetSessionPath())
	if err != nil {
		return models.SessionFile{}, fmt.Errorf("read session file: %w", err)
	}

	var session models.SessionFile
	if err := json.Unmarshal(data, &session); err != nil {
		return models.SessionFile{}, fmt.Errorf("decode session file: %w", err)
	}
	return session, nil
}

func SessionFileExists() bool {
	_, err := os.Stat(paths.GetSessionPath())
	return err == nil
}
