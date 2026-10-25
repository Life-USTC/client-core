package mobile

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

type youngQueryFixture struct {
	Search, Category, Module, ActivityLevel, OrganizerID string
	Active, DateUnknown                                  *bool
	DateFrom, DateTo, TimeBasis                          string
}

func (q youngQueryFixture) binding() *YoungEventQuery {
	result := NewYoungEventQuery()
	result.Search, result.Category, result.Module = q.Search, q.Category, q.Module
	result.ActivityLevel, result.OrganizerID = q.ActivityLevel, q.OrganizerID
	result.DateFrom, result.DateTo, result.TimeBasis = q.DateFrom, q.DateTo, q.TimeBasis
	if q.Active != nil {
		result.HasActive, result.Active = true, *q.Active
	}
	if q.DateUnknown != nil {
		result.HasDateUnknown, result.DateUnknown = true, *q.DateUnknown
	}
	return result
}

func TestYoungQueryBindingAgainstSharedFixtures(t *testing.T) {
	data, err := os.ReadFile("../young/testdata/query_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures struct {
		Requests []struct {
			Name       string
			Query      youngQueryFixture
			Page       *int64
			PageSize   *int64
			Parameters map[string]string
		}
		UnknownDates []struct {
			Name          string
			Query         youngQueryFixture
			ExpectedQuery youngQueryFixture
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
			var page, pageSize int64
			if tc.Page != nil {
				page = *tc.Page
			}
			if tc.PageSize != nil {
				pageSize = *tc.PageSize
			}
			got := tc.Query.binding().Parameters(page, tc.Page != nil, pageSize, tc.PageSize != nil)
			assertYoungParameters(t, got, tc.Parameters)
		})
	}
	for _, tc := range fixtures.UnknownDates {
		t.Run("unknown-dates/"+tc.Name, func(t *testing.T) {
			original := tc.Query.binding()
			got := original.ForUnknownDates()
			if !reflect.DeepEqual(got, tc.ExpectedQuery.binding()) {
				t.Fatalf("ForUnknownDates = %+v, want %+v", got, tc.ExpectedQuery.binding())
			}
			assertYoungParameters(t, got.Parameters(0, false, 0, false), tc.Parameters)
			if !reflect.DeepEqual(original, tc.Query.binding()) {
				t.Fatal("ForUnknownDates changed the original binding")
			}
		})
	}
}

func assertYoungParameters(t *testing.T, got *YoungEventParameters, want map[string]string) {
	t.Helper()
	values := make(map[string]string, got.Count())
	var previous string
	for i := int64(0); i < got.Count(); i++ {
		name := got.Name(i)
		if i > 0 && name <= previous {
			t.Fatalf("unordered or duplicate parameter %q after %q", name, previous)
		}
		previous = name
		values[name] = got.Value(i)
	}
	if !reflect.DeepEqual(values, want) {
		t.Fatalf("Parameters = %#v, want %#v", values, want)
	}
}
