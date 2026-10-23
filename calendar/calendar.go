// Package calendar defines the Shanghai civil-date rules shared by Life@USTC clients.
package calendar

import (
	"fmt"
	"time"
	_ "time/tzdata" // Apple builds must not depend on host zoneinfo files.
)

type Period string

const (
	Day   Period = "day"
	Week  Period = "week"
	Month Period = "month"
)

// DateRange contains inclusive Shanghai calendar dates.
type DateRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

var shanghai = func() *time.Location {
	zone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err) // The zone is supplied by the embedded time/tzdata package.
	}
	return zone
}()

// ShanghaiDate converts an instant to its Shanghai calendar date.
func ShanghaiDate(at time.Time) string {
	return at.In(shanghai).Format(time.DateOnly)
}

// ParseDate parses a canonical YYYY-MM-DD date in Asia/Shanghai.
func ParseDate(value string) (time.Time, error) {
	date, err := parseCivilDate(value)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, shanghai), nil
}

// Bounds returns the complete day, week or month containing anchor.
// firstWeekday is explicit because public and personal calendars differ.
func Bounds(period Period, anchor time.Time, firstWeekday time.Weekday) (DateRange, error) {
	if firstWeekday < time.Sunday || firstWeekday > time.Saturday {
		return DateRange{}, fmt.Errorf("invalid first weekday %d", firstWeekday)
	}
	local := anchor.In(shanghai)
	// UTC holds civil fields here; daylight-saving transitions must not shift dates.
	start := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, time.UTC)
	end := start
	switch period {
	case Day:
	case Week:
		offset := (int(start.Weekday()) - int(firstWeekday) + 7) % 7
		start = start.AddDate(0, 0, -offset)
		end = start.AddDate(0, 0, 6)
	case Month:
		start = time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC)
		end = start.AddDate(0, 1, -1)
	default:
		return DateRange{}, fmt.Errorf("unsupported calendar period %q", period)
	}
	return DateRange{From: start.Format(time.DateOnly), To: end.Format(time.DateOnly)}, nil
}

// RequestWindows divides an inclusive range into adjacent windows of at most 366 days.
func RequestWindows(from, to string) ([]DateRange, error) {
	start, err := parseCivilDate(from)
	if err != nil {
		return nil, err
	}
	end, err := parseCivilDate(to)
	if err != nil {
		return nil, err
	}
	if start.After(end) {
		return nil, fmt.Errorf("calendar range start must not follow its end")
	}
	var windows []DateRange
	for cursor := start; !cursor.After(end); {
		last := cursor.AddDate(0, 0, 365)
		if last.After(end) {
			last = end
		}
		windows = append(windows, DateRange{From: cursor.Format(time.DateOnly), To: last.Format(time.DateOnly)})
		cursor = last.AddDate(0, 0, 1)
	}
	return windows, nil
}

func parseCivilDate(value string) (time.Time, error) {
	date, err := time.Parse(time.DateOnly, value)
	if err != nil || date.Format(time.DateOnly) != value {
		return time.Time{}, fmt.Errorf("invalid calendar date %q (use YYYY-MM-DD)", value)
	}
	return date, nil
}
