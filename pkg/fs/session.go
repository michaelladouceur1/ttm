package fs

import (
	"encoding/json"
	"os"
	"time"
	"ttm/pkg/models"
	"ttm/pkg/paths"
)

func CreateSessionFile(taskId int64, startTime time.Time) (*os.File, error) {
	mutex.Lock()
	defer mutex.Unlock()

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
	mutex.Lock()
	defer mutex.Unlock()

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
	mutex.Lock()
	defer mutex.Unlock()

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
