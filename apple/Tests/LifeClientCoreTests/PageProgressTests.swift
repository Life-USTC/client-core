import Foundation
import LifeClientCore
import LifeCore
import XCTest

final class PageProgressTests: XCTestCase {
    func testSwiftWrapperKeepsPreviousProgressAndCallerContinuationDecision() {
        let original = PageProgress(firstPage: 7)
        let snapshot = original
        let accepted = original.afterPage(itemCount: 100, hasMore: true)
        XCTAssertEqual(accepted.nextPage, 8)
        XCTAssertEqual(accepted.receivedCount, 100)
        XCTAssertEqual(original.nextPage, 7)
        XCTAssertEqual(original.receivedCount, 0)
        XCTAssertEqual(snapshot.nextPage, 7)
        XCTAssertEqual(snapshot.receivedCount, 0)

        let complete = accepted.afterPage(itemCount: 0, hasMore: false)
        XCTAssertEqual(complete.nextPage, 0)
        XCTAssertEqual(complete.receivedCount, 100)
        XCTAssertEqual(accepted.nextPage, 8)
        XCTAssertEqual(accepted.receivedCount, 100)

        let emptyIntermediate = original.afterPage(itemCount: 0, hasMore: true)
        XCTAssertEqual(emptyIntermediate.nextPage, 8)
        XCTAssertEqual(emptyIntermediate.receivedCount, 0)
    }

    private struct Fixture: Decodable {
        struct Page: Decodable {
            let itemCount: Int64
            let hasMore: Bool
            let nextPage: Int64
            let receivedCount: Int64
        }
        let name: String
        let firstPage: Int64
        let pages: [Page]
    }

    func testSharedProgressThroughCompiledGoBinding() throws {
        let url = try XCTUnwrap(Bundle.module.url(
            forResource: "page_progress_cases", withExtension: "json", subdirectory: "Fixtures"
        ))
        let fixtures = try JSONDecoder().decode([Fixture].self, from: Data(contentsOf: url))
        XCTAssertFalse(fixtures.isEmpty)
        for fixture in fixtures {
            var progress = try XCTUnwrap(LifeCoreMobileNewPageProgress(fixture.firstPage))
            XCTAssertEqual(progress.nextPage(), fixture.firstPage, fixture.name)
            XCTAssertEqual(progress.receivedCount(), 0, fixture.name)
            for page in fixture.pages {
                let previousPage = progress.nextPage()
                let previousCount = progress.receivedCount()
                let next = try XCTUnwrap(progress.afterPage(page.itemCount, hasMore: page.hasMore))
                XCTAssertEqual(next.nextPage(), page.nextPage, fixture.name)
                XCTAssertEqual(next.receivedCount(), page.receivedCount, fixture.name)
                XCTAssertEqual(progress.nextPage(), previousPage, fixture.name)
                XCTAssertEqual(progress.receivedCount(), previousCount, fixture.name)
                progress = next
            }
        }
    }
}
