import Foundation
import LifeCore

public struct USTCStaticSemester: Equatable, Sendable {
    public let id: String
    public let name: String
    public let startDate: TimeInterval
    public let endDate: TimeInterval
}

public struct USTCStaticCourseListItem: Equatable, Sendable {
    public let id: Int
    public let name: String
    public let courseCode: String
    public let lessonCode: String
    public let legacyLessonCode: String
    public let teacherName: String
}

public struct USTCStaticLecture: Equatable, Sendable {
    public let startDate: TimeInterval
    public let endDate: TimeInterval
    public let name: String
    public let location: String
    public let teacherName: String
    public let periods: Double
    public let startIndex: Int?
    public let endIndex: Int?
}

public struct USTCStaticCourse: Equatable, Sendable {
    public let id: Int
    public let name: String
    public let courseCode: String
    public let lessonCode: String
    public let teacherName: String
    public let detailText: String?
    public let credit: Double
    public let dateTimePlacePersonText: String?
    public let lectures: [USTCStaticLecture]
}

/// A parsed, immutable app response. Create it before fetching the static list
/// to retain the native helper's existing validation/request ordering.
public struct USTCGraduateAppCourses {
    fileprivate let parsed: LifeCoreMobileUSTCGraduateAppCourses

    public func matchCourseIDs(in courseListData: Data) throws -> [Int] {
        let ids = try parsed.matchCourseIDs(courseListData)
        return (0..<ids.count()).map { Int(ids.at($0)) }
    }
}

public enum USTCStatic {
    public static var graduateAppCoursesURL: URL {
        URL(string: LifeCoreMobileUSTCGraduateAppCoursesURL())!
    }

    public static var semestersPath: String { LifeCoreMobileUSTCStaticSemestersPath() }

    public static func courseListPath(semesterID: String) -> String {
        LifeCoreMobileUSTCStaticCourseListPath(semesterID)
    }

    public static func coursePath(semesterID: String, courseID: Int) -> String {
        LifeCoreMobileUSTCStaticCoursePath(semesterID, Int64(courseID))
    }

    public static func parseGraduateAppCourses(_ data: Data) throws -> USTCGraduateAppCourses {
        var error: NSError?
        let result = LifeCoreMobileParseUSTCGraduateAppCourses(data, &error)
        if let error { throw error }
        return USTCGraduateAppCourses(parsed: result!)
    }

    public static func parseSemesters(_ data: Data) throws -> [USTCStaticSemester] {
        var error: NSError?
        let result = LifeCoreMobileParseUSTCStaticSemesters(data, &error)
        if let error { throw error }
        let semesters = result!
        return (0..<semesters.count()).map { index in
            let value = semesters.at(index)!
            return USTCStaticSemester(
                id: value.identifier, name: value.name,
                startDate: value.startDate, endDate: value.endDate
            )
        }
    }

    public static func parseCourseList(_ data: Data) throws -> [USTCStaticCourseListItem] {
        var error: NSError?
        let result = LifeCoreMobileParseUSTCStaticCourseList(data, &error)
        if let error { throw error }
        let courses = result!
        return (0..<courses.count()).map { index in
            let value = courses.at(index)!
            return USTCStaticCourseListItem(
                id: Int(value.identifier), name: value.name, courseCode: value.courseCode,
                lessonCode: value.lessonCode, legacyLessonCode: value.legacyLessonCode,
                teacherName: value.teacherName
            )
        }
    }

    public static func parseCourse(_ data: Data) throws -> USTCStaticCourse {
        var error: NSError?
        let result = LifeCoreMobileParseUSTCStaticCourse(data, &error)
        if let error { throw error }
        let course = result!
        let lectures = (0..<course.lectureCount()).map { index in
            let value = course.lecture(at: index)!
            return USTCStaticLecture(
                startDate: value.startDate, endDate: value.endDate, name: value.name,
                location: value.location, teacherName: value.teacherName, periods: value.periods,
                startIndex: value.hasStartIndex ? Int(value.startIndex) : nil,
                endIndex: value.hasEndIndex ? Int(value.endIndex) : nil
            )
        }
        return USTCStaticCourse(
            id: Int(course.identifier), name: course.name, courseCode: course.courseCode,
            lessonCode: course.lessonCode, teacherName: course.teacherName,
            detailText: course.hasDetailText ? course.detailText : nil, credit: course.credit,
            dateTimePlacePersonText: course.hasDateTimePlacePersonText ? course.dateTimePlacePersonText : nil,
            lectures: lectures
        )
    }
}
