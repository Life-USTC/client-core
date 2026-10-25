package mobile

import (
	"encoding/json"
	"os"
	"slices"
	"testing"
)

func TestUSTCCourseTableBindingAgainstSharedFixtures(t *testing.T) {
	data, err := os.ReadFile("../ustc/testdata/course_table_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures struct {
		PageURL  string
		Requests []struct {
			Name, StudentID, URL string
			SemesterID           int64
		}
		StudentIDs []struct {
			Name, URL, StudentID string
			Error                bool
		}
		Tables []struct {
			Name, Response string
			LessonIDs      []int64
			Error          bool
		}
	}
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	if fixtures.PageURL == "" || len(fixtures.Requests) == 0 || len(fixtures.StudentIDs) == 0 || len(fixtures.Tables) == 0 {
		t.Fatal("missing shared USTC course table fixtures")
	}
	if got := USTCCourseTablePageURL(); got != fixtures.PageURL {
		t.Fatalf("USTCCourseTablePageURL = %q, want %q", got, fixtures.PageURL)
	}
	for _, tc := range fixtures.Requests {
		t.Run("request/"+tc.Name, func(t *testing.T) {
			if got := USTCCourseTableDataURL(tc.SemesterID, tc.StudentID); got != tc.URL {
				t.Fatalf("USTCCourseTableDataURL = %q, want %q", got, tc.URL)
			}
		})
	}
	for _, tc := range fixtures.StudentIDs {
		t.Run("student-id/"+tc.Name, func(t *testing.T) {
			got, err := USTCStudentIDFromCourseTableURL(tc.URL)
			if (err != nil) != tc.Error || got != tc.StudentID {
				t.Fatalf("USTCStudentIDFromCourseTableURL = %q, %v; want %q, error=%v", got, err, tc.StudentID, tc.Error)
			}
		})
	}
	for _, tc := range fixtures.Tables {
		t.Run("response/"+tc.Name, func(t *testing.T) {
			got, err := ParseUSTCCourseTable([]byte(tc.Response))
			if (err != nil) != tc.Error {
				t.Fatalf("ParseUSTCCourseTable error = %v, want error=%v", err, tc.Error)
			}
			if tc.Error {
				if got != nil {
					t.Fatal("failed parsing returned identifiers")
				}
				return
			}
			if got == nil {
				t.Fatal("successful parsing returned nil instead of an identifier collection")
			}
			ids := make([]int64, got.Count())
			for index := range ids {
				ids[index] = got.At(int64(index))
			}
			if !slices.Equal(ids, tc.LessonIDs) {
				t.Fatalf("lesson identifiers = %v, want %v", ids, tc.LessonIDs)
			}
		})
	}
}
