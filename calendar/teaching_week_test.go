package calendar_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/Life-USTC/client-core/calendar"
)

type teachingWeekFixtures struct {
	Bounds []struct {
		Name          string
		SemesterStart string
		WeekIndex     int
		FirstWeekday  time.Weekday
		From, To      string
	}
	Numbers []struct {
		Name          string
		SemesterStart string
		SelectedDate  string
		FirstWeekday  time.Weekday
		WeekIndex     int
	}
}

func TestSharedTeachingWeekFixtures(t *testing.T) {
	data, err := os.ReadFile("testdata/teaching_week_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures teachingWeekFixtures
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	if len(fixtures.Bounds) == 0 || len(fixtures.Numbers) == 0 {
		t.Fatal("missing shared teaching-week fixtures")
	}
	for _, tc := range fixtures.Bounds {
		t.Run("bounds/"+tc.Name, func(t *testing.T) {
			got, err := calendar.TeachingWeekBounds(tc.SemesterStart, tc.WeekIndex, tc.FirstWeekday)
			want := calendar.DateRange{From: tc.From, To: tc.To}
			if err != nil || got != want {
				t.Fatalf("TeachingWeekBounds = %+v, %v; want %+v", got, err, want)
			}
		})
	}
	for _, tc := range fixtures.Numbers {
		t.Run("number/"+tc.Name, func(t *testing.T) {
			got, err := calendar.TeachingWeekNumber(tc.SemesterStart, tc.SelectedDate, tc.FirstWeekday)
			if err != nil || got != tc.WeekIndex {
				t.Fatalf("TeachingWeekNumber = %d, %v; want %d", got, err, tc.WeekIndex)
			}
		})
	}
}

func TestTeachingWeekRejectsInvalidInputs(t *testing.T) {
	for _, tc := range []struct {
		start string
		week  int
		first time.Weekday
	}{
		{"2026-09-01", 0, time.Monday}, {"2026-09-01", -1, time.Sunday},
		{"2026-02-30", 1, time.Monday}, {"2026-9-01", 1, time.Sunday},
		{"2026-09-01", 1, -1}, {"2026-09-01", 1, 7},
	} {
		if got, err := calendar.TeachingWeekBounds(tc.start, tc.week, tc.first); err == nil || got != (calendar.DateRange{}) {
			t.Errorf("TeachingWeekBounds(%q, %d, %d) = %+v, %v", tc.start, tc.week, tc.first, got, err)
		}
	}
	for _, tc := range []struct {
		start, selected string
		first           time.Weekday
	}{
		{"bad", "2026-09-01", time.Monday}, {"2026-09-01", "2026-09-31", time.Sunday},
		{"2026-09-01", "2026-09-01", -1}, {"2026-09-01", "2026-09-01", 7},
	} {
		if got, err := calendar.TeachingWeekNumber(tc.start, tc.selected, tc.first); err == nil || got != 0 {
			t.Errorf("TeachingWeekNumber(%q, %q, %d) = %d, %v", tc.start, tc.selected, tc.first, got, err)
		}
	}
}
