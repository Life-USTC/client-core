import Foundation
import LifeCore
import XCTest

final class USTCAcademicTests: XCTestCase {
    private struct Fixtures: Decodable {
        struct ScoreReport: Decodable {
            struct Overview: Decodable { let gpa: Double? }
            struct Rank: Decodable {
                let majorName: String?
                let majorRank: Int64?
                let majorStdCount: Int64?
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

    func testSharedAcademicRequestsThroughCompiledGoBinding() throws {
        let cases = try fixtures()
        XCTAssertEqual(LifeCoreMobileUSTCScalarScoreListURL(1), cases.scalarScoreURL)
        XCTAssertEqual(LifeCoreMobileUSTCExamArrangementURL(), cases.examURL)
    }

    func testSharedScalarScoresThroughCompiledGoBinding() throws {
        let cases = try fixtures().scores.filter { $0.format == "scalar" }
        XCTAssertFalse(cases.isEmpty)
        for fixture in cases {
            var error: NSError?
            let result = LifeCoreMobileParseUSTCScalarScores(Data(fixture.response.utf8), &error)
            if fixture.error == true {
                XCTAssertNotNil(error, fixture.name)
                XCTAssertNil(result, fixture.name)
                continue
            }
            XCTAssertNil(error, fixture.name)
            let report = try XCTUnwrap(result)
            let expected = try XCTUnwrap(fixture.expected)
            XCTAssertEqual(report.hasGPA, expected.overview?.gpa != nil, fixture.name)
            if let value = expected.overview?.gpa { XCTAssertEqual(report.gpa, value, fixture.name) }
            XCTAssertEqual(report.hasMajorName, expected.rank?.majorName != nil, fixture.name)
            XCTAssertEqual(report.hasMajorRank, expected.rank?.majorRank != nil, fixture.name)
            XCTAssertEqual(report.hasMajorStdCount, expected.rank?.majorStdCount != nil, fixture.name)
            if let value = expected.rank?.majorName { XCTAssertEqual(report.majorName, value, fixture.name) }
            if let value = expected.rank?.majorRank { XCTAssertEqual(report.majorRank, value, fixture.name) }
            if let value = expected.rank?.majorStdCount { XCTAssertEqual(report.majorStdCount, value, fixture.name) }
            XCTAssertEqual(report.count(), Int64(expected.items.count), fixture.name)
            for (index, entry) in expected.items.enumerated() {
                let actual = try XCTUnwrap(report.at(Int64(index)))
                XCTAssertEqual(actual.semesterID, entry.semesterId ?? "", fixture.name)
                XCTAssertEqual(actual.semesterName, entry.semesterName ?? "", fixture.name)
                XCTAssertEqual(actual.courseName, entry.courseName ?? "", fixture.name)
                XCTAssertEqual(actual.courseCode, entry.courseCode ?? "", fixture.name)
                XCTAssertEqual(actual.lessonCode, entry.lessonCode ?? "", fixture.name)
                XCTAssertEqual(actual.credits, entry.credits ?? 0, fixture.name)
                XCTAssertEqual(actual.score, entry.score ?? "", fixture.name)
                XCTAssertEqual(actual.hasGradePoint, entry.gradePoint != nil, fixture.name)
                if let value = entry.gradePoint { XCTAssertEqual(actual.gradePoint, value, fixture.name) }
            }
        }
    }

    func testSharedRoomExamsThroughCompiledGoBinding() throws {
        let cases = try fixtures().exams.filter { $0.format == "room" }
        XCTAssertFalse(cases.isEmpty)
        for fixture in cases {
            let data: Data
            if let base64 = fixture.responseBase64 {
                data = try XCTUnwrap(Data(base64Encoded: base64))
            } else {
                data = Data(try XCTUnwrap(fixture.response).utf8)
            }
            var error: NSError?
            let result = LifeCoreMobileParseUSTCExamRoomTable(data, &error)
            if fixture.error == true {
                XCTAssertNotNil(error, fixture.name)
                XCTAssertNil(result, fixture.name)
                continue
            }
            XCTAssertNil(error, fixture.name)
            let exams = try XCTUnwrap(result)
            let expected = try XCTUnwrap(fixture.expected)
            XCTAssertEqual(exams.count(), Int64(expected.count), fixture.name)
            for (index, exam) in expected.enumerated() {
                let actual = try XCTUnwrap(exams.at(Int64(index)))
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
