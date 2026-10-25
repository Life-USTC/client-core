package ustc

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// These targets preserve the CLI's existing service URLs. Native WebKit and
// service-session profiles retain their own endpoints and request order.
const (
	CASJWLoginURL         = "https://id.ustc.edu.cn/cas/login?service=https%3A%2F%2Fjw.ustc.edu.cn%2Fucas-sso%2Flogin"
	CASBlackboardLoginURL = "https://passport.ustc.edu.cn/login?service=https%3A%2F%2Fwww.bb.ustc.edu.cn%2Fwebapps%2Fbb-SSOIntegrationDemo-BBLEARN%2Fexecute%2FauthValidate%2FcustomLogin%3FreturnUrl%3Dhttp%3A%2F%2Fwww.bb.ustc.edu.cn%2Fwebapps%2Fportal%2Fframeset.jsp%26authProviderId%3D_103_1"
	CASCatalogLoginURL    = "https://id.ustc.edu.cn/cas/login?service=https%3A%2F%2Fcatalog.ustc.edu.cn%2F"
)

// ErrCASAlreadyAuthenticated retains the CLI's existing instruction to prime
// the target session when its login-page request has reached a non-login host.
var ErrCASAlreadyAuthenticated = errors.New("already authenticated")

type CASLoginPage struct {
	URL       string
	CryptoKey string
	Execution string
}

// CASServiceLoginURL constructs the modern CAS service request used by YJS1.
func CASServiceLoginURL(service string) string {
	return CASLoginURL + "?service=" + url.QueryEscape(service)
}

func BlackboardNginxAuthLoginURL(target string) string {
	next := hex.EncodeToString([]byte(target))
	service := "http://www.bb.ustc.edu.cn/nginx_auth/login.php?next=" + next
	return LegacyCASLoginURL + "?service=" + url.QueryEscape(service)
}

// ParseCASLoginPage consumes the final response URL supplied by the HTTP client.
// Like the original caller, finalURL must be the non-nil response request URL.
func ParseCASLoginPage(data []byte, finalURL *url.URL) (CASLoginPage, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(data)))
	if err != nil {
		return CASLoginPage{}, fmt.Errorf("parse login page: %w", err)
	}
	page := CASLoginPage{
		URL:       finalURL.String(),
		CryptoKey: strings.TrimSpace(doc.Find("#login-croypto").First().Text()),
		Execution: strings.TrimSpace(doc.Find("#login-page-flowkey").First().Text()),
	}
	if page.CryptoKey == "" || page.Execution == "" {
		if !strings.EqualFold(finalURL.Host, "id.ustc.edu.cn") && !strings.EqualFold(finalURL.Host, "passport.ustc.edu.cn") {
			return page, ErrCASAlreadyAuthenticated
		}
		return CASLoginPage{}, fmt.Errorf("parse login page %s: missing auth fields", finalURL)
	}
	return page, nil
}

// CASLoginForm performs the existing CAS password transformation and constructs
// its username/password form. Credentials are supplied transiently by the caller.
func CASLoginForm(page CASLoginPage, username, password string) (url.Values, error) {
	encryptedPassword, err := encryptCASPassword(password, page.CryptoKey)
	if err != nil {
		return nil, fmt.Errorf("encrypt CAS password: %w", err)
	}
	form := url.Values{}
	form.Set("username", username)
	form.Set("password", encryptedPassword)
	form.Set("type", "UsernamePassword")
	form.Set("_eventId", "submit")
	form.Set("geolocation", "")
	form.Set("execution", page.Execution)
	form.Set("croypto", page.CryptoKey)
	return form, nil
}

// CheckCASLoginResult preserves the CLI's expected-host success check and error
// precedence. Recognizing an OTP challenge still reports the existing error;
// this package does not submit or complete one.
func CheckCASLoginResult(finalURL *url.URL, body []byte, expectedHost string) error {
	if finalURL != nil && strings.EqualFold(finalURL.Host, expectedHost) {
		return nil
	}
	html := string(body)
	if strings.Contains(strings.ToLower(html), "dynamic password") || strings.Contains(strings.ToLower(html), ">otp<") {
		return fmt.Errorf("school login requires an OTP challenge that the HTTP client could not complete")
	}
	if msg := extractCASLoginError(html); msg != "" {
		return fmt.Errorf("school login failed: %s", msg)
	}
	if finalURL == nil {
		return fmt.Errorf("school login did not complete")
	}
	return fmt.Errorf("school login did not reach %s; ended at %s", expectedHost, finalURL.Host)
}

func extractCASLoginError(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	for _, selector := range []string{
		".ant-form-item-explain-error",
		".alert-error",
		".alert-danger",
		".login-error",
		".error",
		".message",
	} {
		if text := strings.TrimSpace(doc.Find(selector).First().Text()); text != "" {
			return strings.Join(strings.Fields(text), " ")
		}
	}
	return ""
}

func encryptCASPassword(password, keyBase64 string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(strings.TrimSpace(keyBase64))
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	plain := casPKCS7Pad([]byte(password), block.BlockSize())
	encrypted := make([]byte, len(plain))
	for offset := 0; offset < len(plain); offset += block.BlockSize() {
		block.Encrypt(encrypted[offset:offset+block.BlockSize()], plain[offset:offset+block.BlockSize()])
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func casPKCS7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	out := make([]byte, len(data)+padding)
	copy(out, data)
	for i := len(data); i < len(out); i++ {
		out[i] = byte(padding)
	}
	return out
}
