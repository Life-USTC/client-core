package ustc

import (
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

const (
	CASLoginURL       = "https://id.ustc.edu.cn/cas/login"
	LegacyCASLoginURL = "https://passport.ustc.edu.cn/login"

	NativeCASCookieURL           = "https://id.ustc.edu.cn"
	NativeAASSsoURL              = "https://jw.ustc.edu.cn/ucas-sso/login"
	NativeAASCASLoginURL         = "https://passport.ustc.edu.cn/login?service=https%3A%2F%2Fjw.ustc.edu.cn%2Fucas-sso%2Flogin"
	NativeAASRootURL             = "https://jw.ustc.edu.cn/"
	NativeAASHomeURL             = "https://jw.ustc.edu.cn/home"
	NativeBlackboardHomeURL      = "https://www.bb.ustc.edu.cn/webapps/portal/execute/tabs/tabAction?tab_tab_group_id=_1_1"
	NativeBlackboardCASLoginURL  = "https://passport.ustc.edu.cn/login?service=https%3a%2f%2fwww.bb.ustc.edu.cn%2fwebapps%2fbb-SSOIntegrationDemo-BBLEARN%2fexecute%2fauthValidate%2fcustomLogin%3freturnUrl%3dhttp%3a%2f%2fwww.bb.ustc.edu.cn%2fwebapps%2fportal%2fframeset.jsp%26authProviderId%3d_103_1"
	NativeGraduateAppSessionURL  = "https://app.ustc.edu.cn/a_ustc/api/cas/index?redirect=https://app.ustc.edu.cn/site/examManage/index&from=wap"
	NativeGraduateAppCASLoginURL = "https://id.ustc.edu.cn/cas/login?service=https://app.ustc.edu.cn/a_ustc/api/cas/index?redirect=https%3A%2F%2Fapp.ustc.edu.cn%2Fsite%2FtimeTableQuery%2Findex&from=wap"
)

// NativeCASIsBrowserSuccess receives the original Foundation URL fields. It
// preserves the WebKit delegate's case-sensitive host and absolute URL prefix.
func NativeCASIsBrowserSuccess(host, absoluteURL string) bool {
	return host == "id.ustc.edu.cn" && strings.HasPrefix(absoluteURL, "https://id.ustc.edu.cn/gate/cas-success")
}

// NativeCASIsLoginPage preserves the native page predicate, including its path
// prefix and lack of a scheme restriction. URL parsing remains in Foundation.
func NativeCASIsLoginPage(host, path string) bool {
	return host == "id.ustc.edu.cn" && strings.HasPrefix(path, "/cas/login")
}

// NativeBlackboardLink is a candidate for Foundation URL(string:), not a parsed
// or validated URL. A non-nil link can contain an empty or relative target.
type NativeBlackboardLink struct {
	Target string
}

// ParseNativeBlackboardLoginLink examines only the first anchor. The native
// client follows it only when its href starts with the literal login.php.
func ParseNativeBlackboardLoginLink(data []byte) (*NativeBlackboardLink, error) {
	link, err := ParseNativeBlackboardFinalLink(data)
	if err != nil || link == nil {
		return nil, err
	}
	if !strings.HasPrefix(link.Target, "login.php") {
		return nil, nil
	}
	return &NativeBlackboardLink{Target: "https://www.bb.ustc.edu.cn/nginx_auth/" + link.Target}, nil
}

// ParseNativeBlackboardFinalLink returns the first anchor's href without
// trimming, resolving or validating it. Like String(data:encoding:.utf8) ?? "",
// an invalid UTF-8 response becomes an empty document. A missing href is empty.
func ParseNativeBlackboardFinalLink(data []byte) (*NativeBlackboardLink, error) {
	html := ""
	if utf8.Valid(data) {
		html = string(data)
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	anchor := doc.Find("a").First()
	if anchor.Length() == 0 {
		return nil, nil
	}
	return &NativeBlackboardLink{Target: anchor.AttrOr("href", "")}, nil
}
