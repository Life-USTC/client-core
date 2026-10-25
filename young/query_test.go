package young_test

import (
	"encoding/json"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/Life-USTC/client-core/young"
)

func TestSharedQueryFixtures(t *testing.T) {
	data, err := os.ReadFile("testdata/query_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures struct {
		Requests []struct {
			Name       string
			Query      young.EventQuery
			Page       *int
			PageSize   *int
			Parameters map[string]string
		}
		UnknownDates []struct {
			Name          string
			Query         young.EventQuery
			ExpectedQuery young.EventQuery
			Parameters    map[string]string
		}
	}
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	if len(fixtures.Requests) == 0 || len(fixtures.UnknownDates) == 0 {
		t.Fatal("missing shared Young query fixtures")
	}
	for _, tc := range fixtures.Requests {
		t.Run("parameters/"+tc.Name, func(t *testing.T) {
			assertParameters(t, tc.Query.Parameters(tc.Page, tc.PageSize), tc.Parameters)
		})
	}
	for _, tc := range fixtures.UnknownDates {
		t.Run("unknown-dates/"+tc.Name, func(t *testing.T) {
			before := tc.Query.Parameters(nil, nil)
			got := tc.Query.ForUnknownDates()
			if !reflect.DeepEqual(got, tc.ExpectedQuery) {
				t.Fatalf("ForUnknownDates = %+v, want %+v", got, tc.ExpectedQuery)
			}
			assertParameters(t, got.Parameters(nil, nil), tc.Parameters)
			if !reflect.DeepEqual(tc.Query.Parameters(nil, nil), before) {
				t.Fatal("ForUnknownDates changed the original query")
			}
			if !reflect.DeepEqual(got.ForUnknownDates(), got) {
				t.Fatal("ForUnknownDates changed an already transformed query")
			}
		})
	}
}

func assertParameters(t *testing.T, got url.Values, expected map[string]string) {
	t.Helper()
	want := url.Values{}
	for key, value := range expected {
		want.Set(key, value)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Parameters = %#v, want %#v", got, want)
	}
	decoded, err := url.ParseQuery(got.Encode())
	if err != nil || !reflect.DeepEqual(decoded, want) {
		t.Fatalf("query encoding round trip = %#v, %v; want %#v", decoded, err, want)
	}
}
