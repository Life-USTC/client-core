package mobile

import "github.com/Life-USTC/client-core/ustc"

func USTCExamArrangementURL() string { return ustc.ExamArrangementURL }

type USTCExam struct {
	CourseName, LessonCode, ExamType, DateTime string
	RoomName, BuildingName, District, Detail   string
}

type USTCExams struct{ items []*USTCExam }

func (e *USTCExams) Count() int64 { return int64(len(e.items)) }

func (e *USTCExams) At(index int64) *USTCExam { return e.items[index] }

// ParseUSTCExamRoomTable binds the iOS room-column profile without converting dates.
func ParseUSTCExamRoomTable(data []byte) (*USTCExams, error) {
	exams, err := ustc.ParseExamRoomTable(data)
	if err != nil {
		return nil, err
	}
	result := &USTCExams{}
	for _, exam := range exams {
		result.items = append(result.items, &USTCExam{
			CourseName: exam.CourseName, LessonCode: exam.LessonCode, ExamType: exam.ExamType, DateTime: exam.DateTime,
			RoomName: exam.RoomName, BuildingName: exam.BuildingName, District: exam.District, Detail: exam.Detail,
		})
	}
	return result, nil
}
