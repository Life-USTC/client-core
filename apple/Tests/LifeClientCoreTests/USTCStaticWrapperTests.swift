import Foundation
import LifeClientCore
import XCTest

final class USTCStaticWrapperTests: XCTestCase {
    private struct ResponseCase: Decodable {
        let name: String
        let response: String
    }

    private struct MatchCase: Decodable {
        let name: String
        let appResponse: String
        let courseListResponse: String
    }

    private struct Fixtures: Decodable {
        let graduateAppCoursesURL: String
        let staticSemestersPath: String
        let semesters: [ResponseCase]
        let courseLists: [ResponseCase]
        let courses: [ResponseCase]
        let matches: [MatchCase]
    }

    private func fixtures() throws -> Fixtures {
        let url = try XCTUnwrap(Bundle.module.url(
            forResource: "ustc_static_cases", withExtension: "json", subdirectory: "Fixtures"
        ))
        return try JSONDecoder().decode(Fixtures.self, from: Data(contentsOf: url))
    }

    private func response(_ cases: [ResponseCase], named name: String) throws -> Data {
        Data(try XCTUnwrap(cases.first { $0.name == name }, name).response.utf8)
    }

    func testNativeRequestValues() throws {
        let fixture = try fixtures()
        XCTAssertEqual(USTCStatic.graduateAppCoursesURL.absoluteString, fixture.graduateAppCoursesURL)
        XCTAssertEqual(USTCStatic.semestersPath, fixture.staticSemestersPath)
        XCTAssertEqual(USTCStatic.courseListPath(semesterID: "2026FA"), "curriculum/2026FA/courses.json")
        XCTAssertEqual(USTCStatic.coursePath(semesterID: "2026FA", courseID: 43), "curriculum/2026FA/43.json")
    }

    func testSemesterAndCourseListDTOFields() throws {
        let fixture = try fixtures()
        let semesterData = try response(fixture.semesters, named: "raw-seconds-order-and-duplicates")
        let semesters = try USTCStatic.parseSemesters(semesterData)
        XCTAssertEqual(semesters.map(\.id), ["2026FA", "2025SP", "2026FA"])
        let firstSemester = try XCTUnwrap(semesters.first)
        let secondSemester = try XCTUnwrap(semesters.dropFirst().first)
        XCTAssertEqual(firstSemester.name, " 秋季 ")
        XCTAssertEqual(firstSemester.startDate, 1788105600.25)
        XCTAssertEqual(firstSemester.endDate, 1801238399.5)
        XCTAssertEqual(secondSemester.startDate, 0)
        XCTAssertEqual(secondSemester.endDate, -1)

        let listData = try response(fixture.courseLists, named: "canonical-and-legacy-stay-distinct")
        let courses = try USTCStatic.parseCourseList(listData)
        XCTAssertEqual(courses.map(\.id), [43, 43])
        let firstCourse = try XCTUnwrap(courses.first)
        let secondCourse = try XCTUnwrap(courses.dropFirst().first)
        XCTAssertEqual(firstCourse.name, " Algebra ")
        XCTAssertEqual(firstCourse.courseCode, "MATH")
        XCTAssertEqual(firstCourse.lessonCode, "canonical")
        XCTAssertEqual(firstCourse.legacyLessonCode, "legacy")
        XCTAssertEqual(firstCourse.teacherName, " Teacher ")
        XCTAssertEqual(secondCourse.lessonCode, "")
        XCTAssertEqual(secondCourse.legacyLessonCode, "only-legacy")
    }

    func testCourseDetailPreservesDatesAndOptionalValues() throws {
        let fixture = try fixtures()
        let data = try response(fixture.courses, named: "dated-lectures-and-optional-zero")
        let course = try USTCStatic.parseCourse(data)
        XCTAssertEqual(course.id, 43)
        XCTAssertEqual(course.name, " Algebra ")
        XCTAssertEqual(course.courseCode, "MATH")
        XCTAssertEqual(course.lessonCode, "MATH.01")
        XCTAssertEqual(course.teacherName, "Teacher")
        XCTAssertEqual(course.detailText, "")
        XCTAssertEqual(course.credit, 2.5)
        XCTAssertEqual(course.dateTimePlacePersonText, " raw schedule ")
        XCTAssertEqual(course.lectures.count, 2)
        let lecture = try XCTUnwrap(course.lectures.first)
        XCTAssertEqual(lecture.startDate, 1788105600.25)
        XCTAssertEqual(lecture.endDate, 1788108300.5)
        XCTAssertEqual(lecture.name, "Lecture")
        XCTAssertEqual(lecture.location, "Room")
        XCTAssertEqual(lecture.teacherName, "Teacher")
        XCTAssertEqual(lecture.periods, 1.5)
        XCTAssertEqual(lecture.startIndex, 0)
        XCTAssertNil(lecture.endIndex)
        let secondLecture = try XCTUnwrap(course.lectures.dropFirst().first)
        XCTAssertNil(secondLecture.startIndex)
        XCTAssertEqual(secondLecture.endIndex, 3)

        let missingData = try response(fixture.courses, named: "producer-description-is-not-detail-text")
        let missing = try USTCStatic.parseCourse(missingData)
        XCTAssertNil(missing.detailText)
        XCTAssertNil(missing.dateTimePlacePersonText)
        XCTAssertTrue(missing.lectures.isEmpty)
    }

    func testDefaultingNumbersAndOptionalIndexesRemainDistinct() throws {
        let data = try response(try fixtures().courses, named: "number-value-and-optional-int-differ")
        let course = try USTCStatic.parseCourse(data)
        XCTAssertEqual(course.id, 42)
        XCTAssertEqual(course.credit, 2.5)
        XCTAssertEqual(course.name, "true")
        XCTAssertEqual(course.courseCode, "12")
        XCTAssertEqual(course.lessonCode, "false")
        XCTAssertEqual(course.teacherName, "1.25")
        XCTAssertNil(course.detailText)
        XCTAssertEqual(course.dateTimePlacePersonText, "")
        XCTAssertEqual(course.lectures.count, 4)
        let firstLecture = try XCTUnwrap(course.lectures.first)
        let secondLecture = try XCTUnwrap(course.lectures.dropFirst().first)
        XCTAssertEqual(firstLecture.startDate, 100)
        XCTAssertEqual(firstLecture.endDate, 0)
        XCTAssertNil(firstLecture.startIndex)
        XCTAssertEqual(firstLecture.endIndex, 0)
        XCTAssertEqual(secondLecture.periods, 1.25)
        XCTAssertEqual(secondLecture.startIndex, 2)
        XCTAssertEqual(secondLecture.endIndex, 0)
    }

    func testGraduateSelectionValidatesAppBeforeListAndPreservesIDs() throws {
        let fixture = try XCTUnwrap(try fixtures().matches.first {
            $0.name == "canonical-or-legacy-and-repeated-ids"
        })
        let selection = try USTCStatic.parseGraduateAppCourses(Data(fixture.appResponse.utf8))
        XCTAssertEqual(try selection.matchCourseIDs(in: Data(fixture.courseListResponse.utf8)), [101, 102, 103, 101])
        XCTAssertThrowsError(try USTCStatic.parseGraduateAppCourses(Data("<html>login</html>".utf8)))
        XCTAssertThrowsError(try selection.matchCourseIDs(in: Data("[".utf8)))
    }

    func testJSONShapeDefaultsAndErrorsRemainVisible() throws {
        XCTAssertTrue(try USTCStatic.parseSemesters(Data("{}".utf8)).isEmpty)
        XCTAssertTrue(try USTCStatic.parseCourseList(Data("{}".utf8)).isEmpty)
        let course = try USTCStatic.parseCourse(Data("[]".utf8))
        XCTAssertEqual(course.id, 0)
        XCTAssertEqual(course.name, "")
        XCTAssertEqual(course.credit, 0)
        XCTAssertNil(course.detailText)
        XCTAssertNil(course.dateTimePlacePersonText)
        XCTAssertTrue(course.lectures.isEmpty)
        let emptySelection = try USTCStatic.parseGraduateAppCourses(Data("{}".utf8))
        XCTAssertTrue(try emptySelection.matchCourseIDs(in: Data("[]".utf8)).isEmpty)
        XCTAssertThrowsError(try USTCStatic.parseSemesters(Data("[".utf8)))
        XCTAssertThrowsError(try USTCStatic.parseCourseList(Data("null".utf8)))
        XCTAssertThrowsError(try USTCStatic.parseCourse(Data("null".utf8)))
    }
}
