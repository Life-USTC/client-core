package ustc

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"reflect"
	"testing"
)

type catalogFixture[T any] struct {
	Name  string `json:"name"`
	Input string `json:"input"`
	Want  T      `json:"want"`
	Error string `json:"error"`
}

func TestCatalogProtocolFixtures(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("testdata/catalog_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures struct {
		Semesters []catalogFixture[[]UndergraduateSemester] `json:"semesters"`
		Types     []catalogFixture[[]ScoreSheetType]        `json:"types"`
		Lessons   []catalogFixture[CatalogLessonPage]       `json:"lessons"`
		URLs      []struct {
			Name       string `json:"name"`
			SemesterID int    `json:"semesterId"`
			Page       int    `json:"page"`
			PageSize   int    `json:"pageSize"`
			Want       string `json:"want"`
		} `json:"urls"`
		Endpoints map[string]string `json:"endpoints"`
	}
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	for _, fixture := range fixtures.Semesters {
		t.Run("semesters/"+fixture.Name, func(t *testing.T) {
			got, err := ParseUndergraduateSemesters([]byte(fixture.Input))
			checkCatalogFixture(t, fixture, got, err)
		})
	}
	for _, fixture := range fixtures.Types {
		t.Run("types/"+fixture.Name, func(t *testing.T) {
			got, err := ParseScoreSheetTypes([]byte(fixture.Input))
			checkCatalogFixture(t, fixture, got, err)
		})
	}
	for _, fixture := range fixtures.Lessons {
		t.Run("lessons/"+fixture.Name, func(t *testing.T) {
			got, err := ParseCatalogLessonPage([]byte(fixture.Input))
			checkCatalogFixture(t, fixture, got, err)
		})
	}
	for _, fixture := range fixtures.URLs {
		t.Run("urls/"+fixture.Name, func(t *testing.T) {
			if got := CatalogLessonsURL(fixture.SemesterID, fixture.Page, fixture.PageSize); got != fixture.Want {
				t.Fatalf("URL = %q, want %q", got, fixture.Want)
			}
		})
	}
	endpoints := map[string]string{"semesters": CatalogSemestersURL, "scoreSemesters": ScoreSemestersURL, "scoreTypes": ScoreTypesURL}
	if !reflect.DeepEqual(endpoints, fixtures.Endpoints) {
		t.Fatalf("endpoints = %#v, want %#v", endpoints, fixtures.Endpoints)
	}
}

func checkCatalogFixture[T any](t *testing.T, fixture catalogFixture[T], got T, err error) {
	t.Helper()
	if fixture.Error != "" {
		var syntax *json.SyntaxError
		var wrongType *json.UnmarshalTypeError
		matched := (fixture.Error == "syntax" && errors.As(err, &syntax)) ||
			(fixture.Error == "type" && errors.As(err, &wrongType)) ||
			(fixture.Error == "eof" && errors.Is(err, io.EOF))
		if !matched {
			t.Fatalf("error = %T %v, want %s", err, err, fixture.Error)
		}
		return
	}
	if err != nil || !reflect.DeepEqual(got, fixture.Want) {
		t.Fatalf("parsed = %#v, %v; want %#v", got, err, fixture.Want)
	}
}
