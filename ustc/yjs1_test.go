package ustc_test

import (
	"encoding/json"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/Life-USTC/client-core/ustc"
)

func TestYJS1SharedFixtures(t *testing.T) {
	data, err := os.ReadFile("testdata/yjs1_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures struct {
		URLs  map[string]string
		Forms []struct {
			Name, Kind, SemesterCode, UserID string
			Expected                         url.Values
		}
		Parsers []struct {
			Name, Kind, Response, InfoResponse, ProfileResponse, Error, Now string
			SemesterID, OffsetSeconds                                       int
			Expected                                                        json.RawMessage
		}
		UserIDs []struct {
			Name, Response, Expected string
			Error                    bool
		}
	}
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	t.Run("URLs", func(t *testing.T) {
		got := map[string]string{"scheduleApp": ustc.YJS1ScheduleAppURL, "scheduleTerms": ustc.YJS1ScheduleTermURL, "scheduleRows": ustc.YJS1ScheduleRowsURL, "scoreApp": ustc.YJS1ScoreAppURL, "scoreInfo": ustc.YJS1ScoreInfoURL, "scoreProfile": ustc.YJS1ScoreProfileURL, "scoreRows": ustc.YJS1ScoreRowsURL, "examApp": ustc.YJS1ExamAppURL, "examTerms": ustc.YJS1ExamTermURL, "examRows": ustc.YJS1ExamRowsURL, "homeworkApp": ustc.YJS1HomeworkAppURL, "homeworkTerms": ustc.YJS1HomeworkTermURL, "homeworkRows": ustc.YJS1HomeworkRowsURL}
		if !reflect.DeepEqual(got, fixtures.URLs) {
			t.Fatalf("URLs = %v, want %v", got, fixtures.URLs)
		}
	})
	for _, tc := range fixtures.Forms {
		t.Run(tc.Name, func(t *testing.T) {
			var got url.Values
			switch tc.Kind {
			case "scheduleTerms":
				got = ustc.YJS1ScheduleTermsForm()
			case "schedule":
				got = ustc.YJS1ScheduleRowsForm(tc.SemesterCode)
			case "exam":
				got = ustc.YJS1ExamRowsForm(tc.SemesterCode)
			case "score":
				got = ustc.YJS1ScoreRowsForm()
			case "homeworkTerms":
				got = ustc.YJS1HomeworkTermsForm()
			case "homework":
				got = ustc.YJS1HomeworkRowsForm(tc.SemesterCode, tc.UserID)
			default:
				t.Fatalf("unknown form %s", tc.Kind)
			}
			if !reflect.DeepEqual(got, tc.Expected) {
				t.Fatalf("form = %#v, want %#v", got, tc.Expected)
			}
			if got.Encode() != tc.Expected.Encode() {
				t.Fatal("form encoding changed")
			}
		})
	}
	for _, tc := range fixtures.Parsers {
		t.Run(tc.Name, func(t *testing.T) {
			var got any
			var err error
			switch tc.Kind {
			case "terms":
				got, err = ustc.ParseYJS1Terms([]byte(tc.Response))
			case "curriculum":
				got, err = ustc.ParseYJS1Curriculum([]byte(tc.Response), tc.SemesterID)
			case "exams":
				got, err = ustc.ParseYJS1Exams([]byte(tc.Response))
			case "homework":
				now := time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC)
				if tc.Now != "" {
					now, err = time.Parse(time.RFC3339, tc.Now)
					if err != nil {
						t.Fatal(err)
					}
				}
				got, err = ustc.ParseYJS1Homework([]byte(tc.Response), now, time.FixedZone("caller", tc.OffsetSeconds))
			case "scores":
				infoData := tc.InfoResponse
				if infoData == "" {
					infoData = "{}"
				}
				profileData := tc.ProfileResponse
				if profileData == "" {
					profileData = "{}"
				}
				info, e := ustc.ParseYJS1ScoreInfo([]byte(infoData))
				err = e
				if err == nil {
					profile, e := ustc.ParseYJS1ScoreProfile([]byte(profileData))
					err = e
					if err == nil {
						got, err = ustc.ParseYJS1Scores([]byte(tc.Response), info, profile)
					}
				}
			default:
				t.Fatalf("unknown parser %s", tc.Kind)
			}
			if tc.Error != "" {
				if err == nil || (tc.Error != "*" && err.Error() != tc.Error) {
					t.Fatalf("error = %v, want %q", err, tc.Error)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			encoded, e := json.Marshal(got)
			if e != nil {
				t.Fatal(e)
			}
			var actual, expected any
			if e := json.Unmarshal(encoded, &actual); e != nil {
				t.Fatal(e)
			}
			if e := json.Unmarshal(tc.Expected, &expected); e != nil {
				t.Fatal(e)
			}
			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("result = %s, want %s", encoded, tc.Expected)
			}
		})
	}
	for _, tc := range fixtures.UserIDs {
		t.Run(tc.Name, func(t *testing.T) {
			got, err := ustc.ParseYJS1UserID([]byte(tc.Response), ustc.YJS1HomeworkAppURL)
			if tc.Error {
				if err == nil || err.Error() != "could not resolve graduate student ID from "+ustc.YJS1HomeworkAppURL {
					t.Fatalf("error = %v", err)
				}
				return
			}
			if err != nil || got != tc.Expected {
				t.Fatalf("ID = %q, %v; want %q", got, err, tc.Expected)
			}
		})
	}
}
