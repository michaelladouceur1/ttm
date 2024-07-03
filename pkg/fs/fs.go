package fs

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"time"
	"ttm/pkg/models"
	"ttm/pkg/paths"
)

func CreateSessionFile(taskId int64, startTime time.Time) (*os.File, error) {
	var err error

	sfJson, err := json.Marshal(models.SessionFile{ID: taskId, StartTime: startTime})
	if err != nil {
		return nil, err
	}

	file, err := os.Create(paths.GetSessionPath())
	if err != nil {
		return nil, err
	}

	_, err = file.Write(sfJson)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func RemoveSessionFile() (models.SessionFile, error) {
	var sf models.SessionFile

	file, err := os.Open(paths.GetSessionPath())
	if err != nil {
		return sf, err
	}

	err = json.NewDecoder(file).Decode(&sf)
	if err != nil {
		return sf, err
	}

	file.Close()

	err = os.Remove(paths.GetSessionPath())
	if err != nil {
		return sf, err
	}

	return sf, nil
}

func ReadSessionFile() (models.SessionFile, error) {
	var sf models.SessionFile

	file, err := os.Open(paths.GetSessionPath())
	if err != nil {
		return sf, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&sf)
	if err != nil {
		return sf, err
	}

	return sf, nil
}

func SessionFileExists() bool {
	_, err := os.Stat(paths.GetSessionPath())
	return err == nil
}

type IDMap struct {
	ID     int64
	ListID int64
}

func ReadIDMapFile() ([]IDMap, error) {
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

func GetTaskIDFromListID(listID int64) (int64, error) {
	idMap, err := ReadIDMapFile()
	if err != nil {
		return 0, err
	}

	for _, idMapItem := range idMap {
		if idMapItem.ListID == listID {
			return idMapItem.ID, nil
		}
	}

	return 0, TaskNotFoundError{}
}

func UpdateIDMapFile(tasks []models.Task) error {
	var err error

	if os.MkdirAll(paths.GetTaskStorePath(), os.ModePerm); err != nil {
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

func TasksToCSV(tasks []models.Task) error {
	var err error

	if os.MkdirAll(paths.GetTaskStorePath(), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(paths.GetTaskStoreCSVPath())
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Title", "Description", "Priority", "Status", "Created At"})
	if err != nil {
		return err
	}

	for _, task := range tasks {
		err := writer.Write([]string{task.Title, task.Description, string(task.Priority), string(task.Status), task.CreatedAt.Format("2006-01-02 15:04:05")})
		if err != nil {
			return err
		}
	}

	return nil
}
