// Package mobile exposes the client core through gomobile-supported types.
package mobile

import (
	"encoding/json"
	"time"

	"github.com/Life-USTC/client-core/calendar"
)

// Range has inclusive Asia/Shanghai civil-date bounds in YYYY-MM-DD format.
type Range struct {
	From string
	To   string
}

// CalendarBounds plans a day, week, or month around an instant. FirstWeekday is
// explicit: 0 is Sunday, 1 is Monday, matching time.Weekday and the shared core.
func CalendarBounds(period string, anchorUnixMillis int64, firstWeekday int64) (*Range, error) {
	value, err := calendar.Bounds(calendar.Period(period), time.UnixMilli(anchorUnixMillis), time.Weekday(firstWeekday))
	if err != nil {
		return nil, err
	}
	return &Range{From: value.From, To: value.To}, nil
}

// ShanghaiDate converts an instant to its Asia/Shanghai civil date.
func ShanghaiDate(unixMillis int64) string {
	return calendar.ShanghaiDate(time.UnixMilli(unixMillis))
}

// RequestWindows returns the ordered inclusive ranges as a JSON array. Each
// object has string fields "from" and "to"; no window contains over 366 days.
func RequestWindows(from, to string) (string, error) {
	windows, err := calendar.RequestWindows(from, to)
	if err != nil {
		return "", err
	}
	encoded, err := json.Marshal(windows)
	return string(encoded), err
}
