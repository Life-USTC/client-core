import Foundation
import LifeCore
import XCTest

final class USTCCourseTableTests: XCTestCase {
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
        XCTAssertEqual(LifeCoreMobileUSTCCourseTablePageURL(), fixtures.pageURL)
        for fixture in fixtures.requests {
            XCTAssertEqual(
                LifeCoreMobileUSTCCourseTableDataURL(Int64(fixture.semesterID), fixture.studentID),
                fixture.url, fixture.name
            )
        }
    }

    func testSharedStudentIdentifiersThroughCompiledGoBinding() throws {
        let fixtures = try fixtures()
        XCTAssertFalse(fixtures.studentIDs.isEmpty)
        for fixture in fixtures.studentIDs {
            var error: NSError?
            let value = LifeCoreMobileUSTCStudentIDFromCourseTableURL(fixture.url, &error)
            if fixture.error == true {
                XCTAssertNotNil(error, fixture.name)
                XCTAssertEqual(value, "", fixture.name)
            } else {
                XCTAssertNil(error, fixture.name)
                XCTAssertEqual(value, fixture.studentID, fixture.name)
            }
        }
    }

    func testSharedCourseTablesThroughCompiledGoBinding() throws {
        let fixtures = try fixtures()
        XCTAssertFalse(fixtures.tables.isEmpty)
        for fixture in fixtures.tables {
            let data = Data(fixture.response.utf8)
            var error: NSError?
            let value = LifeCoreMobileParseUSTCCourseTable(data, &error)
            if fixture.error == true {
                XCTAssertNotNil(error, fixture.name)
                XCTAssertNil(value, fixture.name)
            } else {
                XCTAssertNil(error, fixture.name)
                let identifiers = try XCTUnwrap(value)
                let actual = (0..<identifiers.count()).map { Int(identifiers.at($0)) }
                XCTAssertEqual(actual, fixture.lessonIDs, fixture.name)
            }
        }
    }
}
