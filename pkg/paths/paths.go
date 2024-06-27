package paths

import (
	"path/filepath"
	"time"
)

func GetTaskStoreDBPath() string {
	return filepath.Join(GetTaskStorePath(), "ttm.db")
}

func GetConfigPath() string {
	return filepath.Join(GetTaskStorePath(), "config.json")
}

func GetSessionPath() string {
	return filepath.Join(GetTaskStorePath(), "session.json")
}

func GetIDMapPath() string {
	return filepath.Join(GetTaskStorePath(), "id_map.json")
}

func GetTaskStoreCSVPath() string {
	dateTime := time.Now().Format("2006-01-02_15-04-05")
	return filepath.Join(GetTaskStorePath(), "tasks_"+dateTime+".csv")
}
