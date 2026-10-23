import Foundation
import XCTest
@testable import LifeClientCore

final class ClientCalendarTests: XCTestCase {
    private struct Fixtures: Decodable {
        struct Bounds: Decodable {
            let name: String
            let period: String
            let anchor: String
            let firstWeekday: Int64
            let from: String
            let to: String
        }
        struct Windows: Decodable {
            let name: String
            let from: String
            let to: String
            let ranges: [CivilDateRange]
        }
        let bounds: [Bounds]
        let windows: [Windows]
    }

    private func fixtures() throws -> Fixtures {
        let url = try XCTUnwrap(Bundle.module.url(forResource: "calendar_cases", withExtension: "json", subdirectory: "Fixtures"))
        return try JSONDecoder().decode(Fixtures.self, from: Data(contentsOf: url))
    }

    func testSharedCalendarBoundsThroughCompiledGoBinding() throws {
        let formatter = ISO8601DateFormatter()
        for fixture in try fixtures().bounds {
            let period = try XCTUnwrap(CalendarPeriod(rawValue: fixture.period))
            let weekday = try XCTUnwrap(CalendarWeekday(rawValue: fixture.firstWeekday))
            let instant = try XCTUnwrap(formatter.date(from: fixture.anchor))
            let actual = try ClientCalendar.bounds(period: period, at: instant, firstWeekday: weekday)
            XCTAssertEqual(actual.from, fixture.from, fixture.name)
            XCTAssertEqual(actual.to, fixture.to, fixture.name)
            if period == .day {
                XCTAssertEqual(ClientCalendar.shanghaiDate(at: instant), fixture.from, fixture.name)
            }
        }
    }

    func testSharedRequestWindowsThroughCompiledGoBinding() throws {
        for fixture in try fixtures().windows {
            XCTAssertEqual(
                try ClientCalendar.requestWindows(from: fixture.from, to: fixture.to),
                fixture.ranges, fixture.name
            )
        }
    }

    func testGoErrorsReachSwiftWithoutAnEmptySuccess() {
        XCTAssertThrowsError(try ClientCalendar.requestWindows(from: "2026-02-30", to: "2026-03-01"))
        XCTAssertThrowsError(try ClientCalendar.requestWindows(from: "2026-03-02", to: "2026-03-01"))
    }
}
