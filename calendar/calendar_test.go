package calendar_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/Life-USTC/client-core/calendar"
)

type fixtures struct {
	Bounds []struct {
		Name         string
		Period       calendar.Period
		Anchor       time.Time
		FirstWeekday time.Weekday
		From, To     string
	}
	Windows []struct {
		Name     string
		From, To string
		Ranges   []calendar.DateRange
	}
}

func readFixtures(t *testing.T) fixtures {
	t.Helper()
	data, err := os.ReadFile("testdata/calendar_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases fixtures
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases.Bounds) == 0 || len(cases.Windows) == 0 {
		t.Fatal("shared fixtures must cover bounds and request windows")
	}
	return cases
}

func TestSharedBoundsFixtures(t *testing.T) {
	for _, tc := range readFixtures(t).Bounds {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			got, err := calendar.Bounds(tc.Period, tc.Anchor, tc.FirstWeekday)
			if err != nil {
				t.Fatal(err)
			}
			want := calendar.DateRange{From: tc.From, To: tc.To}
			if got != want {
				t.Fatalf("Bounds = %+v, want %+v", got, want)
			}
		})
	}
}

func TestSharedRequestWindowFixtures(t *testing.T) {
	for _, tc := range readFixtures(t).Windows {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			got, err := calendar.RequestWindows(tc.From, tc.To)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tc.Ranges) {
				t.Fatalf("RequestWindows = %+v, want %+v", got, tc.Ranges)
			}
		})
	}
}

func TestParseDateRequiresCanonicalCivilDate(t *testing.T) {
	for _, raw := range []string{"", "2026-2-03", "2026-02-3", "2026-02-30", "2026-13-01", " 2026-10-23", "2026-10-23 ", "2026-10-23T00:00:00Z"} {
		t.Run(raw, func(t *testing.T) {
			if _, err := calendar.ParseDate(raw); err == nil {
				t.Fatalf("ParseDate(%q) accepted an invalid date", raw)
			}
		})
	}
	got, err := calendar.ParseDate("2026-10-23")
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 10, 22, 16, 0, 0, 0, time.UTC)
	if !got.Equal(want) || got.Location().String() != "Asia/Shanghai" {
		t.Fatalf("ParseDate = %v in %s, want Shanghai midnight", got, got.Location())
	}
}

func TestShanghaiDateUsesInstantAndIANAZone(t *testing.T) {
	for _, tc := range []struct{ instant, date string }{
		{"2026-10-22T15:59:59.999Z", "2026-10-22"},
		{"2026-10-22T16:00:00Z", "2026-10-23"},
		{"2026-10-23T00:00:00+08:00", "2026-10-23"},
		{"1986-07-01T15:30:00Z", "1986-07-02"},
	} {
		instant, err := time.Parse(time.RFC3339Nano, tc.instant)
		if err != nil {
			t.Fatal(err)
		}
		if got := calendar.ShanghaiDate(instant); got != tc.date {
			t.Errorf("ShanghaiDate(%s) = %s, want %s", tc.instant, got, tc.date)
		}
	}
}

func TestBoundsRejectsUnsupportedPeriodAndWeekday(t *testing.T) {
	anchor := time.Date(2026, 10, 23, 0, 0, 0, 0, time.UTC)
	for _, period := range []calendar.Period{"", "year", "Week"} {
		if _, err := calendar.Bounds(period, anchor, time.Monday); err == nil {
			t.Errorf("Bounds accepted period %q", period)
		}
	}
	for _, weekday := range []time.Weekday{-1, 7} {
		if _, err := calendar.Bounds(calendar.Week, anchor, weekday); err == nil {
			t.Errorf("Bounds accepted weekday %d", weekday)
		}
	}
}

func TestRequestWindowsRejectsInvalidOrReversedBounds(t *testing.T) {
	for _, bounds := range [][2]string{
		{"2026-10-24", "2026-10-23"},
		{"bad", "2026-10-23"},
		{"2026-10-23", "bad"},
		{"2026-02-30", "2026-03-01"},
	} {
		if _, err := calendar.RequestWindows(bounds[0], bounds[1]); err == nil {
			t.Errorf("RequestWindows accepted %q..%q", bounds[0], bounds[1])
		}
	}
}
