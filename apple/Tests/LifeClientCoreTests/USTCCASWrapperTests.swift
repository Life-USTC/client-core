import Foundation
import LifeClientCore
import XCTest

final class USTCCASWrapperTests: XCTestCase {
    private struct Fixtures: Decodable {
        struct Link: Decodable {
            let name: String
            let response: String?
            let responseBase64: String?
            let loginURL: String?
            let finalURL: String?

            func data() throws -> Data {
                if let response { return Data(response.utf8) }
                let encoded = try XCTUnwrap(responseBase64)
                return try XCTUnwrap(Data(base64Encoded: encoded))
            }
        }
        struct BrowserURL: Decodable {
            let name: String
            let input: String
            let success: Bool
            let loginPage: Bool
        }
        struct LegacyMarkup: Decodable {
            let name: String
            let service: String
            let serviceAbsoluteURL: String
            let url: String
        }
        let endpoints: [String: String]
        let blackboardLinks: [Link]
        let browserURLs: [BrowserURL]
        let legacyMarkup: [LegacyMarkup]
    }

    private func fixtures() throws -> Fixtures {
        let url = try XCTUnwrap(Bundle.module.url(
            forResource: "ustc_cas_native_cases", withExtension: "json", subdirectory: "Fixtures"
        ))
        return try JSONDecoder().decode(Fixtures.self, from: Data(contentsOf: url))
    }

    func testNativeServiceURLsPreserveTheirOriginalSpelling() throws {
        let actual = [
            "casLoginURL": USTCCAS.loginURL,
            "legacyCASLoginURL": USTCCAS.legacyLoginURL,
            "cookieURL": USTCCAS.cookieURL,
            "aasSSOURL": USTCCAS.aasSSOURL,
            "aasCASLoginURL": USTCCAS.aasCASLoginURL,
            "aasRootURL": USTCCAS.aasRootURL,
            "aasHomeURL": USTCCAS.aasHomeURL,
            "blackboardCASLoginURL": USTCCAS.blackboardCASLoginURL,
            "blackboardHomeURL": USTCCAS.blackboardHomeURL,
            "graduateAppSessionURL": USTCCAS.graduateAppSessionURL,
            "graduateAppCASLoginURL": USTCCAS.graduateAppCASLoginURL,
        ].mapValues(\.absoluteString)
        XCTAssertEqual(actual, try fixtures().endpoints)
    }

    func testBrowserPredicatesUseFoundationHostPathAndAbsoluteString() throws {
        let cases = try fixtures().browserURLs
        XCTAssertFalse(cases.isEmpty)
        for fixture in cases {
            let url = try XCTUnwrap(URL(string: fixture.input), fixture.name)
            XCTAssertEqual(USTCCAS.isBrowserSuccess(url), fixture.success, fixture.name)
            XCTAssertEqual(USTCCAS.isBrowserLoginPage(url), fixture.loginPage, fixture.name)
        }
    }

    func testBlackboardLinksKeepFirstAnchorAndFoundationURLBehavior() throws {
        let cases = try fixtures().blackboardLinks
        let selected = [
            "login-query-entities", "first-anchor-only", "href-is-not-trimmed",
            "relative-final-link", "missing-href-is-empty", "invalid-foundation-url-candidate",
            "invalid-utf8-discards-entire-document",
        ]
        for name in selected {
            let fixture = try XCTUnwrap(cases.first { $0.name == name }, name)
            let data = try fixture.data()
            XCTAssertEqual(try USTCCAS.blackboardLoginURL(from: data)?.absoluteString, fixture.loginURL, name)
            XCTAssertEqual(try USTCCAS.blackboardFinalURL(from: data)?.absoluteString, fixture.finalURL, name)
        }
    }

    func testLegacyServiceMarkupKeepsURLComponentsEncoding() throws {
        let cases = try fixtures().legacyMarkup
        XCTAssertFalse(cases.isEmpty)
        for fixture in cases {
            let service = try XCTUnwrap(URL(string: fixture.service), fixture.name)
            XCTAssertEqual(service.absoluteString, fixture.serviceAbsoluteURL, fixture.name)
            let url = USTCCAS.legacyServiceLoginURL(for: service)
            XCTAssertEqual(url.absoluteString, fixture.url, fixture.name)
            let components = try XCTUnwrap(URLComponents(url: url, resolvingAgainstBaseURL: false))
            XCTAssertEqual(components.queryItems, [URLQueryItem(name: "service", value: service.absoluteString)])
        }
    }
}
