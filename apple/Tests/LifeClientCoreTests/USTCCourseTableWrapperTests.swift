import Foundation
import LifeClientCore
import XCTest

final class USTCCourseTableWrapperTests: XCTestCase {
    private struct Fixtures: Decodable {
        struct Request: Decodable {
            let name: String
            let semesterID: Int
            let studentID: String
            let url: String
        }
        struct StudentID: Decodable {
            let name: String
            let url: String
            let studentID: String?
            let error: Bool?
        }
        struct Table: Decodable {
            let name: String
            let response: String
            let lessonIDs: [Int]?
            let error: Bool?
        }
        let pageURL: String
        let requests: [Request]
        let studentIDs: [StudentID]
        let tables: [Table]
    }

    private func fixtures() throws -> Fixtures {
        let url = try XCTUnwrap(Bundle.module.url(
            forResource: "ustc_course_table_cases", withExtension: "json", subdirectory: "Fixtures"
        ))
        return try JSONDecoder().decode(Fixtures.self, from: Data(contentsOf: url))
    }

    func testSharedRequestsThroughCompiledGoBinding() throws {
        let fixtures = try fixtures()
        XCTAssertFalse(fixtures.requests.isEmpty)
        XCTAssertEqual(USTCCourseTable.pageURL.absoluteString, fixtures.pageURL)
        for fixture in fixtures.requests {
            XCTAssertEqual(
                USTCCourseTable.dataURL(semesterID: fixture.semesterID, studentID: fixture.studentID).absoluteString,
                fixture.url, fixture.name
            )
        }
    }

    func testSharedStudentIdentifiersThroughCompiledGoBinding() throws {
        let fixtures = try fixtures()
        XCTAssertFalse(fixtures.studentIDs.isEmpty)
        for fixture in fixtures.studentIDs {
            let url = try XCTUnwrap(URL(string: fixture.url))
            if fixture.error == true {
                XCTAssertThrowsError(try USTCCourseTable.studentID(from: url), fixture.name)
            } else {
                XCTAssertEqual(try USTCCourseTable.studentID(from: url), fixture.studentID, fixture.name)
            }
        }
    }

    func testSharedCourseTablesThroughCompiledGoBinding() throws {
        let fixtures = try fixtures()
        XCTAssertFalse(fixtures.tables.isEmpty)
        for fixture in fixtures.tables {
            let data = Data(fixture.response.utf8)
            if fixture.error == true {
                XCTAssertThrowsError(try USTCCourseTable.parse(data), fixture.name)
            } else {
                XCTAssertEqual(try USTCCourseTable.parse(data), fixture.lessonIDs, fixture.name)
            }
        }
    }
}
