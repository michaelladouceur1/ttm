package paths

import (
	"os"
	"path/filepath"
)

func GetTaskStorePath() string {
	appData := os.Getenv("LOCALAPPDATA")
	return filepath.Join(appData, "ttm")
}
