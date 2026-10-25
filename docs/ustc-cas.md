# Existing CAS and service protocols

[Core #21](https://github.com/Life-USTC/client-core/issues/21) covers the final
portable authentication protocol slice of [Task #12](https://github.com/Life-USTC/client-core/issues/12).
It moves existing university request and response rules. It does not replace
WebKit login with HTTP login or add an OTP flow.

## CLI profile

`CASJWLoginURL`, `CASBlackboardLoginURL` and `CASCatalogLoginURL` preserve the
existing fixed targets, including their modern/legacy CAS hosts and percent
encoding. `CASServiceLoginURL` encodes the YJS1 service URL with the existing Go
query encoder. `BlackboardNginxAuthLoginURL` retains the lowercase hexadecimal
`next` target inside the legacy CAS service request.

`ParseCASLoginPage` receives response bytes and the final HTTP request URL. It
selects the first `#login-croypto` and `#login-page-flowkey` elements, keeps their
descendant text and trims surrounding whitespace. Missing fields on the two
case-insensitively matched login hosts retain the existing parse error. Missing
fields after reaching another host return the parsed page and
`ErrCASAlreadyAuthenticated`, which the existing CLI uses to prime its session.
It deliberately compares `URL.Host`, not `Hostname()`. The caller supplies the
non-nil response URL already required by its existing code.

`CASLoginForm` takes that page and transient username/password strings. It keeps
the existing Base64 key decoding, AES block encryption with PKCS7 padding, and
seven form fields: `username`, `password`, `type`, `_eventId`, `geolocation`,
`execution` and `croypto`. The encoded form and encryption errors match the
original CLI. The package neither loads nor stores credentials. Password
transformation is private; the public function is the form the actual caller
needs.

`CheckCASLoginResult` first accepts the existing case-insensitive expected-host
match. Otherwise it retains the existing OTP-required error before examining
page errors, then checks the ordered error selectors and finally reports the
missing or unexpected final URL. A matching selector's first element supplies
the message, with collapsed whitespace; a blank first element advances to the
next selector. Recognizing an OTP challenge does not submit or complete it.

## Native browser and service profile

The native constants retain the browser login URL, cookie origin, AAS session
URLs, Blackboard legacy URL and completion URL, and the two graduate app session
URLs. Blackboard retains its original lowercase percent escapes, distinct from
the CLI's fixed uppercase-escaped target. No service request is reordered.

`NativeCASIsBrowserSuccess` receives Foundation's host and absolute URL string;
`NativeCASIsLoginPage` receives its host and path. They preserve the existing
case-sensitive comparisons and prefixes. Foundation remains responsible for
URL parsing and representation. The public Swift wrapper keeps native legacy
service markup through `URLComponents.queryItems`, using the shared legacy
endpoint; it does not substitute Go's query encoding.

`ParseNativeBlackboardLoginLink` selects only the first anchor and follows the
existing literal `login.php` prefix rule, adding the existing nginx-auth prefix.
`ParseNativeBlackboardFinalLink` selects the next page's first anchor without
trimming or resolving its href. No anchor returns nil; an anchor with a missing
or empty href returns a present candidate with an empty target. Swift still
performs `URL(string:)` before following a candidate. Invalid UTF-8 bytes retain
the native empty-document behavior. Parsing does not follow any links itself.

The mobile bridge exports only these native endpoint, candidate and predicate
operations. It does not export CLI password/form functions or move credentials
across the Apple binding. WebKit keeps filling and submitting its existing DOM
form; it does not begin using the CLI's HTTP form protocol.

## Adapter ownership

The clients retain HTTP execution, status handling, headers, cookies, session
priming, request order and existing retry behavior. CLI retains credential
loading, command output and synchronization. iOS retains WebKit and DOM
execution, native presentation and navigation, Keychain, cookie transfer and
login continuation management. These functions do not introduce a request
engine, new fallback, host policy or authentication state machine.

## Evidence and release

`ustc/testdata/cas_cases.json` characterizes the existing CLI implementation:
fixed and nested service URLs, first/absent login fields, already-authenticated
classification, form encoding and password transformation, ordered errors and
existing OTP failure. Fixtures use synthetic usernames, passwords and keys.
Ciphertext expectations were generated independently with OpenSSL AES-ECB,
and the complete fixture suite passed a probe running the original CLI pure
functions before the shared implementation was accepted. The probe only adapts
the original login-page parser's HTTP response inputs into supplied bytes/URL.
These are compatibility fixtures, not university-account login evidence.

`ustc/testdata/cas_native_cases.json` records the native first-anchor, URL-field
and legacy-service-markup profiles. A macOS probe used the pinned SwiftSoup
source and Foundation against the original Swift selectors and predicates;
12 anchor inputs, 10 URL projections and 4 service-markup inputs matched.
Native URL construction remains Foundation-owned, including rejecting an empty
candidate with `URL(string:)`. Mobile and raw compiled Swift tests consume the
same native fixture. The separate wrapper draft is applied only with the new
released framework in the distribution change.

Ship source and compiled raw binding tests before publishing the matching
framework. Add public Swift wrappers with the archive URL/checksum in the
subsequent distribution change. Both clients must adopt the released dependency
in their actual login paths and pass their regression suites before this slice
and the parent Task can be complete. Production server OAuth and TestFlight
agreement acceptance remain separate.
