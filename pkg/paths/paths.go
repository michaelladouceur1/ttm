package paths

import (
	"os"
	"path/filepath"
	"time"
)

func GetTTMDirectory() string {
	return filepath.Join(os.Getenv("HOME"), ".ttm")
}

func GetTaskStoreDBPath() string {
	return filepath.Join(GetTTMDirectory(), "ttm.db")
}

func GetConfigPath() string {
	return filepath.Join(GetTTMDirectory(), "config.json")
}

func GetSessionPath() string {
	return filepath.Join(GetTTMDirectory(), "session.json")
}

func GetIDMapPath() string {
	return filepath.Join(GetTTMDirectory(), "id_map.json")
}

func GetTaskStoreCSVPath() string {
	dateTime := time.Now().Format("2006-01-02_15-04-05")
	return filepath.Join(GetTTMDirectory(), "tasks_"+dateTime+".csv")
}
