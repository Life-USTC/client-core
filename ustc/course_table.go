// Package ustc shares USTC request construction and response parsing between clients.
package ustc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strconv"
)

// CourseTablePageURL redirects an authenticated student to their course table.
const CourseTablePageURL = "https://jw.ustc.edu.cn/for-std/course-table"

var studentIDPattern = regexp.MustCompile(`/course-table/(?:info/)?(\d+)`)

// CourseTableDataURL identifies the undergraduate table for the selected semester.
// StudentID is the numeric identifier returned by StudentIDFromCourseTableURL.
func CourseTableDataURL(semesterID int, studentID string) string {
	return fmt.Sprintf("%s/get-data?bizTypeId=2&semesterId=%d&dataId=%s", CourseTablePageURL, semesterID, studentID)
}

// StudentIDFromCourseTableURL reads the student identifier from the final
// course-table redirect, rather than from an unrelated number in the URL.
func StudentIDFromCourseTableURL(rawURL string) (string, error) {
	match := studentIDPattern.FindStringSubmatch(rawURL)
	if len(match) != 2 {
		return "", fmt.Errorf("could not extract jw student id from %s", rawURL)
	}
	return match[1], nil
}

// CourseTable contains the selected lesson identifiers and any embedded details.
// Clients choose how to handle an empty table and how to obtain missing details.
type CourseTable struct {
	LessonIDs []int
	Lessons   []Lesson
}

// ParseCourseTable preserves the JW response precedence used by the CLI:
// studentTableVm fields, then top-level fields, then details indexed by lesson ID.
// Identifiers are sorted and deduplicated; embedded details retain their order.
func ParseCourseTable(data []byte) (CourseTable, error) {
	var payload struct {
		StudentTableVM struct {
			LessonIDs []int    `json:"lessonIds"`
			Lessons   []Lesson `json:"lessons"`
		} `json:"studentTableVm"`
		LessonIDs       []int             `json:"lessonIds"`
		Lessons         []Lesson          `json:"lessons"`
		LessonID2Lesson map[string]Lesson `json:"lessonId2Lesson"`
	}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&payload); err != nil {
		return CourseTable{}, fmt.Errorf("decode course table data: %w", err)
	}

	lessonIDs := payload.StudentTableVM.LessonIDs
	if len(lessonIDs) == 0 {
		lessonIDs = payload.LessonIDs
	}
	if len(lessonIDs) == 0 {
		return CourseTable{LessonIDs: []int{}, Lessons: []Lesson{}}, nil
	}

	slices.Sort(lessonIDs)
	lessons := payload.StudentTableVM.Lessons
	if len(lessons) == 0 {
		lessons = payload.Lessons
	}
	if len(lessons) == 0 && len(payload.LessonID2Lesson) > 0 {
		for _, lessonID := range lessonIDs {
			if lesson, ok := payload.LessonID2Lesson[strconv.Itoa(lessonID)]; ok {
				lessons = append(lessons, lesson)
			}
		}
	}
	return CourseTable{LessonIDs: slices.Compact(lessonIDs), Lessons: lessons}, nil
}
