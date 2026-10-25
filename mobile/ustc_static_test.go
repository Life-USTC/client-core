package mobile

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/Life-USTC/client-core/ustc"
)

type staticBindingFixtures struct {
	GraduateAppCoursesURL, StaticSemestersPath string
	Paths                                      []struct {
		SemesterID, CourseList, Course string
		CourseID                       int64
	}
	Semesters []struct {
		Name, Response string
		Semesters      []ustc.StaticSemester
		Error          bool
	}
	CourseLists []struct {
		Name, Response string
		Courses        []ustc.StaticCourseListItem
		Error          bool
	}
	Courses []struct {
		Name, Response string
		Course         *ustc.StaticCourse
		Error          bool
	}
	Matches []struct {
		Name, AppResponse, CourseListResponse, ErrorStage string
		IDs                                               []int
		Error                                             bool
	}
}

func readStaticBindingFixtures(t *testing.T) staticBindingFixtures {
	t.Helper()
	data, err := os.ReadFile("../ustc/testdata/static_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures staticBindingFixtures
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	if len(fixtures.Paths) == 0 || len(fixtures.Semesters) == 0 || len(fixtures.CourseLists) == 0 || len(fixtures.Courses) == 0 || len(fixtures.Matches) == 0 {
		t.Fatal("missing static binding fixtures")
	}
	return fixtures
}

func TestStaticBindingPaths(t *testing.T) {
	fixtures := readStaticBindingFixtures(t)
	if USTCGraduateAppCoursesURL() != fixtures.GraduateAppCoursesURL || USTCStaticSemestersPath() != fixtures.StaticSemestersPath {
		t.Fatal("static/app endpoint differs")
	}
	for _, tc := range fixtures.Paths {
		if got := USTCStaticCourseListPath(tc.SemesterID); got != tc.CourseList {
			t.Errorf("course list = %q, want %q", got, tc.CourseList)
		}
		if got := USTCStaticCoursePath(tc.SemesterID, tc.CourseID); got != tc.Course {
			t.Errorf("course = %q, want %q", got, tc.Course)
		}
	}
}

func TestStaticSemesterBinding(t *testing.T) {
	for _, tc := range readStaticBindingFixtures(t).Semesters {
		t.Run(tc.Name, func(t *testing.T) {
			value, err := ParseUSTCStaticSemesters([]byte(tc.Response))
			if (err != nil) != tc.Error {
				t.Fatalf("error=%v, want error=%v", err, tc.Error)
			}
			if tc.Error {
				if value != nil {
					t.Fatal("failed parsing returned semesters")
				}
				return
			}
			actual := make([]ustc.StaticSemester, value.Count())
			for i := range actual {
				v := value.At(int64(i))
				actual[i] = ustc.StaticSemester{ID: v.Identifier, Name: v.Name, StartDate: v.StartDate, EndDate: v.EndDate}
			}
			if !reflect.DeepEqual(actual, tc.Semesters) {
				t.Fatalf("semesters=%+v, want %+v", actual, tc.Semesters)
			}
		})
	}
}

func TestStaticCourseListBinding(t *testing.T) {
	for _, tc := range readStaticBindingFixtures(t).CourseLists {
		t.Run(tc.Name, func(t *testing.T) {
			value, err := ParseUSTCStaticCourseList([]byte(tc.Response))
			if (err != nil) != tc.Error {
				t.Fatalf("error=%v, want error=%v", err, tc.Error)
			}
			if tc.Error {
				if value != nil {
					t.Fatal("failed parsing returned courses")
				}
				return
			}
			actual := make([]ustc.StaticCourseListItem, value.Count())
			for i := range actual {
				v := value.At(int64(i))
				actual[i] = ustc.StaticCourseListItem{ID: int(v.Identifier), Name: v.Name, CourseCode: v.CourseCode, LessonCode: v.LessonCode, LegacyLessonCode: v.LegacyLessonCode, TeacherName: v.TeacherName}
			}
			if !reflect.DeepEqual(actual, tc.Courses) {
				t.Fatalf("courses=%+v, want %+v", actual, tc.Courses)
			}
		})
	}
}

func TestStaticCourseBinding(t *testing.T) {
	for _, tc := range readStaticBindingFixtures(t).Courses {
		t.Run(tc.Name, func(t *testing.T) {
			v, err := ParseUSTCStaticCourse([]byte(tc.Response))
			if (err != nil) != tc.Error {
				t.Fatalf("error=%v, want error=%v", err, tc.Error)
			}
			if tc.Error {
				if v != nil {
					t.Fatal("failed parsing returned course")
				}
				return
			}
			actual := ustc.StaticCourse{ID: int(v.Identifier), Name: v.Name, CourseCode: v.CourseCode, LessonCode: v.LessonCode, TeacherName: v.TeacherName, Credit: v.Credit, Lectures: []ustc.StaticLecture{}}
			if v.HasDetailText {
				actual.DetailText = &v.DetailText
			}
			if v.HasDateTimePlacePersonText {
				actual.DateTimePlacePersonText = &v.DateTimePlacePersonText
			}
			for i := int64(0); i < v.LectureCount(); i++ {
				x := v.LectureAt(i)
				lecture := ustc.StaticLecture{StartDate: x.StartDate, EndDate: x.EndDate, Name: x.Name, Location: x.Location, TeacherName: x.TeacherName, Periods: x.Periods}
				if x.HasStartIndex {
					n := int(x.StartIndex)
					lecture.StartIndex = &n
				}
				if x.HasEndIndex {
					n := int(x.EndIndex)
					lecture.EndIndex = &n
				}
				actual.Lectures = append(actual.Lectures, lecture)
			}
			if tc.Course == nil || !reflect.DeepEqual(actual, *tc.Course) {
				t.Fatalf("course=%+v, want %+v", actual, tc.Course)
			}
		})
	}
}

func TestGraduateAppCourseIDBinding(t *testing.T) {
	for _, tc := range readStaticBindingFixtures(t).Matches {
		t.Run(tc.Name, func(t *testing.T) {
			selection, err := ParseUSTCGraduateAppCourses([]byte(tc.AppResponse))
			stage := "app"
			var value *USTCLessonIDs
			if err == nil {
				stage = "courseList"
				value, err = selection.MatchCourseIDs([]byte(tc.CourseListResponse))
			} else if selection != nil {
				t.Fatal("failed app parsing returned a selection")
			}
			if (err != nil) != tc.Error {
				t.Fatalf("error=%v, want error=%v", err, tc.Error)
			}
			if tc.Error {
				if stage != tc.ErrorStage {
					t.Fatalf("error stage=%q, want %q", stage, tc.ErrorStage)
				}
				if value != nil {
					t.Fatal("failed parsing returned IDs")
				}
				return
			}
			actual := make([]int, value.Count())
			for i := range actual {
				actual[i] = int(value.At(int64(i)))
			}
			if !reflect.DeepEqual(actual, tc.IDs) {
				t.Fatalf("IDs=%v, want %v", actual, tc.IDs)
			}
		})
	}
}
