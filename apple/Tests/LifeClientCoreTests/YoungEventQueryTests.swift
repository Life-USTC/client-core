import Foundation
import XCTest
@testable import LifeClientCore

final class YoungEventQueryTests: XCTestCase {
    private struct QueryInput: Decodable {
        let search: String?
        let category: String?
        let module: String?
        let activityLevel: String?
        let organizerId: String?
        let active: Bool?
        let dateUnknown: Bool?
        let dateFrom: String?
        let dateTo: String?
        let timeBasis: String?

        var query: YoungEventQuery {
            YoungEventQuery(
                search: search ?? "", category: category ?? "", module: module ?? "",
                activityLevel: activityLevel ?? "", organizerID: organizerId ?? "",
                active: active, dateUnknown: dateUnknown,
                dateFrom: dateFrom ?? "", dateTo: dateTo ?? "",
                timeBasis: timeBasis.flatMap(YoungTimeBasis.init(rawValue:))
            )
        }
    }

    private struct Fixtures: Decodable {
        struct Request: Decodable {
            let name: String
            let query: QueryInput
            let page: Int?
            let pageSize: Int?
            let parameters: [String: String]
        }
        struct UnknownDates: Decodable {
            let name: String
            let query: QueryInput
            let expectedQuery: QueryInput
            let parameters: [String: String]
        }
        let requests: [Request]
        let unknownDates: [UnknownDates]
    }

    private func fixtures() throws -> Fixtures {
        let url = try XCTUnwrap(Bundle.module.url(
            forResource: "young_query_cases", withExtension: "json", subdirectory: "Fixtures"
        ))
        let values = try JSONDecoder().decode(Fixtures.self, from: Data(contentsOf: url))
        XCTAssertFalse(values.requests.isEmpty)
        XCTAssertFalse(values.unknownDates.isEmpty)
        return values
    }

    func testSharedYoungQueryParametersThroughCompiledGoBinding() throws {
        for fixture in try fixtures().requests {
            let items = fixture.query.query.queryItems(page: fixture.page, pageSize: fixture.pageSize)
            let actual = Dictionary(uniqueKeysWithValues: items.map { ($0.name, $0.value!) })
            XCTAssertEqual(actual, fixture.parameters, fixture.name)
            XCTAssertEqual(items.map(\.name), fixture.parameters.keys.sorted(), fixture.name)
        }
    }

    func testSharedUnknownDateTransitionThroughCompiledGoBinding() throws {
        for fixture in try fixtures().unknownDates {
            let query = fixture.query.query
            let actual = query.forUnknownDates()
            XCTAssertEqual(actual, fixture.expectedQuery.query, fixture.name)
            XCTAssertEqual(query, fixture.query.query, fixture.name)
            let values = Dictionary(uniqueKeysWithValues: actual.queryItems().map { ($0.name, $0.value!) })
            XCTAssertEqual(values, fixture.parameters, fixture.name)
        }
    }
}
