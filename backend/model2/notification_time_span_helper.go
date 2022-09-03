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
