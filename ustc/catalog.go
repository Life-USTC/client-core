package ustc

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const CatalogSemestersURL = "https://catalog.ustc.edu.cn/api/teach/semester/list"

// CatalogLessonsURL leaves page selection and page size to the caller.
func CatalogLessonsURL(semesterID, page, pageSize int) string {
	return fmt.Sprintf("https://catalog.ustc.edu.cn/api/teach/lesson/list-for-teach/%d?page=%d&pageSize=%d", semesterID, page, pageSize)
}

// UndergraduateSemester is the existing catalog and JW score-metadata profile.
// Selection, display-name precedence and fallback remain client policies.
type UndergraduateSemester struct {
	ID         int    `json:"id"`
	Code       string `json:"code,omitempty"`
	SemesterCn string `json:"semesterCn,omitempty"`
	SemesterEn string `json:"semesterEn,omitempty"`
	StartDate  string `json:"startDate,omitempty"`
	EndDate    string `json:"endDate,omitempty"`
	IsLast     bool   `json:"isLast,omitempty"`
}

func ParseUndergraduateSemesters(data []byte) ([]UndergraduateSemester, error) {
	var semesters []UndergraduateSemester
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&semesters); err != nil {
		return nil, err
	}
	return semesters, nil
}

type ScoreSheetType struct {
	ID int `json:"id"`
}

func ParseScoreSheetTypes(data []byte) ([]ScoreSheetType, error) {
	var types []ScoreSheetType
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&types); err != nil {
		return nil, err
	}
	return types, nil
}

// CatalogLessonPage distinguishes a direct array/null response from a page.
// A direct response replaces earlier pages in the existing CLI consumer.
type CatalogLessonPage struct {
	Rows      []Lesson `json:"rows"`
	TotalPage int      `json:"totalPage"`
	Direct    bool     `json:"direct"`
}

// ParseCatalogLessonPage requires one complete JSON document, unlike semester
// and transcript metadata. It preserves the existing direct-then-page attempts.
func ParseCatalogLessonPage(data []byte) (CatalogLessonPage, error) {
	var direct []Lesson
	if err := json.Unmarshal(data, &direct); err == nil {
		return CatalogLessonPage{Rows: direct, Direct: true}, nil
	}
	var page struct {
		Rows      []Lesson `json:"rows"`
		TotalPage int      `json:"totalPage"`
	}
	if err := json.Unmarshal(data, &page); err != nil {
		return CatalogLessonPage{}, err
	}
	return CatalogLessonPage{Rows: page.Rows, TotalPage: page.TotalPage}, nil
}
