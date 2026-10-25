package mobile

import (
	"time"

	"github.com/Life-USTC/client-core/calendar"
)

// TeachingWeekBounds returns a complete seven-day range. Inputs are Shanghai
// YYYY-MM-DD dates; firstWeekday uses 0 for Sunday and 1 for Monday.
func TeachingWeekBounds(semesterStart string, weekIndex, firstWeekday int64) (*Range, error) {
	value, err := calendar.TeachingWeekBounds(semesterStart, int(weekIndex), time.Weekday(firstWeekday))
	if err != nil {
		return nil, err
	}
	return &Range{From: value.From, To: value.To}, nil
}

// TeachingWeekNumber counts the week containing semesterStart as one.
func TeachingWeekNumber(semesterStart, selectedDate string, firstWeekday int64) (int64, error) {
	value, err := calendar.TeachingWeekNumber(semesterStart, selectedDate, time.Weekday(firstWeekday))
	return int64(value), err
}
