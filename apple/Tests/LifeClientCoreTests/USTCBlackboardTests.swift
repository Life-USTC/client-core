import Foundation
import LifeCore
import XCTest

final class USTCBlackboardTests: XCTestCase {
    private struct Fixtures: Decodable {
        struct Event: Decodable, Equatable {
            let title: String
            let courseName: String
            let endDate: String
        }
        struct SelectedEvents: Decodable {
            let name: String
            let response: String
            let events: [Event]?
            let error: Bool?
        }
        let selectedEventsURL: String
        let selectedEvents: [SelectedEvents]
    }

    private func fixtures() throws -> Fixtures {
        let url = try XCTUnwrap(Bundle.module.url(
            forResource: "ustc_blackboard_cases", withExtension: "json", subdirectory: "Fixtures"
        ))
        return try JSONDecoder().decode(Fixtures.self, from: Data(contentsOf: url))
    }

    func testSharedSelectedEventsURLThroughCompiledGoBinding() throws {
        let fixtures = try fixtures()
        XCTAssertFalse(fixtures.selectedEventsURL.isEmpty)
        XCTAssertEqual(LifeCoreMobileUSTCBlackboardSelectedCalendarEventsURL(), fixtures.selectedEventsURL)
    }

    func testSharedSelectedEventsThroughCompiledGoBinding() throws {
        let fixtures = try fixtures()
        XCTAssertFalse(fixtures.selectedEvents.isEmpty)
        for fixture in fixtures.selectedEvents {
            let data = Data(fixture.response.utf8)
            var error: NSError?
            let value = LifeCoreMobileParseUSTCBlackboardSelectedCalendarEvents(data, &error)
            if fixture.error == true {
                XCTAssertNotNil(error, fixture.name)
                XCTAssertNil(value, fixture.name)
            } else {
                XCTAssertNil(error, fixture.name)
                let events = try XCTUnwrap(value, fixture.name)
                let expected = try XCTUnwrap(fixture.events, fixture.name)
                let actual = (0..<events.count()).map { index in
                    Fixtures.Event(
                        title: events.title(index),
                        courseName: events.courseName(index),
                        endDate: events.endDate(index)
                    )
                }
                XCTAssertEqual(actual, expected, fixture.name)
            }
        }
    }
}
