package paths

import (
	"os"
	"path/filepath"
)

func GetTaskStorePath() string {
	appData := os.Getenv("HOME")
	return filepath.Join(appData, ".ttm")
}
