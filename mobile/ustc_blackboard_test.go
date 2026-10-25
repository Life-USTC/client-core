package mobile

import (
	"encoding/json"
	"os"
	"slices"
	"testing"
)

func TestUSTCBlackboardBindingAgainstSharedFixtures(t *testing.T) {
	data, err := os.ReadFile("../ustc/testdata/blackboard_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	type event struct {
		Title      string `json:"title"`
		CourseName string `json:"courseName"`
		EndDate    string `json:"endDate"`
	}
	var fixtures struct {
		SelectedEventsURL string
		SelectedEvents    []struct {
			Name, Response string
			Events         []event
			Error          bool
		}
	}
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	if fixtures.SelectedEventsURL == "" || len(fixtures.SelectedEvents) == 0 {
		t.Fatal("missing shared USTC Blackboard fixtures")
	}
	if got := USTCBlackboardSelectedCalendarEventsURL(); got != fixtures.SelectedEventsURL {
		t.Fatalf("USTCBlackboardSelectedCalendarEventsURL = %q, want %q", got, fixtures.SelectedEventsURL)
	}
	for _, tc := range fixtures.SelectedEvents {
		t.Run(tc.Name, func(t *testing.T) {
			got, err := ParseUSTCBlackboardSelectedCalendarEvents([]byte(tc.Response))
			if (err != nil) != tc.Error {
				t.Fatalf("ParseUSTCBlackboardSelectedCalendarEvents error = %v, want error=%v", err, tc.Error)
			}
			if tc.Error {
				if got != nil {
					t.Fatal("failed parsing returned calendar events")
				}
				return
			}
			if got == nil {
				t.Fatal("successful parsing returned nil instead of a calendar event collection")
			}
			events := make([]event, got.Count())
			for index := range events {
				events[index] = event{
					Title:      got.Title(int64(index)),
					CourseName: got.CourseName(int64(index)),
					EndDate:    got.EndDate(int64(index)),
				}
			}
			if !slices.Equal(events, tc.Events) {
				t.Fatalf("calendar events = %+v, want %+v", events, tc.Events)
			}
		})
	}
}
