package mobile

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"

	"github.com/Life-USTC/client-core/ustc"
)

func TestAcademicBindingAgainstSharedFixtures(t *testing.T) {
	data, err := os.ReadFile("../ustc/testdata/academic_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures struct {
		ScalarScoreURL, ExamURL string
		Scores                  []struct {
			Name, Format, Response string
			Expected               ustc.ScoreReport
			Error                  bool
		}
		Exams []struct {
			Name, Format, Response, ResponseBase64 string
			Expected                               []ustc.Exam
			Error                                  bool
		}
	}
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	if len(fixtures.Scores) == 0 || len(fixtures.Exams) == 0 {
		t.Fatal("missing academic fixtures")
	}
	if USTCScalarScoreListURL(1) != fixtures.ScalarScoreURL || USTCExamArrangementURL() != fixtures.ExamURL {
		t.Fatal("binding request URL changed")
	}
	for _, tc := range fixtures.Scores {
		if tc.Format != "scalar" {
			continue
		}
		t.Run(tc.Name, func(t *testing.T) {
			report, err := ParseUSTCScalarScores([]byte(tc.Response))
			if (err != nil) != tc.Error {
				t.Fatalf("error = %v, want error=%v", err, tc.Error)
			}
			if tc.Error {
				if report != nil {
					t.Fatal("failed parsing returned a report")
				}
				return
			}
			if report == nil {
				t.Fatal("successful parsing returned nil")
			}
			if report.HasGPA != (tc.Expected.Overview.GPA != nil) {
				t.Fatal("GPA presence changed")
			}
			if tc.Expected.Overview.GPA != nil && report.GPA != *tc.Expected.Overview.GPA {
				t.Fatal("GPA value changed")
			}
			rank := tc.Expected.Rank
			if report.HasMajorName != (rank != nil && rank.MajorName != nil) || report.HasMajorRank != (rank != nil && rank.MajorRank != nil) || report.HasMajorStdCount != (rank != nil && rank.MajorStdCount != nil) {
				t.Fatal("rank field presence changed")
			}
			if rank != nil {
				if rank.MajorName != nil && report.MajorName != *rank.MajorName {
					t.Fatal("major name changed")
				}
				if rank.MajorRank != nil && report.MajorRank != int64(*rank.MajorRank) {
					t.Fatal("major rank changed")
				}
				if rank.MajorStdCount != nil && report.MajorStdCount != int64(*rank.MajorStdCount) {
					t.Fatal("major count changed")
				}
			}
			if report.Count() != int64(len(tc.Expected.Items)) {
				t.Fatal("entry count changed")
			}
			for i, expected := range tc.Expected.Items {
				entry := report.At(int64(i))
				if entry.SemesterID != expected.SemesterID || entry.SemesterName != expected.SemesterName || entry.CourseName != expected.CourseName || entry.CourseCode != expected.CourseCode || entry.LessonCode != expected.LessonCode || entry.Credits != expected.Credits || entry.Score != expected.Score {
					t.Fatalf("entry %d changed: %+v", i, entry)
				}
				if entry.HasGradePoint != (expected.GradePoint != nil) {
					t.Fatal("grade-point presence changed")
				}
				if expected.GradePoint != nil && entry.GradePoint != *expected.GradePoint {
					t.Fatal("grade-point value changed")
				}
			}
		})
	}
	for _, tc := range fixtures.Exams {
		if tc.Format != "room" {
			continue
		}
		t.Run(tc.Name, func(t *testing.T) {
			data := []byte(tc.Response)
			if tc.ResponseBase64 != "" {
				var err error
				data, err = base64.StdEncoding.DecodeString(tc.ResponseBase64)
				if err != nil {
					t.Fatal(err)
				}
			}
			exams, err := ParseUSTCExamRoomTable(data)
			if (err != nil) != tc.Error {
				t.Fatalf("error = %v, want error=%v", err, tc.Error)
			}
			if tc.Error {
				if exams != nil {
					t.Fatal("failed parsing returned exams")
				}
				return
			}
			if exams == nil || exams.Count() != int64(len(tc.Expected)) {
				t.Fatal("exam count changed")
			}
			for i, expected := range tc.Expected {
				exam := exams.At(int64(i))
				if exam.CourseName != expected.CourseName || exam.LessonCode != expected.LessonCode || exam.ExamType != expected.ExamType || exam.DateTime != expected.DateTime || exam.RoomName != expected.RoomName || exam.BuildingName != expected.BuildingName || exam.District != expected.District || exam.Detail != expected.Detail {
					t.Fatalf("exam %d changed: %+v", i, exam)
				}
			}
		})
	}
}
