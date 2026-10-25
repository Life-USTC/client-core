package mobile

import "github.com/Life-USTC/client-core/ustc"

func USTCGraduateAppCoursesURL() string                 { return ustc.GraduateAppCoursesURL }
func USTCStaticSemestersPath() string                   { return ustc.StaticSemestersPath }
func USTCStaticCourseListPath(semesterID string) string { return ustc.StaticCourseListPath(semesterID) }
func USTCStaticCoursePath(semesterID string, courseID int64) string {
	return ustc.StaticCoursePath(semesterID, int(courseID))
}

type USTCStaticSemester struct {
	Identifier, Name   string
	StartDate, EndDate float64
}

type USTCStaticSemesters struct{ items []*USTCStaticSemester }

func (s *USTCStaticSemesters) Count() int64                       { return int64(len(s.items)) }
func (s *USTCStaticSemesters) At(index int64) *USTCStaticSemester { return s.items[index] }

func ParseUSTCStaticSemesters(data []byte) (*USTCStaticSemesters, error) {
	semesters, err := ustc.ParseStaticSemesters(data)
	if err != nil {
		return nil, err
	}
	result := &USTCStaticSemesters{}
	for _, semester := range semesters {
		result.items = append(result.items, &USTCStaticSemester{
			Identifier: semester.ID, Name: semester.Name, StartDate: semester.StartDate, EndDate: semester.EndDate,
		})
	}
	return result, nil
}

type USTCStaticCourseListItem struct {
	Identifier                                                  int64
	Name, CourseCode, LessonCode, LegacyLessonCode, TeacherName string
}

type USTCStaticCourseList struct{ items []*USTCStaticCourseListItem }

func (c *USTCStaticCourseList) Count() int64                             { return int64(len(c.items)) }
func (c *USTCStaticCourseList) At(index int64) *USTCStaticCourseListItem { return c.items[index] }

func ParseUSTCStaticCourseList(data []byte) (*USTCStaticCourseList, error) {
	courses, err := ustc.ParseStaticCourseList(data)
	if err != nil {
		return nil, err
	}
	result := &USTCStaticCourseList{}
	for _, course := range courses {
		result.items = append(result.items, &USTCStaticCourseListItem{
			Identifier: int64(course.ID), Name: course.Name, CourseCode: course.CourseCode,
			LessonCode: course.LessonCode, LegacyLessonCode: course.LegacyLessonCode, TeacherName: course.TeacherName,
		})
	}
	return result, nil
}

// Optional values use presence flags so zero and empty strings stay distinct
// from missing values through the Objective-C binding.
type USTCStaticLecture struct {
	StartDate, EndDate, Periods float64
	Name, Location, TeacherName string
	StartIndex, EndIndex        int64
	HasStartIndex, HasEndIndex  bool
}

type USTCStaticCourse struct {
	Identifier                                int64
	Name, CourseCode, LessonCode, TeacherName string
	Credit                                    float64
	DetailText, DateTimePlacePersonText       string
	HasDetailText, HasDateTimePlacePersonText bool
	lectures                                  []*USTCStaticLecture
}

func (c *USTCStaticCourse) LectureCount() int64                      { return int64(len(c.lectures)) }
func (c *USTCStaticCourse) LectureAt(index int64) *USTCStaticLecture { return c.lectures[index] }

func ParseUSTCStaticCourse(data []byte) (*USTCStaticCourse, error) {
	course, err := ustc.ParseStaticCourse(data)
	if err != nil {
		return nil, err
	}
	result := &USTCStaticCourse{
		Identifier: int64(course.ID), Name: course.Name, CourseCode: course.CourseCode,
		LessonCode: course.LessonCode, TeacherName: course.TeacherName, Credit: course.Credit,
	}
	if course.DetailText != nil {
		result.DetailText, result.HasDetailText = *course.DetailText, true
	}
	if course.DateTimePlacePersonText != nil {
		result.DateTimePlacePersonText, result.HasDateTimePlacePersonText = *course.DateTimePlacePersonText, true
	}
	for _, lecture := range course.Lectures {
		value := &USTCStaticLecture{
			StartDate: lecture.StartDate, EndDate: lecture.EndDate, Name: lecture.Name,
			Location: lecture.Location, TeacherName: lecture.TeacherName, Periods: lecture.Periods,
		}
		if lecture.StartIndex != nil {
			value.StartIndex, value.HasStartIndex = int64(*lecture.StartIndex), true
		}
		if lecture.EndIndex != nil {
			value.EndIndex, value.HasEndIndex = int64(*lecture.EndIndex), true
		}
		result.lectures = append(result.lectures, value)
	}
	return result, nil
}

// USTCGraduateAppCourses retains only parsed codes. The immutable handle lets
// native code validate the app response before requesting the static course list.
type USTCGraduateAppCourses struct {
	codes []string
}

func ParseUSTCGraduateAppCourses(data []byte) (*USTCGraduateAppCourses, error) {
	codes, err := ustc.ParseGraduateAppCourseCodes(data)
	if err != nil {
		return nil, err
	}
	return &USTCGraduateAppCourses{codes: codes}, nil
}

func (c *USTCGraduateAppCourses) MatchCourseIDs(courseListData []byte) (*USTCLessonIDs, error) {
	courses, err := ustc.ParseStaticCourseList(courseListData)
	if err != nil {
		return nil, err
	}
	return &USTCLessonIDs{values: ustc.MatchStaticCourses(courses, c.codes)}, nil
}
