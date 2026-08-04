package fs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"ttm/pkg/models"
	"ttm/pkg/paths"
)

type IDMap struct {
	ID     int64
	ListID int64
}

func ReadIDMapFile() ([]IDMap, error) {
	mutex.Lock()
	defer mutex.Unlock()

	data, err := os.ReadFile(paths.GetIDMapPath())
	if err != nil {
		return nil, fmt.Errorf("read task ID map: %w", err)
	}

	var idMap []IDMap
	if err := json.Unmarshal(data, &idMap); err != nil {
		return nil, fmt.Errorf("decode task ID map: %w", err)
	}

	return idMap, nil
}

var ErrTaskNotFound = errors.New("task not found")

func GetTaskIDFromTempID(listID int64) (int64, error) {
	idMap, err := ReadIDMapFile()
	if err != nil {
		return 0, err
	}

	for _, idMapItem := range idMap {
		if idMapItem.ListID == listID {
			return idMapItem.ID, nil
		}
	}

	return 0, ErrTaskNotFound
}

func UpdateIDMapFile(tasks []models.Task) error {
	mutex.Lock()
	defer mutex.Unlock()

	var idMap []IDMap
	for _, task := range tasks {
		idMap = append(idMap, IDMap{ID: task.ID, ListID: task.ListID})
	}

	idMapJSON, err := json.Marshal(idMap)
	if err != nil {
		return fmt.Errorf("encode task ID map: %w", err)
	}

	if err := os.MkdirAll(paths.GetTTMDirectory(), 0o755); err != nil {
		return fmt.Errorf("create TTM directory: %w", err)
	}
	if err := os.WriteFile(paths.GetIDMapPath(), idMapJSON, 0o600); err != nil {
		return fmt.Errorf("write task ID map: %w", err)
	}
	return nil
}
