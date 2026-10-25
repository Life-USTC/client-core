package mobile

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestUSTCNativeCASBindingAgainstSharedFixtures(t *testing.T) {
	data, err := os.ReadFile("../ustc/testdata/cas_native_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures struct {
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
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	if len(fixtures.Endpoints) == 0 || len(fixtures.BlackboardLinks) == 0 || len(fixtures.BrowserURLs) == 0 {
		t.Fatal("missing native CAS fixtures")
	}
	endpoints := map[string]string{
		"casLoginURL": USTCCASLoginURL(), "legacyCASLoginURL": USTCLegacyCASLoginURL(),
		"cookieURL": USTCNativeCASCookieURL(), "aasSSOURL": USTCNativeAASSsoURL(),
		"aasCASLoginURL": USTCNativeAASCASLoginURL(), "aasRootURL": USTCNativeAASRootURL(),
		"aasHomeURL": USTCNativeAASHomeURL(), "blackboardCASLoginURL": USTCNativeBlackboardCASLoginURL(),
		"blackboardHomeURL": USTCNativeBlackboardHomeURL(), "graduateAppSessionURL": USTCNativeGraduateAppSessionURL(),
		"graduateAppCASLoginURL": USTCNativeGraduateAppCASLoginURL(),
	}
	if !reflect.DeepEqual(endpoints, fixtures.Endpoints) {
		t.Fatalf("endpoints = %#v, want %#v", endpoints, fixtures.Endpoints)
	}
	for _, fixture := range fixtures.BrowserURLs {
		t.Run(fixture.Name, func(t *testing.T) {
			if got := USTCNativeCASIsBrowserSuccess(fixture.Host, fixture.AbsoluteURL); got != fixture.Success {
				t.Errorf("browser success = %v, want %v", got, fixture.Success)
			}
			if got := USTCNativeCASIsLoginPage(fixture.Host, fixture.Path); got != fixture.LoginPage {
				t.Errorf("login page = %v, want %v", got, fixture.LoginPage)
			}
		})
	}
	for _, fixture := range fixtures.BlackboardLinks {
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
				parse func([]byte) (*USTCNativeBlackboardLink, error)
				want  *string
			}{
				{"login", ParseUSTCNativeBlackboardLoginLink, fixture.LoginTarget},
				{"final", ParseUSTCNativeBlackboardFinalLink, fixture.FinalTarget},
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
