package ustc

import (
	"fmt"
	"regexp"
	"slices"
)

const (
	GraduateAppCoursesURL = "https://app.ustc.edu.cn/xkjg/wap/default/get-index"
	StaticSemestersPath   = "curriculum/semesters.json"
)

func StaticCourseListPath(semesterID string) string {
	return fmt.Sprintf("curriculum/%s/courses.json", semesterID)
}

func StaticCoursePath(semesterID string, courseID int) string {
	return fmt.Sprintf("curriculum/%s/%d.json", semesterID, courseID)
}

// StaticSemester retains epoch seconds for the native Date and current-term policy.
type StaticSemester struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	StartDate float64 `json:"startDate"`
	EndDate   float64 `json:"endDate"`
}

// StaticCourseListItem keeps both raw lesson-code fields: graduate matching uses
// either, while the additional-course view displays only the canonical field.
type StaticCourseListItem struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	CourseCode       string `json:"courseCode"`
	LessonCode       string `json:"lessonCode"`
	LegacyLessonCode string `json:"legacyLessonCode"`
	TeacherName      string `json:"teacherName"`
}

type StaticLecture struct {
	StartDate   float64 `json:"startDate"`
	EndDate     float64 `json:"endDate"`
	Name        string  `json:"name"`
	Location    string  `json:"location"`
	TeacherName string  `json:"teacherName"`
	Periods     float64 `json:"periods"`
	StartIndex  *int    `json:"startIndex"`
	EndIndex    *int    `json:"endIndex"`
}

type StaticCourse struct {
	ID                      int             `json:"id"`
	Name                    string          `json:"name"`
	CourseCode              string          `json:"courseCode"`
	LessonCode              string          `json:"lessonCode"`
	TeacherName             string          `json:"teacherName"`
	DetailText              *string         `json:"detailText"`
	Credit                  float64         `json:"credit"`
	DateTimePlacePersonText *string         `json:"dateTimePlacePersonText"`
	Lectures                []StaticLecture `json:"lectures"`
}

// ParseStaticSemesters preserves the actual iOS helpers' arrayValue, stringValue
// and numeric numberValue conversions, including their per-entry defaults.
func ParseStaticSemesters(data []byte) ([]StaticSemester, error) {
	value, err := decodeScalarJSON(data)
	if err != nil {
		return nil, fmt.Errorf("decode static semesters: %w", err)
	}
	array, _ := value.([]any)
	semesters := make([]StaticSemester, 0, len(array))
	for _, value := range array {
		item := scalarJSONObject(value)
		semesters = append(semesters, StaticSemester{
			ID: scalarJSONText(item["id"]), Name: scalarJSONText(item["name"]),
			StartDate: scalarJSONNumberValue(item["startDate"]), EndDate: scalarJSONNumberValue(item["endDate"]),
		})
	}
	return semesters, nil
}

func ParseStaticCourseList(data []byte) ([]StaticCourseListItem, error) {
	value, err := decodeScalarJSON(data)
	if err != nil {
		return nil, fmt.Errorf("decode static course list: %w", err)
	}
	array, _ := value.([]any)
	courses := make([]StaticCourseListItem, 0, len(array))
	for _, value := range array {
		item := scalarJSONObject(value)
		courses = append(courses, StaticCourseListItem{
			ID: scalarJSONIntValue(item["id"]), Name: scalarJSONText(item["name"]),
			CourseCode: scalarJSONText(item["courseCode"]), LessonCode: scalarJSONText(item["lessonCode"]),
			LegacyLessonCode: scalarJSONText(item["lesson_code"]), TeacherName: scalarJSONText(item["teacherName"]),
		})
	}
	return courses, nil
}

func ParseStaticCourse(data []byte) (StaticCourse, error) {
	value, err := decodeScalarJSON(data)
	if err != nil {
		return StaticCourse{}, fmt.Errorf("decode static course: %w", err)
	}
	item := scalarJSONObject(value)
	course := StaticCourse{
		ID: scalarJSONIntValue(item["id"]), Name: scalarJSONText(item["name"]),
		CourseCode: scalarJSONText(item["courseCode"]), LessonCode: scalarJSONText(item["lessonCode"]),
		TeacherName: scalarJSONText(item["teacherName"]), DetailText: scalarJSONString(item["detailText"]),
		Credit: scalarJSONNumberValue(item["credit"]), DateTimePlacePersonText: scalarJSONString(item["dateTimePlacePersonText"]),
		Lectures: []StaticLecture{},
	}
	lectures, _ := item["lectures"].([]any)
	for _, value := range lectures {
		item := scalarJSONObject(value)
		course.Lectures = append(course.Lectures, StaticLecture{
			StartDate: scalarJSONNumberValue(item["startDate"]), EndDate: scalarJSONNumberValue(item["endDate"]),
			Name: scalarJSONText(item["name"]), Location: scalarJSONText(item["location"]), TeacherName: scalarJSONText(item["teacherName"]),
			Periods: scalarJSONNumberValue(item["periods"]), StartIndex: scalarJSONInt(item["startIndex"]), EndIndex: scalarJSONInt(item["endIndex"]),
		})
	}
	return course, nil
}

// Swift's original dot does not cross newline characters. Keep the first ASCII
// parenthesis pair, including empty captures, whitespace and repeated codes.
var graduateAppCourseCode = regexp.MustCompile(`\(([^\n\r\v\f\x{0085}\x{2028}\x{2029}]*?)\)`)

func ParseGraduateAppCourseCodes(data []byte) ([]string, error) {
	value, err := decodeScalarJSON(data)
	if err != nil {
		return nil, fmt.Errorf("decode graduate app courses: %w", err)
	}
	root := scalarJSONObject(value)
	array, _ := scalarJSONObject(root["d"])["lists"].([]any)
	codes := make([]string, 0, len(array))
	for _, value := range array {
		name := scalarJSONText(scalarJSONObject(value)["course_name"])
		if match := graduateAppCourseCode.FindStringSubmatch(name); match != nil {
			codes = append(codes, match[1])
		}
	}
	return codes, nil
}

// MatchStaticCourses retains static-list order and repeated IDs. Neither code
// field substitutes for the other, and an explicitly captured empty code matches.
func MatchStaticCourses(courses []StaticCourseListItem, lessonCodes []string) []int {
	ids := make([]int, 0, len(courses))
	for _, course := range courses {
		if slices.Contains(lessonCodes, course.LessonCode) || slices.Contains(lessonCodes, course.LegacyLessonCode) {
			ids = append(ids, course.ID)
		}
	}
	return ids
}
