package mobile

import "github.com/Life-USTC/client-core/ustc"

func USTCCourseTablePageURL() string {
	return ustc.CourseTablePageURL
}

func USTCCourseTableDataURL(semesterID int64, studentID string) string {
	return ustc.CourseTableDataURL(int(semesterID), studentID)
}

func USTCStudentIDFromCourseTableURL(rawURL string) (string, error) {
	return ustc.StudentIDFromCourseTableURL(rawURL)
}

// USTCLessonIDs exposes the parsed identifiers needed by the iOS curriculum import.
type USTCLessonIDs struct {
	values []int
}

func (ids *USTCLessonIDs) Count() int64 {
	return int64(len(ids.values))
}

func (ids *USTCLessonIDs) At(index int64) int64 {
	return int64(ids.values[index])
}

func ParseUSTCCourseTable(data []byte) (*USTCLessonIDs, error) {
	table, err := ustc.ParseCourseTable(data)
	if err != nil {
		return nil, err
	}
	return &USTCLessonIDs{values: table.LessonIDs}, nil
}
