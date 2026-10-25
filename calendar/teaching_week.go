package calendar

import (
	"fmt"
	"time"
)

// TeachingWeekBounds returns the complete week at a one-based teaching-week
// index. Week one contains semesterStart; callers retain semester clipping.
func TeachingWeekBounds(semesterStart string, weekIndex int, firstWeekday time.Weekday) (DateRange, error) {
	if weekIndex < 1 {
		return DateRange{}, fmt.Errorf("teaching week index must be positive")
	}
	start, err := teachingWeekStart(semesterStart, firstWeekday)
	if err != nil {
		return DateRange{}, err
	}
	start = start.AddDate(0, 0, (weekIndex-1)*7)
	return DateRange{From: start.Format(time.DateOnly), To: start.AddDate(0, 0, 6).Format(time.DateOnly)}, nil
}

// TeachingWeekNumber returns the week containing selectedDate, counting the
// week containing semesterStart as one. Earlier weeks have zero/negative
// numbers; the caller decides whether that week belongs to the semester.
func TeachingWeekNumber(semesterStart, selectedDate string, firstWeekday time.Weekday) (int, error) {
	start, err := teachingWeekStart(semesterStart, firstWeekday)
	if err != nil {
		return 0, err
	}
	selected, err := teachingWeekStart(selectedDate, firstWeekday)
	if err != nil {
		return 0, err
	}
	return int((selected.Unix()-start.Unix())/(7*24*60*60)) + 1, nil
}

func teachingWeekStart(value string, firstWeekday time.Weekday) (time.Time, error) {
	date, err := parseCivilDate(value)
	if err != nil {
		return time.Time{}, err
	}
	if firstWeekday < time.Sunday || firstWeekday > time.Saturday {
		return time.Time{}, fmt.Errorf("invalid first weekday %d", firstWeekday)
	}
	offset := (int(date.Weekday()) - int(firstWeekday) + 7) % 7
	return date.AddDate(0, 0, -offset), nil
}
