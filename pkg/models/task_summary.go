package models

import "time"

type TaskSummary struct {
	Days []TaskSummaryDay
}

type TaskSummaryDay struct {
	Day   time.Time
	Tasks []Task
}

func (ts *TaskSummary) AddDay(day TaskSummaryDay) {
	ts.Days = append(ts.Days, day)
}
