import Foundation
import XCTest
@testable import LifeClientCore

final class TeachingWeekTests: XCTestCase {
    private struct Fixtures: Decodable {
        struct Bounds: Decodable {
            let name: String
            let semesterStart: String
            let weekIndex: Int
            let firstWeekday: Int64
            let from: String
            let to: String
        }
        struct Number: Decodable {
            let name: String
            let semesterStart: String
            let selectedDate: String
            let firstWeekday: Int64
            let weekIndex: Int
        }
        let bounds: [Bounds]
        let numbers: [Number]
    }

    private func fixtures() throws -> Fixtures {
        let url = try XCTUnwrap(Bundle.module.url(
            forResource: "teaching_week_cases", withExtension: "json", subdirectory: "Fixtures"
        ))
        return try JSONDecoder().decode(Fixtures.self, from: Data(contentsOf: url))
    }

    func testSharedTeachingWeekBoundsThroughCompiledGoBinding() throws {
        for fixture in try fixtures().bounds {
            let weekday = try XCTUnwrap(CalendarWeekday(rawValue: fixture.firstWeekday))
            let actual = try ClientCalendar.teachingWeekBounds(
                semesterStart: fixture.semesterStart, weekIndex: fixture.weekIndex, firstWeekday: weekday
            )
            XCTAssertEqual(actual.from, fixture.from, fixture.name)
            XCTAssertEqual(actual.to, fixture.to, fixture.name)
        }
    }

    func testSharedTeachingWeekNumbersThroughCompiledGoBinding() throws {
        for fixture in try fixtures().numbers {
            let weekday = try XCTUnwrap(CalendarWeekday(rawValue: fixture.firstWeekday))
            XCTAssertEqual(
                try ClientCalendar.teachingWeekNumber(
                    semesterStart: fixture.semesterStart, selectedDate: fixture.selectedDate, firstWeekday: weekday
                ),
                fixture.weekIndex, fixture.name
            )
        }
    }

    func testTeachingWeekErrorsReachSwift() {
        XCTAssertThrowsError(try ClientCalendar.teachingWeekBounds(
            semesterStart: "2026-09-01", weekIndex: 0, firstWeekday: .monday
        ))
        XCTAssertThrowsError(try ClientCalendar.teachingWeekBounds(
            semesterStart: "2026-02-30", weekIndex: 1, firstWeekday: .monday
        ))
        XCTAssertThrowsError(try ClientCalendar.teachingWeekNumber(
            semesterStart: "2026-09-01", selectedDate: "bad", firstWeekday: .sunday
        ))
    }
}
