package ustc

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

type staticFixtures struct {
	GraduateAppCoursesURL, StaticSemestersPath string
	Paths                                      []struct {
		SemesterID, CourseList, Course string
		CourseID                       int
	}
	Semesters []struct {
		Name, Response string
		Semesters      []StaticSemester
		Error          bool
	}
	CourseLists []struct {
		Name, Response string
		Courses        []StaticCourseListItem
		Error          bool
	}
	Courses []struct {
		Name, Response string
		Course         *StaticCourse
		Error          bool
	}
	AppCodes []struct {
		Name, Response string
		Codes          []string
		Error          bool
	}
	Matches []struct {
		Name, AppResponse, CourseListResponse string
		IDs                                   []int
		Error                                 bool
	}
}

func readStaticFixtures(t *testing.T) staticFixtures {
	t.Helper()
	data, err := os.ReadFile("testdata/static_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures staticFixtures
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	if len(fixtures.Paths) == 0 || len(fixtures.Semesters) == 0 || len(fixtures.CourseLists) == 0 || len(fixtures.Courses) == 0 || len(fixtures.AppCodes) == 0 || len(fixtures.Matches) == 0 {
		t.Fatal("missing shared USTC static fixtures")
	}
	return fixtures
}

func TestStaticPaths(t *testing.T) {
	fixtures := readStaticFixtures(t)
	if GraduateAppCoursesURL != fixtures.GraduateAppCoursesURL || StaticSemestersPath != fixtures.StaticSemestersPath {
		t.Fatal("static/app endpoint differs from existing native transport")
	}
	for _, tc := range fixtures.Paths {
		if got := StaticCourseListPath(tc.SemesterID); got != tc.CourseList {
			t.Errorf("course list path = %q, want %q", got, tc.CourseList)
		}
		if got := StaticCoursePath(tc.SemesterID, tc.CourseID); got != tc.Course {
			t.Errorf("course path = %q, want %q", got, tc.Course)
		}
	}
}

func TestStaticSemesters(t *testing.T) {
	for _, tc := range readStaticFixtures(t).Semesters {
		t.Run(tc.Name, func(t *testing.T) {
			got, err := ParseStaticSemesters([]byte(tc.Response))
			if (err != nil) != tc.Error {
				t.Fatalf("error = %v, want error=%v", err, tc.Error)
			}
			if !tc.Error && !reflect.DeepEqual(got, tc.Semesters) {
				t.Fatalf("semesters = %+v, want %+v", got, tc.Semesters)
			}
		})
	}
}

func TestStaticCourseList(t *testing.T) {
	for _, tc := range readStaticFixtures(t).CourseLists {
		t.Run(tc.Name, func(t *testing.T) {
			got, err := ParseStaticCourseList([]byte(tc.Response))
			if (err != nil) != tc.Error {
				t.Fatalf("error = %v, want error=%v", err, tc.Error)
			}
			if !tc.Error && !reflect.DeepEqual(got, tc.Courses) {
				t.Fatalf("courses = %+v, want %+v", got, tc.Courses)
			}
		})
	}
}

func TestStaticCourse(t *testing.T) {
	for _, tc := range readStaticFixtures(t).Courses {
		t.Run(tc.Name, func(t *testing.T) {
			got, err := ParseStaticCourse([]byte(tc.Response))
			if (err != nil) != tc.Error {
				t.Fatalf("error = %v, want error=%v", err, tc.Error)
			}
			if !tc.Error && (tc.Course == nil || !reflect.DeepEqual(got, *tc.Course)) {
				t.Fatalf("course = %+v, want %+v", got, tc.Course)
			}
		})
	}
}

func TestGraduateAppCourseCodes(t *testing.T) {
	for _, tc := range readStaticFixtures(t).AppCodes {
		t.Run(tc.Name, func(t *testing.T) {
			got, err := ParseGraduateAppCourseCodes([]byte(tc.Response))
			if (err != nil) != tc.Error {
				t.Fatalf("error = %v, want error=%v", err, tc.Error)
			}
			if !tc.Error && !reflect.DeepEqual(got, tc.Codes) {
				t.Fatalf("codes = %#v, want %#v", got, tc.Codes)
			}
		})
	}
}

func TestMatchStaticCourses(t *testing.T) {
	for _, tc := range readStaticFixtures(t).Matches {
		t.Run(tc.Name, func(t *testing.T) {
			codes, err := ParseGraduateAppCourseCodes([]byte(tc.AppResponse))
			var courses []StaticCourseListItem
			if err == nil {
				courses, err = ParseStaticCourseList([]byte(tc.CourseListResponse))
			}
			if (err != nil) != tc.Error {
				t.Fatalf("error = %v, want error=%v", err, tc.Error)
			}
			if !tc.Error {
				if got := MatchStaticCourses(courses, codes); !reflect.DeepEqual(got, tc.IDs) {
					t.Fatalf("IDs = %v, want %v", got, tc.IDs)
				}
			}
		})
	}
}
