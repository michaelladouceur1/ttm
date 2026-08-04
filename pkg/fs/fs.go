package fs

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	"ttm/pkg/models"
	"ttm/pkg/paths"
)

var mutex = sync.Mutex{}

func TasksToCSV(tasks []models.Task) error {
	if err := os.MkdirAll(paths.GetTTMDirectory(), 0o755); err != nil {
		return fmt.Errorf("create TTM directory: %w", err)
	}

	file, err := os.Create(paths.GetTaskStoreCSVPath())
	if err != nil {
		return fmt.Errorf("create CSV file: %w", err)
	}

	writer := csv.NewWriter(file)
	if err := writer.Write([]string{"Title", "Description", "Priority", "Status", "Created At"}); err != nil {
		file.Close()
		return fmt.Errorf("write CSV header: %w", err)
	}

	for _, task := range tasks {
		if err := writer.Write([]string{task.Title, task.Description, string(task.Priority), string(task.Status), task.CreatedAt.Format("2006-01-02 15:04:05")}); err != nil {
			file.Close()
			return fmt.Errorf("write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		file.Close()
		return fmt.Errorf("flush CSV file: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close CSV file: %w", err)
	}
	return nil
}
