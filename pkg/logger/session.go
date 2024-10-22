package logger

import (
	"fmt"
	"strconv"
	"time"
	"ttm/pkg/models"
)

func LogSessionStart(taskId int64) {
	taskIdString := strconv.Itoa(int(taskId))
	fmt.Printf("Session started for task: %s \n", taskIdString)
}

func LogSessionEnd(session models.SessionFile, task models.Task) {
	LogSessionInfo(session, task)
}

func LogSessionInfo(session models.SessionFile, task models.Task) {
	timeSince := time.Since(session.StartTime).Round(time.Second)

	var startTime string
	if timeSince.Hours() > 12 {
		startTime = session.StartTime.Round(time.Second).Format("2006-01-02 15:04:05")
	} else {
		startTime = session.StartTime.Round(time.Second).Format("15:04:05")
	}

	fmt.Printf("Current session for task: %s \n", task.Title)
	fmt.Printf("Start time: %s \n", startTime)
	fmt.Printf("Duration: %s \n", timeSince)
}

func LogSessionCancel() {
	fmt.Printf("Session cancelled\n")
}
