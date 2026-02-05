package timekit

import "time"

func TruncateToHour(t time.Time) time.Time {
	return time.Date(
		t.Year(), t.Month(), t.Day(), t.Hour(),
		0, 0, 0, t.Location(),
	)
}
