import XCTest
import LifeClientCore

final class USTCBlackboardWrapperTests: XCTestCase {
    private struct Fixtures: Decodable {
        let selectedEventsURL: String
        let selectedEvents: [EventCase]
    }

    private struct EventCase: Decodable {
        let name: String
        let response: String
        let events: [Event]?
        let error: Bool?
    }

    private struct Event: Decodable {
        let title: String
        let courseName: String
        let endDate: String
    }

    private func fixtures() throws -> Fixtures {
        let url = try XCTUnwrap(Bundle.module.url(
            forResource: "ustc_blackboard_cases", withExtension: "json", subdirectory: "Fixtures"
        ))
        return try JSONDecoder().decode(Fixtures.self, from: Data(contentsOf: url))
    }

    func testSelectedCalendarEventsURL() throws {
        XCTAssertEqual(USTCBlackboard.selectedCalendarEventsURL.absoluteString, try fixtures().selectedEventsURL)
    }

    func testSelectedCalendarEventsAgainstSharedFixtures() throws {
        for fixture in try fixtures().selectedEvents {
            let data = Data(fixture.response.utf8)
            if fixture.error == true {
                XCTAssertThrowsError(try USTCBlackboard.parseSelectedCalendarEvents(data), fixture.name)
                continue
            }
            let expected = try XCTUnwrap(fixture.events, fixture.name).map {
                USTCBlackboardSelectedCalendarEvent(title: $0.title, courseName: $0.courseName, endDate: $0.endDate)
            }
            XCTAssertEqual(try USTCBlackboard.parseSelectedCalendarEvents(data), expected, fixture.name)
        }
    }
}
