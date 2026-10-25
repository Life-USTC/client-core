package mobile

import "github.com/Life-USTC/client-core/ustc"

func USTCScalarScoreListURL(trainTypeID int64) string {
	return ustc.ScalarScoreListURL(int(trainTypeID))
}

type USTCScore struct {
	SemesterID, SemesterName, CourseName, LessonCode, CourseCode string
	Credits, GradePoint                                          float64
	HasGradePoint                                                bool
	Score                                                        string
}

type USTCScoreReport struct {
	GPA                            float64
	HasGPA                         bool
	MajorName                      string
	HasMajorName                   bool
	MajorRank, MajorStdCount       int64
	HasMajorRank, HasMajorStdCount bool
	items                          []*USTCScore
}

func (r *USTCScoreReport) Count() int64 { return int64(len(r.items)) }

func (r *USTCScoreReport) At(index int64) *USTCScore { return r.items[index] }

// ParseUSTCScalarScores converts the existing iOS score profile to binding values.
func ParseUSTCScalarScores(data []byte) (*USTCScoreReport, error) {
	report, err := ustc.ParseScalarScores(data)
	if err != nil {
		return nil, err
	}
	result := &USTCScoreReport{}
	if report.Overview.GPA != nil {
		result.GPA, result.HasGPA = *report.Overview.GPA, true
	}
	if report.Rank != nil {
		if report.Rank.MajorName != nil {
			result.MajorName, result.HasMajorName = *report.Rank.MajorName, true
		}
		if report.Rank.MajorRank != nil {
			result.MajorRank, result.HasMajorRank = int64(*report.Rank.MajorRank), true
		}
		if report.Rank.MajorStdCount != nil {
			result.MajorStdCount, result.HasMajorStdCount = int64(*report.Rank.MajorStdCount), true
		}
	}
	for _, item := range report.Items {
		entry := &USTCScore{
			SemesterID: item.SemesterID, SemesterName: item.SemesterName,
			CourseName: item.CourseName, CourseCode: item.CourseCode, LessonCode: item.LessonCode,
			Credits: item.Credits, Score: item.Score,
		}
		if item.GradePoint != nil {
			entry.GradePoint, entry.HasGradePoint = *item.GradePoint, true
		}
		result.items = append(result.items, entry)
	}
	return result, nil
}
