import Foundation
import LifeCore

/// The existing native CAS and service-session protocol profile.
/// Networking, cookies, WebKit and login decisions remain in the caller.
public enum USTCCAS {
    public static var loginURL: URL { URL(string: LifeCoreMobileUSTCCASLoginURL())! }
    public static var legacyLoginURL: URL { URL(string: LifeCoreMobileUSTCLegacyCASLoginURL())! }
    public static var cookieURL: URL { URL(string: LifeCoreMobileUSTCNativeCASCookieURL())! }
    public static var aasSSOURL: URL { URL(string: LifeCoreMobileUSTCNativeAASSsoURL())! }
    public static var aasCASLoginURL: URL { URL(string: LifeCoreMobileUSTCNativeAASCASLoginURL())! }
    public static var aasRootURL: URL { URL(string: LifeCoreMobileUSTCNativeAASRootURL())! }
    public static var aasHomeURL: URL { URL(string: LifeCoreMobileUSTCNativeAASHomeURL())! }
    public static var blackboardCASLoginURL: URL { URL(string: LifeCoreMobileUSTCNativeBlackboardCASLoginURL())! }
    public static var blackboardHomeURL: URL { URL(string: LifeCoreMobileUSTCNativeBlackboardHomeURL())! }
    public static var graduateAppSessionURL: URL { URL(string: LifeCoreMobileUSTCNativeGraduateAppSessionURL())! }
    public static var graduateAppCASLoginURL: URL { URL(string: LifeCoreMobileUSTCNativeGraduateAppCASLoginURL())! }

    public static func isBrowserSuccess(_ url: URL) -> Bool {
        LifeCoreMobileUSTCNativeCASIsBrowserSuccess(url.host ?? "", url.absoluteString)
    }

    public static func isBrowserLoginPage(_ url: URL) -> Bool {
        LifeCoreMobileUSTCNativeCASIsLoginPage(url.host ?? "", url.path)
    }

    public static func legacyServiceLoginURL(for service: URL) -> URL {
        var components = URLComponents(url: legacyLoginURL, resolvingAgainstBaseURL: false)!
        components.queryItems = [.init(name: "service", value: service.absoluteString)]
        return components.url!
    }

    public static func blackboardLoginURL(from data: Data) throws -> URL? {
        var error: NSError?
        let link = LifeCoreMobileParseUSTCNativeBlackboardLoginLink(data, &error)
        if let error { throw error }
        return link.flatMap { URL(string: $0.target) }
    }

    public static func blackboardFinalURL(from data: Data) throws -> URL? {
        var error: NSError?
        let link = LifeCoreMobileParseUSTCNativeBlackboardFinalLink(data, &error)
        if let error { throw error }
        return link.flatMap { URL(string: $0.target) }
    }
}
