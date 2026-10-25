package ustc

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

type nativeCASFixtures struct {
	Endpoints       map[string]string
	BlackboardLinks []struct {
		Name, Response, ResponseBase64 string
		LoginTarget, FinalTarget       *string
	}
	BrowserURLs []struct {
		Name, Host, Path, AbsoluteURL string
		Success, LoginPage            bool
	}
}

func loadNativeCASFixtures(t *testing.T) nativeCASFixtures {
	t.Helper()
	data, err := os.ReadFile("testdata/cas_native_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures nativeCASFixtures
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	if len(fixtures.Endpoints) == 0 || len(fixtures.BlackboardLinks) == 0 || len(fixtures.BrowserURLs) == 0 {
		t.Fatal("missing native CAS fixtures")
	}
	return fixtures
}

func TestNativeCASEndpoints(t *testing.T) {
	got := map[string]string{
		"casLoginURL": CASLoginURL, "legacyCASLoginURL": LegacyCASLoginURL,
		"cookieURL": NativeCASCookieURL, "aasSSOURL": NativeAASSsoURL,
		"aasCASLoginURL": NativeAASCASLoginURL, "aasRootURL": NativeAASRootURL,
		"aasHomeURL": NativeAASHomeURL, "blackboardCASLoginURL": NativeBlackboardCASLoginURL,
		"blackboardHomeURL": NativeBlackboardHomeURL, "graduateAppSessionURL": NativeGraduateAppSessionURL,
		"graduateAppCASLoginURL": NativeGraduateAppCASLoginURL,
	}
	if want := loadNativeCASFixtures(t).Endpoints; !reflect.DeepEqual(got, want) {
		t.Fatalf("native CAS endpoints = %#v, want %#v", got, want)
	}
}

func TestNativeCASBrowserPredicates(t *testing.T) {
	for _, fixture := range loadNativeCASFixtures(t).BrowserURLs {
		t.Run(fixture.Name, func(t *testing.T) {
			if got := NativeCASIsBrowserSuccess(fixture.Host, fixture.AbsoluteURL); got != fixture.Success {
				t.Errorf("browser success = %v, want %v", got, fixture.Success)
			}
			if got := NativeCASIsLoginPage(fixture.Host, fixture.Path); got != fixture.LoginPage {
				t.Errorf("login page = %v, want %v", got, fixture.LoginPage)
			}
		})
	}
}

func TestNativeBlackboardLinks(t *testing.T) {
	for _, fixture := range loadNativeCASFixtures(t).BlackboardLinks {
		t.Run(fixture.Name, func(t *testing.T) {
			data := []byte(fixture.Response)
			if fixture.ResponseBase64 != "" {
				var err error
				data, err = base64.StdEncoding.DecodeString(fixture.ResponseBase64)
				if err != nil {
					t.Fatal(err)
				}
			}
			stages := []struct {
				name  string
				parse func([]byte) (*NativeBlackboardLink, error)
				want  *string
			}{
				{"login", ParseNativeBlackboardLoginLink, fixture.LoginTarget},
				{"final", ParseNativeBlackboardFinalLink, fixture.FinalTarget},
			}
			for _, stage := range stages {
				t.Run(stage.name, func(t *testing.T) {
					got, err := stage.parse(data)
					if err != nil {
						t.Fatal(err)
					}
					if stage.want == nil {
						if got != nil {
							t.Fatalf("link = %+v, want nil", got)
						}
						return
					}
					if got == nil || got.Target != *stage.want {
						t.Fatalf("link = %+v, want %q", got, *stage.want)
					}
				})
			}
		})
	}
}
