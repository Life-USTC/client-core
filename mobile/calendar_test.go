package mobile

import (
	"testing"
	"time"
)

func TestScalarBindingUsesMillisecondsAndExplicitWeekday(t *testing.T) {
	anchor := time.Date(2026, 10, 25, 4, 0, 0, 0, time.UTC).UnixMilli()
	got, err := CalendarBounds("week", anchor, 0)
	if err != nil || got.From != "2026-10-25" || got.To != "2026-10-31" {
		t.Fatalf("Sunday bounds = %#v, %v", got, err)
	}
	if got := ShanghaiDate(anchor); got != "2026-10-25" {
		t.Fatalf("ShanghaiDate = %q", got)
	}
}

func TestBindingErrorsReturnNoSuccessValue(t *testing.T) {
	if got, err := CalendarBounds("invalid", 0, 1); err == nil || got != nil {
		t.Fatalf("invalid bounds = %#v, %v", got, err)
	}
	if got, err := RequestWindows("2026-02-30", "2026-03-01"); err == nil || got != "" {
		t.Fatalf("invalid windows = %q, %v", got, err)
	}
}

func TestWindowBindingPreservesInclusiveJSONContract(t *testing.T) {
	got, err := RequestWindows("2024-01-01", "2025-01-02")
	want := `[{"from":"2024-01-01","to":"2024-12-31"},{"from":"2025-01-01","to":"2025-01-02"}]`
	if err != nil || got != want {
		t.Fatalf("windows = %s, %v; want %s", got, err, want)
	}
}
