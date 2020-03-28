package batch

import (
	"time"
)

type ScheduleItem struct {
	StartedAt uint16 // minute from 00:00
	Window    time.Duration
}

type Schedule []ScheduleItem

func NewSchedule(cfg []Rule) (*Schedule, error) {
	schedule := make(Schedule, len(cfg))

	for i, rule := range cfg {
		t, err := time.Parse("15:04", rule.StartedAt)
		if err != nil {
			return nil, err
		}

		schedule[i] = ScheduleItem{
			StartedAt: uint16(t.Hour()*60 + t.Minute()),
			Window:    rule.Window,
		}
	}

	return &schedule, nil
}

func (s Schedule) GetWindow(t time.Time) time.Duration {
	minutesFromMidnight := uint16(t.Hour()*60 + t.Minute())

	index := s.indexOf(minutesFromMidnight)
	if index == -1 {
		return 0
	}

	if len(s) == 1 {
		return s[index].Window
	}

	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	// find minutes from midnight for the next schedule rule
	minutes := s[0].StartedAt
	if (index + 1) < len(s) {
		minutes = s[index+1].StartedAt
	}

	// if the rule is based on the midnight (00:00) we should use next midnight for calculation
	if minutes == 0 {
		minutes = 1440 // 24*60
	}

	nextRuleTime := midnight.Add(time.Duration(minutes) * time.Minute)
	d := nextRuleTime.Sub(t)
	if d < 0 {
		return s[index].Window
	}

	if d >= s[index].Window {
		return s[index].Window
	}

	return d
}

func (s Schedule) indexOf(minutesFromMidnight uint16) int {
	if len(s) == 0 {
		return -1
	}

	for i := len(s) - 1; i >= 0; i-- {
		if s[i].StartedAt > minutesFromMidnight {
			continue
		}

		return i
	}

	return 0
}
