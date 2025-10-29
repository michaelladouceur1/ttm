package fs

import (
	"encoding/csv"
	"os"
	"sync"
	"ttm/pkg/models"
	"ttm/pkg/paths"
)

var mutex = sync.Mutex{}

func TasksToCSV(tasks []models.Task) error {
	var err error

	if os.MkdirAll(paths.GetTTMDirectory(), os.ModePerm); err != nil {
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
