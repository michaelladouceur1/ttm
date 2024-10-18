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

func LogSessionEnd(session models.SessionFile) {
	LogSessionInfo(session)
}

func LogSessionInfo(session models.SessionFile) {
	taskIdString := strconv.Itoa(int(session.ID))
	fmt.Printf("Current session for task: %s \n", taskIdString)
	fmt.Printf("Start time: %s \n", session.StartTime)
	fmt.Printf("Duration: %s \n", time.Since(session.StartTime))
}

func LogSessionCancel() {
	fmt.Printf("Session cancelled\n")
}
