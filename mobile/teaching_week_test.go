package mobile

import "testing"

func TestTeachingWeekBindingPreservesPolicyAndSignedNumbers(t *testing.T) {
	got, err := TeachingWeekBounds("2026-09-01", 2, 0)
	if err != nil || got.From != "2026-09-06" || got.To != "2026-09-12" {
		t.Fatalf("TeachingWeekBounds = %#v, %v", got, err)
	}
	number, err := TeachingWeekNumber("2026-09-01", "2026-08-23", 1)
	if err != nil || number != -1 {
		t.Fatalf("TeachingWeekNumber = %d, %v", number, err)
	}
}

func TestTeachingWeekBindingDoesNotReturnAnEmptySuccess(t *testing.T) {
	if got, err := TeachingWeekBounds("2026-09-01", 0, 1); err == nil || got != nil {
		t.Fatalf("TeachingWeekBounds = %#v, %v", got, err)
	}
	if got, err := TeachingWeekNumber("2026-09-01", "bad", 1); err == nil || got != 0 {
		t.Fatalf("TeachingWeekNumber = %d, %v", got, err)
	}
}
