package paths

import (
	"os"
	"path/filepath"
)

func GetTTMDirectory() string {
	return filepath.Join(os.Getenv("HOME"), ".ttm")
}
