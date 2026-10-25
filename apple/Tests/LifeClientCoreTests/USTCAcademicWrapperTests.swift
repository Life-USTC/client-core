import Foundation
import LifeClientCore
import XCTest

final class USTCAcademicWrapperTests: XCTestCase {
    private struct Fixtures: Decodable {
        struct ScoreReport: Decodable {
            struct Overview: Decodable { let gpa: Double? }
            struct Rank: Decodable {
                let majorName: String?
                let majorRank: Int?
                let majorStdCount: Int?
            }
            struct Entry: Decodable {
                let semesterId: String?
                let semesterName: String?
                let courseName: String?
                let courseCode: String?
                let lessonCode: String?
                let credits: Double?
                let gradePoint: Double?
                let score: String?
            }
            let overview: Overview?
            let rank: Rank?
            let items: [Entry]
        }
        struct ScoreCase: Decodable {
            let name: String
            let format: String
            let response: String
            let expected: ScoreReport?
            let error: Bool?
        }
        struct Exam: Decodable {
            let courseName: String?
            let lessonCode: String?
            let examType: String?
            let dateTime: String?
            let roomName: String?
            let buildingName: String?
            let district: String?
            let detail: String?
        }
        struct ExamCase: Decodable {
            let name: String
            let format: String
            let response: String?
            let responseBase64: String?
            let expected: [Exam]?
            let error: Bool?
        }
        let scalarScoreURL: String
        let examURL: String
        let scores: [ScoreCase]
        let exams: [ExamCase]
    }

    private func fixtures() throws -> Fixtures {
        let url = try XCTUnwrap(Bundle.module.url(
            forResource: "ustc_academic_cases", withExtension: "json", subdirectory: "Fixtures"
        ))
        return try JSONDecoder().decode(Fixtures.self, from: Data(contentsOf: url))
    }

    func testSharedAcademicRequestURLs() throws {
        let cases = try fixtures()
        XCTAssertEqual(USTCAcademic.scoreURL().absoluteString, cases.scalarScoreURL)
        XCTAssertEqual(USTCAcademic.examURL.absoluteString, cases.examURL)
    }

    func testSharedScoresKeepNativeOptionalValues() throws {
        let cases = try fixtures().scores.filter { $0.format == "scalar" }
        XCTAssertFalse(cases.isEmpty)
        for fixture in cases {
            let data = Data(fixture.response.utf8)
            if fixture.error == true {
                XCTAssertThrowsError(try USTCAcademic.parseScores(data), fixture.name)
                continue
            }
            let report = try USTCAcademic.parseScores(data)
            let expected = try XCTUnwrap(fixture.expected)
            XCTAssertEqual(report.gpa, expected.overview?.gpa, fixture.name)
            XCTAssertEqual(report.majorName, expected.rank?.majorName, fixture.name)
            XCTAssertEqual(report.majorRank, expected.rank?.majorRank, fixture.name)
            XCTAssertEqual(report.majorStdCount, expected.rank?.majorStdCount, fixture.name)
            XCTAssertEqual(report.entries.count, expected.items.count, fixture.name)
            for (actual, entry) in zip(report.entries, expected.items) {
                XCTAssertEqual(actual.semesterID, entry.semesterId ?? "", fixture.name)
                XCTAssertEqual(actual.semesterName, entry.semesterName ?? "", fixture.name)
                XCTAssertEqual(actual.courseName, entry.courseName ?? "", fixture.name)
                XCTAssertEqual(actual.courseCode, entry.courseCode ?? "", fixture.name)
                XCTAssertEqual(actual.lessonCode, entry.lessonCode ?? "", fixture.name)
                XCTAssertEqual(actual.credits, entry.credits ?? 0, fixture.name)
                XCTAssertEqual(actual.gradePoint, entry.gradePoint, fixture.name)
                XCTAssertEqual(actual.score, entry.score ?? "", fixture.name)
            }
        }
    }

    func testSharedExamsKeepSeparateRoomFieldsAndDateText() throws {
        let cases = try fixtures().exams.filter { $0.format == "room" }
        XCTAssertFalse(cases.isEmpty)
        for fixture in cases {
            let data: Data
            if let base64 = fixture.responseBase64 {
                data = try XCTUnwrap(Data(base64Encoded: base64))
            } else {
                data = Data(try XCTUnwrap(fixture.response).utf8)
            }
            if fixture.error == true {
                XCTAssertThrowsError(try USTCAcademic.parseExams(data), fixture.name)
                continue
            }
            let exams = try USTCAcademic.parseExams(data)
            let expected = try XCTUnwrap(fixture.expected)
            XCTAssertEqual(exams.count, expected.count, fixture.name)
            for (actual, exam) in zip(exams, expected) {
                XCTAssertEqual(actual.courseName, exam.courseName ?? "", fixture.name)
                XCTAssertEqual(actual.lessonCode, exam.lessonCode ?? "", fixture.name)
                XCTAssertEqual(actual.examType, exam.examType ?? "", fixture.name)
                XCTAssertEqual(actual.dateTime, exam.dateTime ?? "", fixture.name)
                XCTAssertEqual(actual.roomName, exam.roomName ?? "", fixture.name)
                XCTAssertEqual(actual.buildingName, exam.buildingName ?? "", fixture.name)
                XCTAssertEqual(actual.district, exam.district ?? "", fixture.name)
                XCTAssertEqual(actual.detail, exam.detail ?? "", fixture.name)
            }
        }
    }
}
