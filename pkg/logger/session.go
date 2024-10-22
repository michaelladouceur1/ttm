package logger

import (
	"fmt"
	"time"
	"ttm/pkg/models"
)

func LogSessionStart(task models.Task) {
	fmt.Printf("Session started for task: %s \n", task.Title)
}

func LogSessionEnd(session models.SessionFile, task models.Task) {
	LogSessionInfo(session, task)
}

func LogSessionInfo(session models.SessionFile, task models.Task) {
	startTime := session.StartTime.Round(time.Second).Format("15:04:05")
	timeSince := time.Since(session.StartTime).Round(time.Second)

	fmt.Printf("Current session for task: %s \n", task.Title)
	fmt.Printf("Start time: %s \n", startTime)
	fmt.Printf("Duration: %s \n", timeSince)
}

func LogSessionCancel() {
	fmt.Printf("Session cancelled\n")
}
