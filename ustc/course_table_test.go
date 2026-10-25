package ustc_test

import (
	"encoding/json"
	"os"
	"slices"
	"testing"

	"github.com/Life-USTC/client-core/ustc"
)

func TestCourseTableSharedFixtures(t *testing.T) {
	data, err := os.ReadFile("testdata/course_table_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures struct {
		PageURL  string
		Requests []struct {
			Name, StudentID, URL string
			SemesterID           int
		}
		StudentIDs []struct {
			Name, URL, StudentID string
			Error                bool
		}
		Tables []struct {
			Name, Response               string
			LessonIDs, EmbeddedLessonIDs []int
			LessonCodes                  []string
			Error                        bool
		}
	}
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	if fixtures.PageURL == "" || len(fixtures.Requests) == 0 || len(fixtures.StudentIDs) == 0 || len(fixtures.Tables) == 0 {
		t.Fatal("missing shared USTC course table fixtures")
	}
	if ustc.CourseTablePageURL != fixtures.PageURL {
		t.Fatalf("CourseTablePageURL = %q, want %q", ustc.CourseTablePageURL, fixtures.PageURL)
	}
	for _, tc := range fixtures.Requests {
		t.Run("request/"+tc.Name, func(t *testing.T) {
			if got := ustc.CourseTableDataURL(tc.SemesterID, tc.StudentID); got != tc.URL {
				t.Fatalf("CourseTableDataURL = %q, want %q", got, tc.URL)
			}
		})
	}
	for _, tc := range fixtures.StudentIDs {
		t.Run("student-id/"+tc.Name, func(t *testing.T) {
			got, err := ustc.StudentIDFromCourseTableURL(tc.URL)
			if (err != nil) != tc.Error || got != tc.StudentID {
				t.Fatalf("StudentIDFromCourseTableURL = %q, %v; want %q, error=%v", got, err, tc.StudentID, tc.Error)
			}
		})
	}
	for _, tc := range fixtures.Tables {
		t.Run("response/"+tc.Name, func(t *testing.T) {
			got, err := ustc.ParseCourseTable([]byte(tc.Response))
			if (err != nil) != tc.Error {
				t.Fatalf("ParseCourseTable error = %v, want error=%v", err, tc.Error)
			}
			if tc.Error {
				return
			}
			if !slices.Equal(got.LessonIDs, tc.LessonIDs) {
				t.Fatalf("LessonIDs = %v, want %v", got.LessonIDs, tc.LessonIDs)
			}
			var embeddedIDs []int
			var lessonCodes []string
			for _, lesson := range got.Lessons {
				embeddedIDs = append(embeddedIDs, lesson.ID)
				lessonCodes = append(lessonCodes, lesson.Code)
			}
			if !slices.Equal(embeddedIDs, tc.EmbeddedLessonIDs) || !slices.Equal(lessonCodes, tc.LessonCodes) {
				t.Fatalf("embedded lessons = %v, %v; want %v, %v", embeddedIDs, lessonCodes, tc.EmbeddedLessonIDs, tc.LessonCodes)
			}
		})
	}
}

func TestCourseTableRetainsCurriculumFields(t *testing.T) {
	table, err := ustc.ParseCourseTable([]byte(`{
		"lessonIds": [10],
		"lessons": [{
			"id": 10, "code": "MATH1001.01",
			"courseCode": "MATH1001", "courseName": "高等数学", "credits": 4,
			"dateTimePlacePersonText": {"cn": " 周一 3-4节 ", "en": "Monday"},
			"course": {
				"code": "MATH1001", "nameZh": "高等数学", "nameEn": "Calculus",
				"cn": "数学", "en": "Math", "credits": 4.5
			},
			"teacherAssignmentList": [{
				"cn": "张三", "en": "Zhang San",
				"teacher": {"nameZh": "李四", "nameEn": "Li Si"}
			}]
		}]
	}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(table.Lessons) != 1 {
		t.Fatalf("got %d lessons, want one", len(table.Lessons))
	}
	lesson := table.Lessons[0]
	if lesson.ID != 10 || lesson.Code != "MATH1001.01" || lesson.CourseCode != "MATH1001" || lesson.CourseName != "高等数学" || lesson.Credits != 4 {
		t.Fatalf("lesson fields changed: %+v", lesson)
	}
	if lesson.DateTimePlacePersonText.Value != "周一 3-4节" {
		t.Fatalf("schedule = %q", lesson.DateTimePlacePersonText.Value)
	}
	wantCourse := ustc.Course{Code: "MATH1001", NameZh: "高等数学", NameEn: "Calculus", CN: "数学", EN: "Math", Credits: 4.5}
	if lesson.Course != wantCourse {
		t.Fatalf("embedded course = %+v, want %+v", lesson.Course, wantCourse)
	}
	if len(lesson.TeacherAssignmentList) != 1 {
		t.Fatalf("got %d teachers, want one", len(lesson.TeacherAssignmentList))
	}
	teacher := lesson.TeacherAssignmentList[0]
	if teacher.CN != "张三" || teacher.EN != "Zhang San" || teacher.Teacher.NameZh != "李四" || teacher.Teacher.NameEn != "Li Si" {
		t.Fatalf("teacher fields changed: %+v", teacher)
	}
}

func TestLocalizedLessonText(t *testing.T) {
	for _, tc := range []struct {
		name, input, want string
		error             bool
	}{
		{name: "plain text keeps whitespace", input: `" 周一 3-4节 "`, want: " 周一 3-4节 "},
		{name: "Chinese first", input: `{"cn":" 中文 ","en":"English"}`, want: "中文"},
		{name: "English without Chinese", input: `{"cn":" ","en":" English "}`, want: "English"},
		{name: "empty object", input: `{}`},
		{name: "null text", input: `null`},
		{name: "invalid text", input: `42`, error: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var value ustc.LocalizedText
			err := json.Unmarshal([]byte(tc.input), &value)
			if (err != nil) != tc.error || value.Value != tc.want {
				t.Fatalf("localized text = %q, %v; want %q, error=%v", value.Value, err, tc.want, tc.error)
			}
		})
	}
}
