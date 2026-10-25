package mobile

import "github.com/Life-USTC/client-core/ustc"

func USTCCASLoginURL() string                  { return ustc.CASLoginURL }
func USTCLegacyCASLoginURL() string            { return ustc.LegacyCASLoginURL }
func USTCNativeCASCookieURL() string           { return ustc.NativeCASCookieURL }
func USTCNativeAASSsoURL() string              { return ustc.NativeAASSsoURL }
func USTCNativeAASCASLoginURL() string         { return ustc.NativeAASCASLoginURL }
func USTCNativeAASRootURL() string             { return ustc.NativeAASRootURL }
func USTCNativeAASHomeURL() string             { return ustc.NativeAASHomeURL }
func USTCNativeBlackboardCASLoginURL() string  { return ustc.NativeBlackboardCASLoginURL }
func USTCNativeBlackboardHomeURL() string      { return ustc.NativeBlackboardHomeURL }
func USTCNativeGraduateAppSessionURL() string  { return ustc.NativeGraduateAppSessionURL }
func USTCNativeGraduateAppCASLoginURL() string { return ustc.NativeGraduateAppCASLoginURL }

func USTCNativeCASIsBrowserSuccess(host, absoluteURL string) bool {
	return ustc.NativeCASIsBrowserSuccess(host, absoluteURL)
}

func USTCNativeCASIsLoginPage(host, path string) bool {
	return ustc.NativeCASIsLoginPage(host, path)
}

// USTCNativeBlackboardLink carries an unvalidated candidate back to Foundation.
type USTCNativeBlackboardLink struct {
	Target string
}

func ParseUSTCNativeBlackboardLoginLink(data []byte) (*USTCNativeBlackboardLink, error) {
	link, err := ustc.ParseNativeBlackboardLoginLink(data)
	if err != nil || link == nil {
		return nil, err
	}
	return &USTCNativeBlackboardLink{Target: link.Target}, nil
}

func ParseUSTCNativeBlackboardFinalLink(data []byte) (*USTCNativeBlackboardLink, error) {
	link, err := ustc.ParseNativeBlackboardFinalLink(data)
	if err != nil || link == nil {
		return nil, err
	}
	return &USTCNativeBlackboardLink{Target: link.Target}, nil
}
