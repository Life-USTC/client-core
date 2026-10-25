package ustc_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/Life-USTC/client-core/ustc"
)

type blackboardCalendarEventExpectation struct {
	ID           string `json:"id"`
	CalendarID   string `json:"calendarId"`
	ItemSourceID string `json:"itemSourceId"`
	Title        string `json:"title"`
	Start        string `json:"start"`
	End          string `json:"end"`
	EventType    string `json:"eventType"`
	CalendarName string `json:"calendarName"`
}

type blackboardSelectedEventExpectation struct {
	Title      string `json:"title"`
	CourseName string `json:"courseName"`
	EndDate    string `json:"endDate"`
}

func TestBlackboardSharedFixtures(t *testing.T) {
	data, err := os.ReadFile("testdata/blackboard_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures struct {
		SelectedEventsURL, PageURL, CalendarsURL, CoursesURL string
		EventsURLs                                           []struct {
			IDs []string
			URL string
		}
		GradesURLs []struct {
			ID, URL string
		}
		CalendarIDs []struct {
			Name, Response string
			IDs            []string
			Error          bool
		}
		CalendarEvents []struct {
			Name, Response string
			Events         []blackboardCalendarEventExpectation
			Error          bool
		}
		SelectedEvents []struct {
			Name, Response string
			Events         []blackboardSelectedEventExpectation
			Error          bool
		}
		CourseIDs []struct {
			Name, HTML string
			IDs        map[string]string
		}
		GradeStatuses []struct {
			Name, HTML string
			Statuses   map[string]string
		}
		CalendarIdentities []struct {
			ID, LessonCode, SemesterCode string
		}
		Batches []struct {
			Name    string
			IDs     []string
			Max     int
			Batches [][]string
		}
	}
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	for name, count := range map[string]int{
		"events URLs":         len(fixtures.EventsURLs),
		"grades URLs":         len(fixtures.GradesURLs),
		"calendar IDs":        len(fixtures.CalendarIDs),
		"calendar events":     len(fixtures.CalendarEvents),
		"selected events":     len(fixtures.SelectedEvents),
		"course IDs":          len(fixtures.CourseIDs),
		"grade statuses":      len(fixtures.GradeStatuses),
		"calendar identities": len(fixtures.CalendarIdentities),
		"batches":             len(fixtures.Batches),
	} {
		if count == 0 {
			t.Fatalf("missing shared Blackboard %s fixtures", name)
		}
	}
	for _, tc := range []struct {
		name, got, want string
	}{
		{"selected-events", ustc.BlackboardSelectedCalendarEventsURL, fixtures.SelectedEventsURL},
		{"page", ustc.BlackboardCalendarPageURL, fixtures.PageURL},
		{"calendars", ustc.BlackboardCalendarsURL, fixtures.CalendarsURL},
		{"courses", ustc.BlackboardCoursesURL, fixtures.CoursesURL},
	} {
		t.Run("url/"+tc.name, func(t *testing.T) {
			if tc.want == "" {
				t.Fatal("missing shared Blackboard URL fixture")
			}
			if tc.got != tc.want {
				t.Fatalf("URL = %q, want %q", tc.got, tc.want)
			}
		})
	}
	for index, tc := range fixtures.EventsURLs {
		t.Run(fmt.Sprintf("events-url/%d", index), func(t *testing.T) {
			if got := ustc.BlackboardCalendarEventsURL(tc.IDs); got != tc.URL {
				t.Fatalf("BlackboardCalendarEventsURL = %q, want %q", got, tc.URL)
			}
		})
	}
	for index, tc := range fixtures.GradesURLs {
		t.Run(fmt.Sprintf("grades-url/%d", index), func(t *testing.T) {
			if got := ustc.BlackboardGradesURL(tc.ID); got != tc.URL {
				t.Fatalf("BlackboardGradesURL = %q, want %q", got, tc.URL)
			}
		})
	}
	for _, tc := range fixtures.CalendarIDs {
		t.Run("calendar-ids/"+tc.Name, func(t *testing.T) {
			got, err := ustc.ParseBlackboardCalendarIDs([]byte(tc.Response))
			if (err != nil) != tc.Error {
				t.Fatalf("ParseBlackboardCalendarIDs error = %v, want error=%v", err, tc.Error)
			}
			if !reflect.DeepEqual(got, tc.IDs) {
				t.Fatalf("calendar IDs = %#v, want %#v", got, tc.IDs)
			}
		})
	}
	for _, tc := range fixtures.CalendarEvents {
		t.Run("calendar-events/"+tc.Name, func(t *testing.T) {
			got, err := ustc.ParseBlackboardCalendarEvents([]byte(tc.Response))
			if (err != nil) != tc.Error {
				t.Fatalf("ParseBlackboardCalendarEvents error = %v, want error=%v", err, tc.Error)
			}
			var actual []blackboardCalendarEventExpectation
			if got != nil {
				actual = make([]blackboardCalendarEventExpectation, len(got))
				for index, event := range got {
					actual[index] = blackboardCalendarEventExpectation{
						ID:           event.ID,
						CalendarID:   event.CalendarID,
						ItemSourceID: event.ItemSourceID,
						Title:        event.Title,
						Start:        event.Start,
						End:          event.End,
						EventType:    event.EventType,
						CalendarName: event.CalendarName,
					}
				}
			}
			if !reflect.DeepEqual(actual, tc.Events) {
				t.Fatalf("calendar events = %#v, want %#v", actual, tc.Events)
			}
		})
	}
	for _, tc := range fixtures.SelectedEvents {
		t.Run("selected-events/"+tc.Name, func(t *testing.T) {
			got, err := ustc.ParseBlackboardSelectedCalendarEvents([]byte(tc.Response))
			if (err != nil) != tc.Error {
				t.Fatalf("ParseBlackboardSelectedCalendarEvents error = %v, want error=%v", err, tc.Error)
			}
			var actual []blackboardSelectedEventExpectation
			if got != nil {
				actual = make([]blackboardSelectedEventExpectation, len(got))
				for index, event := range got {
					actual[index] = blackboardSelectedEventExpectation{
						Title:      event.Title,
						CourseName: event.CourseName,
						EndDate:    event.EndDate,
					}
				}
			}
			if !reflect.DeepEqual(actual, tc.Events) {
				t.Fatalf("selected calendar events = %#v, want %#v", actual, tc.Events)
			}
		})
	}
	for _, tc := range fixtures.CourseIDs {
		t.Run("course-ids/"+tc.Name, func(t *testing.T) {
			got, err := ustc.ParseBlackboardCourseIDs(strings.NewReader(tc.HTML))
			if err != nil {
				t.Fatalf("ParseBlackboardCourseIDs error = %v", err)
			}
			if !reflect.DeepEqual(got, tc.IDs) {
				t.Fatalf("course IDs = %#v, want %#v", got, tc.IDs)
			}
		})
	}
	for _, tc := range fixtures.GradeStatuses {
		t.Run("grade-statuses/"+tc.Name, func(t *testing.T) {
			got := ustc.ParseBlackboardGradeStatuses(strings.NewReader(tc.HTML))
			if !reflect.DeepEqual(got, tc.Statuses) {
				t.Fatalf("grade statuses = %#v, want %#v", got, tc.Statuses)
			}
		})
	}
	for index, tc := range fixtures.CalendarIdentities {
		t.Run(fmt.Sprintf("calendar-identity/%d", index), func(t *testing.T) {
			lessonCode, semesterCode := ustc.SplitBlackboardCalendarID(tc.ID)
			if lessonCode != tc.LessonCode || semesterCode != tc.SemesterCode {
				t.Fatalf("SplitBlackboardCalendarID(%q) = %q, %q; want %q, %q", tc.ID, lessonCode, semesterCode, tc.LessonCode, tc.SemesterCode)
			}
		})
	}
	for _, tc := range fixtures.Batches {
		t.Run("batches/"+tc.Name, func(t *testing.T) {
			got := ustc.BatchBlackboardCalendarIDs(tc.IDs, tc.Max)
			if !reflect.DeepEqual(got, tc.Batches) {
				t.Fatalf("calendar batches = %#v, want %#v", got, tc.Batches)
			}
		})
	}
}

type blackboardErrorReader struct {
	err error
}

func (reader blackboardErrorReader) Read([]byte) (int, error) {
	return 0, reader.err
}

func TestBlackboardHTMLReaderErrors(t *testing.T) {
	readError := errors.New("Blackboard fixture read failed")
	t.Run("course-ids", func(t *testing.T) {
		got, err := ustc.ParseBlackboardCourseIDs(blackboardErrorReader{err: readError})
		if !errors.Is(err, readError) {
			t.Fatalf("ParseBlackboardCourseIDs error = %v, want %v", err, readError)
		}
		if got != nil {
			t.Fatalf("failed course ID parsing returned %#v, want nil", got)
		}
	})
	t.Run("grade-statuses", func(t *testing.T) {
		got := ustc.ParseBlackboardGradeStatuses(blackboardErrorReader{err: readError})
		if got != nil {
			t.Fatalf("failed grade status parsing returned %#v, want nil", got)
		}
	})
}

func TestBlackboardSelectedCalendarEventsRejectInvalidUTF8(t *testing.T) {
	data := append([]byte(`[{"title":"`), 0xff)
	data = append(data, []byte(`"}]`)...)
	got, err := ustc.ParseBlackboardSelectedCalendarEvents(data)
	if err == nil {
		t.Fatal("invalid UTF-8 selected calendar events did not return an error")
	}
	if got != nil {
		t.Fatalf("invalid UTF-8 selected calendar events returned %#v, want nil", got)
	}
}

func TestBlackboardDisabledBatchBudgetCopiesInput(t *testing.T) {
	for _, maxQueryLength := range []int{0, -1} {
		t.Run(fmt.Sprintf("max=%d", maxQueryLength), func(t *testing.T) {
			ids := []string{"A", "B"}
			got := ustc.BatchBlackboardCalendarIDs(ids, maxQueryLength)
			want := [][]string{{"A", "B"}}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("calendar batches = %#v, want %#v", got, want)
			}
			got[0][0] = "changed batch"
			if !reflect.DeepEqual(ids, []string{"A", "B"}) {
				t.Fatalf("changing a batch mutated the input: %#v", ids)
			}
		})
	}
}
