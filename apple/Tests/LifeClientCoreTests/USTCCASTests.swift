import Foundation
import LifeCore
import XCTest

final class USTCCASTests: XCTestCase {
    private struct Fixtures: Decodable {
        struct Link: Decodable {
            let name: String
            let response: String?
            let responseBase64: String?
            let loginTarget: String?
            let finalTarget: String?
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
            let host: String
            let path: String
            let absoluteURL: String
            let success: Bool
            let loginPage: Bool
        }
        let endpoints: [String: String]
        let blackboardLinks: [Link]
        let browserURLs: [BrowserURL]
    }

    private func fixtures() throws -> Fixtures {
        let url = try XCTUnwrap(Bundle.module.url(
            forResource: "ustc_cas_native_cases", withExtension: "json", subdirectory: "Fixtures"
        ))
        return try JSONDecoder().decode(Fixtures.self, from: Data(contentsOf: url))
    }

    func testNativeServiceEndpointsThroughCompiledGoBinding() throws {
        let actual = [
            "casLoginURL": LifeCoreMobileUSTCCASLoginURL(),
            "legacyCASLoginURL": LifeCoreMobileUSTCLegacyCASLoginURL(),
            "cookieURL": LifeCoreMobileUSTCNativeCASCookieURL(),
            "aasSSOURL": LifeCoreMobileUSTCNativeAASSsoURL(),
            "aasCASLoginURL": LifeCoreMobileUSTCNativeAASCASLoginURL(),
            "aasRootURL": LifeCoreMobileUSTCNativeAASRootURL(),
            "aasHomeURL": LifeCoreMobileUSTCNativeAASHomeURL(),
            "blackboardCASLoginURL": LifeCoreMobileUSTCNativeBlackboardCASLoginURL(),
            "blackboardHomeURL": LifeCoreMobileUSTCNativeBlackboardHomeURL(),
            "graduateAppSessionURL": LifeCoreMobileUSTCNativeGraduateAppSessionURL(),
            "graduateAppCASLoginURL": LifeCoreMobileUSTCNativeGraduateAppCASLoginURL(),
        ]
        let expected = try fixtures().endpoints
        XCTAssertFalse(expected.isEmpty)
        XCTAssertEqual(actual, expected)
    }

    func testFoundationURLFieldsAndBrowserPredicatesThroughCompiledGoBinding() throws {
        let cases = try fixtures().browserURLs
        XCTAssertFalse(cases.isEmpty)
        for fixture in cases {
            let url = try XCTUnwrap(URL(string: fixture.input), fixture.name)
            XCTAssertEqual(url.host ?? "", fixture.host, fixture.name)
            XCTAssertEqual(url.path, fixture.path, fixture.name)
            XCTAssertEqual(url.absoluteString, fixture.absoluteURL, fixture.name)
            XCTAssertEqual(
                LifeCoreMobileUSTCNativeCASIsBrowserSuccess(url.host ?? "", url.absoluteString),
                fixture.success, fixture.name
            )
            XCTAssertEqual(
                LifeCoreMobileUSTCNativeCASIsLoginPage(url.host ?? "", url.path),
                fixture.loginPage, fixture.name
            )
        }
    }

    func testFirstAnchorCandidatesThroughCompiledGoBinding() throws {
        let cases = try fixtures().blackboardLinks
        XCTAssertFalse(cases.isEmpty)
        for fixture in cases {
            let data = try fixture.data()
            var error: NSError?
            let login = LifeCoreMobileParseUSTCNativeBlackboardLoginLink(data, &error)
            XCTAssertNil(error, fixture.name)
            XCTAssertEqual(login?.target, fixture.loginTarget, fixture.name)
            XCTAssertEqual(login.flatMap { URL(string: $0.target) }?.absoluteString, fixture.loginURL, fixture.name)
            error = nil
            let final = LifeCoreMobileParseUSTCNativeBlackboardFinalLink(data, &error)
            XCTAssertNil(error, fixture.name)
            XCTAssertEqual(final?.target, fixture.finalTarget, fixture.name)
            XCTAssertEqual(final.flatMap { URL(string: $0.target) }?.absoluteString, fixture.finalURL, fixture.name)
        }
    }
}
