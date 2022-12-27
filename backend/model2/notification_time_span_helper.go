package model2

import (
	"fmt"
	"time"
)

func NotificationTimeSpanTimeFormat(hour, minute int) string {
	return fmt.Sprintf("%02d:%02d:00", hour, minute)
}

func NotificationTimeSpanTimeParse(value string) (time.Time, error) {
	return time.Parse("15:04:05", value)
}

type NotificationTimeSpanList []*NotificationTimeSpan

func (l NotificationTimeSpanList) Within(t time.Time) bool {
	target := t
	for _, timeSpan := range l {
		if timeSpan.Within(target) {
			return true
		}
	}
	return false
}

func (s *NotificationTimeSpan) Within(t time.Time) bool {
	from, err := time.Parse("15:04:05", s.FromTime)
	if err != nil {
		return false
	}
	to, err := time.Parse("15:04:05", s.ToTime)
	if err != nil {
		return false
	}
	if from.Before(to) {
		if (t.After(from) || t.Equal(from)) && (t.Before(to) || t.Equal(to)) {
			return true
		}
	} else {
		// Add 24 hour to s.to if from > to. from=04:00, to=03:00 -> from=04:00, to=27:00
		toTime := to
		toTime = toTime.Add(time.Hour * 24)
		if (t.After(from) || t.Equal(from)) && (t.Before(toTime) || t.Equal(toTime)) {
			return true
		}
	}
	return false
}
