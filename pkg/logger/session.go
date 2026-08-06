package logger

import (
	"fmt"
	"time"
	"ttm/pkg/models"
)

func (l *Logger) LogSessionStart(task models.Task, start time.Time) {
	fmt.Println(l.style.RenderSummary("Session Started", []SummaryItem{
		{Key: "Task Title", Value: task.Title},
		{Key: "Task Description", Value: task.Description},
		{Key: "Start Time", Value: start.Round(time.Second).Format("2006-01-02 15:04:05")},
	}))
}

func (l *Logger) LogSessionEnd(session models.SessionFile, task models.Task) {
	l.logSessionSummary(session, task, "Session Ended")
}

func (l *Logger) LogSessionInfo(session models.SessionFile, task models.Task) {
	l.logSessionSummary(session, task, "Session Info")
}

func (l *Logger) LogSessionCancel() {
	l.LogMessage("Session cancelled.")
}

func (l *Logger) logSessionSummary(session models.SessionFile, task models.Task, title string) {
	fmt.Println(l.style.RenderSummary(title, []SummaryItem{
		{Key: "Task Title", Value: task.Title},
		{Key: "Start Time", Value: session.StartTime.Round(time.Second).Format("2006-01-02 15:04:05")},
		{Key: "Duration", Value: time.Since(session.StartTime).Round(time.Second).String()},
	}))
}

func LogSessionStart(task models.Task, start time.Time) {
	defaultLogger.LogSessionStart(task, start)
}

func LogSessionEnd(session models.SessionFile, task models.Task) {
	defaultLogger.LogSessionEnd(session, task)
}

func LogSessionInfo(session models.SessionFile, task models.Task) {
	defaultLogger.LogSessionInfo(session, task)
}

func LogSessionCancel() {
	defaultLogger.LogSessionCancel()
}
