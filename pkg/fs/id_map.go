package fs

import (
	"encoding/json"
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

	var idMap []IDMap

	file, err := os.Open(paths.GetIDMapPath())
	if err != nil {
		return idMap, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&idMap)
	if err != nil {
		return idMap, err
	}

	return idMap, nil
}

type TaskNotFoundError struct{}

func (e TaskNotFoundError) Error() string {
	return "task not found"
}

func GetTaskIDFromTempID(tempID int64) (int64, error) {
	idMap, err := ReadIDMapFile()
	if err != nil {
		return 0, err
	}

	for _, idMapItem := range idMap {
		if idMapItem.ListID == tempID {
			return idMapItem.ID, nil
		}
	}

	return 0, TaskNotFoundError{}
}

func UpdateIDMapFile(tasks []models.Task) error {
	mutex.Lock()
	defer mutex.Unlock()

	var err error

	if os.MkdirAll(paths.GetTTMDirectory(), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(paths.GetIDMapPath())
	if err != nil {
		return err
	}
	defer file.Close()

	var idMap []IDMap
	for _, task := range tasks {
		idMap = append(idMap, IDMap{ID: task.ID, ListID: task.ListID})
	}

	idMapJson, err := json.Marshal(idMap)
	if err != nil {
		return err
	}

	_, err = file.Write(idMapJson)
	if err != nil {
		return err
	}

	return nil
}
