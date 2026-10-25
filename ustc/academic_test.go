package ustc_test

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/Life-USTC/client-core/ustc"
)

func TestAcademicSharedFixtures(t *testing.T) {
	data, err := os.ReadFile("testdata/academic_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures struct {
		ScoreSemestersURL, ScoreTypesURL, ExamURL, NumericScoreURL, ScalarScoreURL string
		Scores                                                                     []struct {
			Name, Format, Response string
			SemesterNames          map[int]string
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
	t.Run("request URLs", func(t *testing.T) {
		if ustc.ScoreSemestersURL != fixtures.ScoreSemestersURL || ustc.ScoreTypesURL != fixtures.ScoreTypesURL || ustc.ExamArrangementURL != fixtures.ExamURL {
			t.Fatal("academic endpoint changed")
		}
		if got := ustc.NumericScoreListURL(1, []int{301, 9}); got != fixtures.NumericScoreURL {
			t.Fatalf("numeric request = %q, want %q", got, fixtures.NumericScoreURL)
		}
		if got := ustc.ScalarScoreListURL(1); got != fixtures.ScalarScoreURL {
			t.Fatalf("scalar request = %q, want %q", got, fixtures.ScalarScoreURL)
		}
	})
	for _, tc := range fixtures.Scores {
		t.Run(tc.Name, func(t *testing.T) {
			var report ustc.ScoreReport
			var err error
			switch tc.Format {
			case "numeric":
				report, err = ustc.ParseNumericScores([]byte(tc.Response), tc.SemesterNames)
			case "scalar":
				report, err = ustc.ParseScalarScores([]byte(tc.Response))
			default:
				t.Fatalf("unknown score profile %q", tc.Format)
			}
			if (err != nil) != tc.Error {
				t.Fatalf("error = %v, want error=%v", err, tc.Error)
			}
			if tc.Error {
				return
			}
			var summary, expectedSummary any
			if len(report.Summary) > 0 {
				if err := json.Unmarshal(report.Summary, &summary); err != nil {
					t.Fatal(err)
				}
			}
			if len(tc.Expected.Summary) > 0 {
				if err := json.Unmarshal(tc.Expected.Summary, &expectedSummary); err != nil {
					t.Fatal(err)
				}
			}
			if !reflect.DeepEqual(summary, expectedSummary) {
				t.Fatalf("summary = %#v, want %#v", summary, expectedSummary)
			}
			report.Summary, tc.Expected.Summary = nil, nil
			if !reflect.DeepEqual(report, tc.Expected) {
				got, _ := json.Marshal(report)
				want, _ := json.Marshal(tc.Expected)
				t.Fatalf("report = %s, want %s", got, want)
			}
		})
	}
	for _, tc := range fixtures.Exams {
		t.Run(tc.Name, func(t *testing.T) {
			data := []byte(tc.Response)
			if tc.ResponseBase64 != "" {
				var err error
				data, err = base64.StdEncoding.DecodeString(tc.ResponseBase64)
				if err != nil {
					t.Fatal(err)
				}
			}
			var exams []ustc.Exam
			var err error
			switch tc.Format {
			case "status":
				exams, err = ustc.ParseExamStatusTable(data)
			case "room":
				exams, err = ustc.ParseExamRoomTable(data)
			default:
				t.Fatalf("unknown exam profile %q", tc.Format)
			}
			if (err != nil) != tc.Error {
				t.Fatalf("error = %v, want error=%v", err, tc.Error)
			}
			if !tc.Error && !reflect.DeepEqual(exams, tc.Expected) {
				t.Fatalf("exams = %#v, want %#v", exams, tc.Expected)
			}
		})
	}
}
