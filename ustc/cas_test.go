package ustc

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"reflect"
	"testing"
)

type casFixtures struct {
	LoginURLs   map[string]string
	ServiceURLs []struct{ Name, Service, URL string }
	NginxURLs   []struct{ Name, Target, URL string }
	Pages       []struct {
		Name, URL, Response, Error string
		Page                       CASLoginPage
		AlreadyAuthenticated       bool
	}
	Forms []struct {
		Name, Username, Password, Encoded, Error string
		Page                                     CASLoginPage
	}
	Results []struct {
		Name, ExpectedHost, Response, Error string
		URL                                 *string
	}
}

func readCASFixtures(t *testing.T) casFixtures {
	t.Helper()
	data, err := os.ReadFile("testdata/cas_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures casFixtures
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	if len(fixtures.LoginURLs) == 0 || len(fixtures.ServiceURLs) == 0 || len(fixtures.NginxURLs) == 0 || len(fixtures.Pages) == 0 || len(fixtures.Forms) == 0 || len(fixtures.Results) == 0 {
		t.Fatal("missing CAS protocol fixtures")
	}
	return fixtures
}

func TestCASServiceURLs(t *testing.T) {
	fixtures := readCASFixtures(t)
	actual := map[string]string{"jw": CASJWLoginURL, "blackboard": CASBlackboardLoginURL, "catalog": CASCatalogLoginURL}
	if !reflect.DeepEqual(actual, fixtures.LoginURLs) {
		t.Fatalf("login URLs = %v, want %v", actual, fixtures.LoginURLs)
	}
	for _, tc := range fixtures.ServiceURLs {
		t.Run(tc.Name, func(t *testing.T) {
			if got := CASServiceLoginURL(tc.Service); got != tc.URL {
				t.Fatalf("URL = %q, want %q", got, tc.URL)
			}
		})
	}
	for _, tc := range fixtures.NginxURLs {
		t.Run(tc.Name, func(t *testing.T) {
			if got := BlackboardNginxAuthLoginURL(tc.Target); got != tc.URL {
				t.Fatalf("URL = %q, want %q", got, tc.URL)
			}
		})
	}
}

func TestCASLoginPageFixtures(t *testing.T) {
	for _, tc := range readCASFixtures(t).Pages {
		t.Run(tc.Name, func(t *testing.T) {
			finalURL, err := url.Parse(tc.URL)
			if err != nil {
				t.Fatal(err)
			}
			page, err := ParseCASLoginPage([]byte(tc.Response), finalURL)
			assertCASError(t, err, tc.Error)
			if got := errors.Is(err, ErrCASAlreadyAuthenticated); got != tc.AlreadyAuthenticated {
				t.Errorf("already authenticated = %v, want %v", got, tc.AlreadyAuthenticated)
			}
			if page != tc.Page {
				t.Fatalf("page = %+v, want %+v", page, tc.Page)
			}
		})
	}
}

func TestCASLoginFormFixtures(t *testing.T) {
	for _, tc := range readCASFixtures(t).Forms {
		t.Run(tc.Name, func(t *testing.T) {
			form, err := CASLoginForm(tc.Page, tc.Username, tc.Password)
			assertCASError(t, err, tc.Error)
			if tc.Error != "" {
				if form != nil {
					t.Fatal("failed form construction returned a form")
				}
				return
			}
			if got := form.Encode(); got != tc.Encoded {
				t.Fatalf("form = %q, want %q", got, tc.Encoded)
			}
		})
	}
}

func TestCASLoginResultFixtures(t *testing.T) {
	for _, tc := range readCASFixtures(t).Results {
		t.Run(tc.Name, func(t *testing.T) {
			var finalURL *url.URL
			if tc.URL != nil {
				var err error
				finalURL, err = url.Parse(*tc.URL)
				if err != nil {
					t.Fatal(err)
				}
			}
			err := CheckCASLoginResult(finalURL, []byte(tc.Response), tc.ExpectedHost)
			assertCASError(t, err, tc.Error)
		})
	}
}

func assertCASError(t *testing.T, err error, expected string) {
	t.Helper()
	if expected == "" {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	} else if err == nil || err.Error() != expected {
		t.Fatalf("error = %v, want %q", err, expected)
	}
}
