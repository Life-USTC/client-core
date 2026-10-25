import Foundation
import LifeCore

/// Request and response conversion for the shared undergraduate course-table logic.
public enum USTCCourseTable {
    public static var pageURL: URL {
        URL(string: LifeCoreMobileUSTCCourseTablePageURL())!
    }

    public static func dataURL(semesterID: Int, studentID: String) -> URL {
        URL(string: LifeCoreMobileUSTCCourseTableDataURL(Int64(semesterID), studentID))!
    }

    public static func studentID(from url: URL) throws -> String {
        var error: NSError?
        let value = LifeCoreMobileUSTCStudentIDFromCourseTableURL(url.absoluteString, &error)
        if let error { throw error }
        return value
    }

    public static func parse(_ data: Data) throws -> [Int] {
        var error: NSError?
        let value = LifeCoreMobileParseUSTCCourseTable(data, &error)
        if let error { throw error }
        // The Go function returns a collection whenever its error is nil.
        return (0..<value!.count()).map { Int(value!.at($0)) }
    }
}
