import Foundation
import LifeCore

public struct CivilDateRange: Codable, Equatable {
    public let from: String
    public let to: String
}

public enum CalendarPeriod: String {
    case day, week, month
}

public enum CalendarWeekday: Int64 {
    case sunday, monday, tuesday, wednesday, thursday, friday, saturday
}

/// Swift value conversion for the shared Go calendar implementation.
public enum ClientCalendar {
    public static func bounds(
        period: CalendarPeriod, at instant: Date, firstWeekday: CalendarWeekday
    ) throws -> CivilDateRange {
        var error: NSError?
        let result = LifeCoreMobileCalendarBounds(
            period.rawValue, milliseconds(instant), firstWeekday.rawValue, &error
        )
        if let error { throw error }
        // The Go function returns a range whenever its error is nil.
        return CivilDateRange(from: result!.from, to: result!.to)
    }

    public static func shanghaiDate(at instant: Date) -> String {
        LifeCoreMobileShanghaiDate(milliseconds(instant))
    }

    public static func requestWindows(from: String, to: String) throws -> [CivilDateRange] {
        var error: NSError?
        let json = LifeCoreMobileRequestWindows(from, to, &error)
        if let error { throw error }
        return try JSONDecoder().decode([CivilDateRange].self, from: Data(json.utf8))
    }

    private static func milliseconds(_ date: Date) -> Int64 {
        Int64((date.timeIntervalSince1970 * 1_000).rounded(.down))
    }
}
