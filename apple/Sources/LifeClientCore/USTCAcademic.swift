import Foundation
import LifeCore

public struct USTCScoreEntry: Equatable, Sendable {
    public let semesterID: String
    public let semesterName: String
    public let courseName: String
    public let courseCode: String
    public let lessonCode: String
    public let credits: Double
    public let gradePoint: Double?
    public let score: String
}

public struct USTCScoreReport: Equatable, Sendable {
    public let gpa: Double?
    public let majorName: String?
    public let majorRank: Int?
    public let majorStdCount: Int?
    public let entries: [USTCScoreEntry]
}

public struct USTCExamEntry: Equatable, Sendable {
    public let courseName: String
    public let lessonCode: String
    public let examType: String
    public let dateTime: String
    public let roomName: String
    public let buildingName: String
    public let district: String
    public let detail: String
}

/// Native value conversion for the existing iOS score and exam response formats.
public enum USTCAcademic {
    public static func scoreURL(trainTypeID: Int = 1) -> URL {
        URL(string: LifeCoreMobileUSTCScalarScoreListURL(Int64(trainTypeID)))!
    }

    public static var examURL: URL {
        URL(string: LifeCoreMobileUSTCExamArrangementURL())!
    }

    public static func parseScores(_ data: Data) throws -> USTCScoreReport {
        var error: NSError?
        let result = LifeCoreMobileParseUSTCScalarScores(data, &error)
        if let error { throw error }
        // The Go parser returns a report whenever its error is nil.
        let report = result!
        let entries = (0..<report.count()).map { index in
            let entry = report.at(index)!
            return USTCScoreEntry(
                semesterID: entry.semesterID, semesterName: entry.semesterName,
                courseName: entry.courseName, courseCode: entry.courseCode,
                lessonCode: entry.lessonCode, credits: entry.credits,
                gradePoint: entry.hasGradePoint ? entry.gradePoint : nil, score: entry.score
            )
        }
        return USTCScoreReport(
            gpa: report.hasGPA ? report.gpa : nil,
            majorName: report.hasMajorName ? report.majorName : nil,
            majorRank: report.hasMajorRank ? Int(report.majorRank) : nil,
            majorStdCount: report.hasMajorStdCount ? Int(report.majorStdCount) : nil,
            entries: entries
        )
    }

    public static func parseExams(_ data: Data) throws -> [USTCExamEntry] {
        var error: NSError?
        let result = LifeCoreMobileParseUSTCExamRoomTable(data, &error)
        if let error { throw error }
        // Dates remain text for the caller's existing timezone/import policy.
        let exams = result!
        return (0..<exams.count()).map { index in
            let exam = exams.at(index)!
            return USTCExamEntry(
                courseName: exam.courseName, lessonCode: exam.lessonCode,
                examType: exam.examType, dateTime: exam.dateTime,
                roomName: exam.roomName, buildingName: exam.buildingName,
                district: exam.district, detail: exam.detail
            )
        }
    }
}
