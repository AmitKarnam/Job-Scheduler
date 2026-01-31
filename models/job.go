package models

import (
	"fmt"
	"time"
)

// JobType represents the kind of job to be executed.
type JobType string

const (
	EmailJob           JobType = "email_job"
	MobileNotification JobType = "mobile_notification"
)

// Status represents the current state of a job.
type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

// ScheduleType is a lightweight identifier for schedule kinds.
type ScheduleType string

const (
	ScheduleOneTime ScheduleType = "one_time"
	ScheduleFixed   ScheduleType = "fixed_interval"
	ScheduleMonthly ScheduleType = "monthly"
	ScheduleDaily   ScheduleType = "daily"
	ScheduleWeekly  ScheduleType = "weekly"
	ScheduleCron    ScheduleType = "cron"
)

// Schedule defines how a job should be scheduled. Implementations return
// the next execution time after a provided timestamp. If no next run exists
// the boolean returned should be false.
type Schedule interface {
	NextAfter(time.Time) (time.Time, bool, error)
	IsRecurring() bool
	Type() ScheduleType
}

// OneTimeSchedule is a single execution at a specific timestamp.
type OneTimeSchedule struct {
	At time.Time `json:"at"`
}

func (s OneTimeSchedule) NextAfter(after time.Time) (time.Time, bool, error) {
	if s.At.After(after) {
		return s.At, true, nil
	}
	return time.Time{}, false, nil
}

func (s OneTimeSchedule) IsRecurring() bool  { return false }
func (s OneTimeSchedule) Type() ScheduleType { return ScheduleOneTime }

// FixedIntervalSchedule executes starting at `Start` and repeats every
// `Interval` duration. If MaxRuns is > 0 it's a bounded repeating schedule.
type FixedIntervalSchedule struct {
	Start    time.Time     `json:"start"`
	Interval time.Duration `json:"interval"`
	MaxRuns  int           `json:"max_runs,omitempty"` // 0 = unbounded
}

func (s FixedIntervalSchedule) NextAfter(after time.Time) (time.Time, bool, error) {
	if s.Interval <= 0 {
		return time.Time{}, false, fmt.Errorf("invalid interval")
	}
	// If after is before start, next is start
	if after.Before(s.Start) {
		return s.Start, true, nil
	}
	elapsed := after.Sub(s.Start)
	// number of intervals passed
	n := int(elapsed / s.Interval)
	next := s.Start.Add(time.Duration(n+1) * s.Interval)
	if s.MaxRuns > 0 {
		// compute how many runs would have happened including this next
		if n+1 > s.MaxRuns {
			return time.Time{}, false, nil
		}
	}
	return next, true, nil
}

func (s FixedIntervalSchedule) IsRecurring() bool  { return s.MaxRuns == 0 || s.MaxRuns > 1 }
func (s FixedIntervalSchedule) Type() ScheduleType { return ScheduleFixed }

// MonthlySchedule schedules a job once per month on a specified day/time.
// Day is 1-31; if the month doesn't have that day it will pick the last day.
type MonthlySchedule struct {
	Day    int `json:"day"`    // 1..31
	Hour   int `json:"hour"`   // 0..23
	Minute int `json:"minute"` // 0..59
}

func (s MonthlySchedule) NextAfter(after time.Time) (time.Time, bool, error) {
	// Start searching from next month up to N months ahead (safe guard)
	candidate := time.Date(after.Year(), after.Month(), 1, s.Hour, s.Minute, 0, 0, after.Location())
	// ensure candidate is not before `after`
	for i := 0; i < 120; i++ { // search up to 10 years
		year, month := candidate.Year(), candidate.Month()
		// determine day: clamp to last day of month
		day := s.Day
		// get last day of month
		firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, candidate.Location())
		lastDay := firstOfNextMonth.Add(-time.Hour * 24).Day()
		if day > lastDay {
			day = lastDay
		}
		cand := time.Date(year, month, day, s.Hour, s.Minute, 0, 0, candidate.Location())
		if cand.After(after) {
			return cand, true, nil
		}
		// move to first of next month
		candidate = firstOfNextMonth
	}
	return time.Time{}, false, fmt.Errorf("could not find next monthly occurrence")
}

func (s MonthlySchedule) IsRecurring() bool  { return true }
func (s MonthlySchedule) Type() ScheduleType { return ScheduleMonthly }

// DailySchedule schedules a job every day at a specified hour/minute.
type DailySchedule struct {
	Hour   int `json:"hour"`   // 0..23
	Minute int `json:"minute"` // 0..59
}

func (s DailySchedule) NextAfter(after time.Time) (time.Time, bool, error) {
	loc := after.Location()
	candidate := time.Date(after.Year(), after.Month(), after.Day(), s.Hour, s.Minute, 0, 0, loc)
	if candidate.After(after) {
		return candidate, true, nil
	}
	next := candidate.Add(24 * time.Hour)
	return next, true, nil
}

func (s DailySchedule) IsRecurring() bool  { return true }
func (s DailySchedule) Type() ScheduleType { return ScheduleDaily }

// WeeklySchedule schedules a job on a specific weekday at hour/minute.
type WeeklySchedule struct {
	Day    time.Weekday `json:"day"`    // time.Sunday..time.Saturday
	Hour   int          `json:"hour"`   // 0..23
	Minute int          `json:"minute"` // 0..59
}

func (s WeeklySchedule) NextAfter(after time.Time) (time.Time, bool, error) {
	loc := after.Location()
	// compute days until target weekday
	current := after.Weekday()
	daysUntil := (int(s.Day) - int(current) + 7) % 7
	candidate := time.Date(after.Year(), after.Month(), after.Day(), s.Hour, s.Minute, 0, 0, loc).Add(time.Duration(daysUntil) * 24 * time.Hour)
	if candidate.After(after) {
		return candidate, true, nil
	}
	// if candidate is not after (same day earlier time), move to next week's occurrence
	candidate = candidate.Add(7 * 24 * time.Hour)
	return candidate, true, nil
}

func (s WeeklySchedule) IsRecurring() bool  { return true }
func (s WeeklySchedule) Type() ScheduleType { return ScheduleWeekly }

// CronSchedule stores a cron expression. Computing next times requires a
// cron parsing library (e.g. robfig/cron). For now this is a placeholder
// that preserves the spec and reports unimplemented for NextAfter.
type CronSchedule struct {
	Spec string `json:"spec"`
}

func (s CronSchedule) NextAfter(after time.Time) (time.Time, bool, error) {
	return time.Time{}, false, fmt.Errorf("cron schedule next computation not implemented; add a cron parser")
}

func (s CronSchedule) IsRecurring() bool  { return true }
func (s CronSchedule) Type() ScheduleType { return ScheduleCron }

// Job is the canonical job definition used across the scheduler.
type Job struct {
	ID         string                 `json:"id"`
	Type       JobType                `json:"type"`
	Schedule   Schedule               `json:"-"` // interface; JSON marshalling requires a helper
	Payload    map[string]interface{} `json:"payload,omitempty"`
	RetryCount int                    `json:"retry_count"`
	MaxRetries int                    `json:"max_retries"`
	Status     Status                 `json:"status"`
	LastError  string                 `json:"last_error,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// NextRun returns the next scheduled execution time after the provided time.
// If no next run exists the boolean will be false.
func (j *Job) NextRun(after time.Time) (time.Time, bool, error) {
	if j.Schedule == nil {
		return time.Time{}, false, fmt.Errorf("no schedule defined")
	}
	return j.Schedule.NextAfter(after)
}

// Execute runs the job action. This is a placeholder implementation and
// should be replaced by concrete job behaviour or by implementing a
// separate jobs interface for specific job types.
func (j *Job) Execute() error {
	// TODO: implement actual execution logic (worker, HTTP call, etc.)
	return nil
}
