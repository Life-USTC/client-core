import Foundation
import LifeCore
import XCTest

final class USTCStaticTests: XCTestCase {
    private struct Fixtures: Decodable {
        struct Path: Decodable {
            let semesterID: String
            let courseID: Int64
            let courseList: String
            let course: String
        }
        struct Semester: Decodable {
            let id: String
            let name: String
            let startDate: Double
            let endDate: Double
        }
        struct SemesterCase: Decodable {
            let name: String
            let response: String
            let semesters: [Semester]?
            let error: Bool?
        }
        struct CourseItem: Decodable {
            let id: Int64
            let name: String
            let courseCode: String
            let lessonCode: String
            let legacyLessonCode: String
            let teacherName: String
        }
        struct CourseListCase: Decodable {
            let name: String
            let response: String
            let courses: [CourseItem]?
            let error: Bool?
        }
        struct Lecture: Decodable {
            let startDate: Double
            let endDate: Double
            let name: String
            let location: String
            let teacherName: String
            let periods: Double
            let startIndex: Int64?
            let endIndex: Int64?
        }
        struct Course: Decodable {
            let id: Int64
            let name: String
            let courseCode: String
            let lessonCode: String
            let teacherName: String
            let detailText: String?
            let credit: Double
            let dateTimePlacePersonText: String?
            let lectures: [Lecture]
        }
        struct CourseCase: Decodable {
            let name: String
            let response: String
            let course: Course?
            let error: Bool?
        }
        struct MatchCase: Decodable {
            let name: String
            let appResponse: String
            let courseListResponse: String
            let errorStage: String?
            let ids: [Int64]?
            let error: Bool?
        }
        let graduateAppCoursesURL: String
        let staticSemestersPath: String
        let paths: [Path]
        let semesters: [SemesterCase]
        let courseLists: [CourseListCase]
        let courses: [CourseCase]
        let matches: [MatchCase]
    }

    private func fixtures() throws -> Fixtures {
        let url = try XCTUnwrap(Bundle.module.url(
            forResource: "ustc_static_cases", withExtension: "json", subdirectory: "Fixtures"
        ))
        return try JSONDecoder().decode(Fixtures.self, from: Data(contentsOf: url))
    }

    func testSharedStaticPathsThroughCompiledGoBinding() throws {
        let fixtures = try fixtures()
        XCTAssertEqual(LifeCoreMobileUSTCGraduateAppCoursesURL(), fixtures.graduateAppCoursesURL)
        XCTAssertEqual(LifeCoreMobileUSTCStaticSemestersPath(), fixtures.staticSemestersPath)
        XCTAssertFalse(fixtures.paths.isEmpty)
        for fixture in fixtures.paths {
            XCTAssertEqual(LifeCoreMobileUSTCStaticCourseListPath(fixture.semesterID), fixture.courseList)
            XCTAssertEqual(LifeCoreMobileUSTCStaticCoursePath(fixture.semesterID, fixture.courseID), fixture.course)
        }
    }

    func testSharedStaticSemestersThroughCompiledGoBinding() throws {
        let cases = try fixtures().semesters
        XCTAssertFalse(cases.isEmpty)
        for fixture in cases {
            var error: NSError?
            let result = LifeCoreMobileParseUSTCStaticSemesters(Data(fixture.response.utf8), &error)
            if fixture.error == true {
                XCTAssertNotNil(error, fixture.name)
                XCTAssertNil(result, fixture.name)
                continue
            }
            XCTAssertNil(error, fixture.name)
            let actual = try XCTUnwrap(result)
            let expected = try XCTUnwrap(fixture.semesters)
            XCTAssertEqual(actual.count(), Int64(expected.count), fixture.name)
            for (index, semester) in expected.enumerated() {
                let value = try XCTUnwrap(actual.at(Int64(index)))
                XCTAssertEqual(value.identifier, semester.id, fixture.name)
                XCTAssertEqual(value.name, semester.name, fixture.name)
                XCTAssertEqual(value.startDate, semester.startDate, fixture.name)
                XCTAssertEqual(value.endDate, semester.endDate, fixture.name)
            }
        }
    }

    func testSharedStaticCourseListThroughCompiledGoBinding() throws {
        let cases = try fixtures().courseLists
        XCTAssertFalse(cases.isEmpty)
        for fixture in cases {
            var error: NSError?
            let result = LifeCoreMobileParseUSTCStaticCourseList(Data(fixture.response.utf8), &error)
            if fixture.error == true {
                XCTAssertNotNil(error, fixture.name)
                XCTAssertNil(result, fixture.name)
                continue
            }
            XCTAssertNil(error, fixture.name)
            let actual = try XCTUnwrap(result)
            let expected = try XCTUnwrap(fixture.courses)
            XCTAssertEqual(actual.count(), Int64(expected.count), fixture.name)
            for (index, course) in expected.enumerated() {
                let value = try XCTUnwrap(actual.at(Int64(index)))
                XCTAssertEqual(value.identifier, course.id, fixture.name)
                XCTAssertEqual(value.name, course.name, fixture.name)
                XCTAssertEqual(value.courseCode, course.courseCode, fixture.name)
                XCTAssertEqual(value.lessonCode, course.lessonCode, fixture.name)
                XCTAssertEqual(value.legacyLessonCode, course.legacyLessonCode, fixture.name)
                XCTAssertEqual(value.teacherName, course.teacherName, fixture.name)
            }
        }
    }

    func testSharedStaticCourseThroughCompiledGoBinding() throws {
        let cases = try fixtures().courses
        XCTAssertFalse(cases.isEmpty)
        for fixture in cases {
            var error: NSError?
            let result = LifeCoreMobileParseUSTCStaticCourse(Data(fixture.response.utf8), &error)
            if fixture.error == true {
                XCTAssertNotNil(error, fixture.name)
                XCTAssertNil(result, fixture.name)
                continue
            }
            XCTAssertNil(error, fixture.name)
            let actual = try XCTUnwrap(result)
            let expected = try XCTUnwrap(fixture.course)
            XCTAssertEqual(actual.identifier, expected.id, fixture.name)
            XCTAssertEqual(actual.name, expected.name, fixture.name)
            XCTAssertEqual(actual.courseCode, expected.courseCode, fixture.name)
            XCTAssertEqual(actual.lessonCode, expected.lessonCode, fixture.name)
            XCTAssertEqual(actual.teacherName, expected.teacherName, fixture.name)
            XCTAssertEqual(actual.credit, expected.credit, fixture.name)
            XCTAssertEqual(actual.hasDetailText ? actual.detailText : nil, expected.detailText, fixture.name)
            XCTAssertEqual(actual.hasDateTimePlacePersonText ? actual.dateTimePlacePersonText : nil, expected.dateTimePlacePersonText, fixture.name)
            XCTAssertEqual(actual.lectureCount(), Int64(expected.lectures.count), fixture.name)
            for (index, lecture) in expected.lectures.enumerated() {
                let value = try XCTUnwrap(actual.lecture(at: Int64(index)))
                XCTAssertEqual(value.startDate, lecture.startDate, fixture.name)
                XCTAssertEqual(value.endDate, lecture.endDate, fixture.name)
                XCTAssertEqual(value.name, lecture.name, fixture.name)
                XCTAssertEqual(value.location, lecture.location, fixture.name)
                XCTAssertEqual(value.teacherName, lecture.teacherName, fixture.name)
                XCTAssertEqual(value.periods, lecture.periods, fixture.name)
                XCTAssertEqual(value.hasStartIndex ? value.startIndex : nil, lecture.startIndex, fixture.name)
                XCTAssertEqual(value.hasEndIndex ? value.endIndex : nil, lecture.endIndex, fixture.name)
            }
        }
    }

    func testSharedGraduateAppMatchingThroughCompiledGoBinding() throws {
        let cases = try fixtures().matches
        XCTAssertFalse(cases.isEmpty)
        for fixture in cases {
            var error: NSError?
            let selection = LifeCoreMobileParseUSTCGraduateAppCourses(Data(fixture.appResponse.utf8), &error)
            if fixture.errorStage == "app" {
                XCTAssertEqual(fixture.error, true, fixture.name)
                XCTAssertNotNil(error, fixture.name)
                XCTAssertNil(selection, fixture.name)
                continue
            }
            XCTAssertNil(error, fixture.name)
            let parsed = try XCTUnwrap(selection)
            do {
                let actual = try XCTUnwrap(parsed.matchCourseIDs(Data(fixture.courseListResponse.utf8)))
                XCTAssertNotEqual(fixture.error, true, fixture.name)
                let expected = try XCTUnwrap(fixture.ids)
                XCTAssertEqual((0..<actual.count()).map { actual.at($0) }, expected, fixture.name)
            } catch {
                XCTAssertEqual(fixture.error, true, fixture.name)
                XCTAssertEqual(fixture.errorStage, "courseList", fixture.name)
            }
        }
    }
}
