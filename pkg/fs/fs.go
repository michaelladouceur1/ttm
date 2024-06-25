package fs

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"time"
	"ttm/pkg/models"
	"ttm/pkg/paths"
)

func CreateSessionFile(taskId string, startTime time.Time) (*os.File, error) {
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
