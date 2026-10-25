package ustc

import (
	"encoding/json"
	"strings"
)

// Lesson retains the JW/catalog fields needed by the existing curriculum adapter.
type Lesson struct {
	ID                      int                 `json:"id"`
	Code                    string              `json:"code"`
	DateTimePlacePersonText LocalizedText       `json:"dateTimePlacePersonText"`
	CourseCode              string              `json:"courseCode"`
	CourseName              string              `json:"courseName"`
	Credits                 float64             `json:"credits"`
	Course                  Course              `json:"course"`
	TeacherAssignmentList   []TeacherAssignment `json:"teacherAssignmentList"`
}

// TeacherAssignment supports the teacher fields returned by JW and catalog.
type TeacherAssignment struct {
	CN      string `json:"cn"`
	EN      string `json:"en"`
	Teacher struct {
		NameZh string `json:"nameZh"`
		NameEn string `json:"nameEn"`
	} `json:"teacher"`
}

// Course is the embedded course information associated with a lesson.
type Course struct {
	Code    string  `json:"code"`
	NameZh  string  `json:"nameZh"`
	NameEn  string  `json:"nameEn"`
	CN      string  `json:"cn"`
	EN      string  `json:"en"`
	Credits float64 `json:"credits"`
}

// LocalizedText accepts a plain string or a cn/en object. A nonempty Chinese
// value takes priority over English, matching the existing CLI adapter.
type LocalizedText struct {
	Value string
}

func (l *LocalizedText) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		l.Value = text
		return nil
	}

	var localized struct {
		CN string `json:"cn"`
		EN string `json:"en"`
	}
	if err := json.Unmarshal(data, &localized); err != nil {
		return err
	}
	l.Value = strings.TrimSpace(localized.CN)
	if l.Value == "" {
		l.Value = strings.TrimSpace(localized.EN)
	}
	return nil
}
