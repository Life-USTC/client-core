import Foundation
import LifeCore

/// The fields imported by the existing selected-calendar-events homework path.
public struct USTCBlackboardSelectedCalendarEvent: Equatable, Sendable {
    public let title: String
    public let courseName: String
    public let endDate: String

    public init(title: String, courseName: String, endDate: String) {
        self.title = title
        self.courseName = courseName
        self.endDate = endDate
    }
}

public enum USTCBlackboard {
    public static var selectedCalendarEventsURL: URL {
        URL(string: LifeCoreMobileUSTCBlackboardSelectedCalendarEventsURL())!
    }

    /// Dates remain raw strings for the native adapter's existing DateFormatter.
    public static func parseSelectedCalendarEvents(_ data: Data) throws -> [USTCBlackboardSelectedCalendarEvent] {
        var error: NSError?
        let value = LifeCoreMobileParseUSTCBlackboardSelectedCalendarEvents(data, &error)
        if let error { throw error }
        // The Go function returns a collection whenever its error is nil.
        return (0..<value!.count()).map { index in
            USTCBlackboardSelectedCalendarEvent(
                title: value!.title(index),
                courseName: value!.courseName(index),
                endDate: value!.endDate(index)
            )
        }
    }
}
