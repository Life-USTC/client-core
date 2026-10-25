import LifeCore

/// Accepted row count and next page, with traversal policy supplied by the caller.
public struct PageProgress {
    private let binding: LifeCoreMobilePageProgress

    public init(firstPage: Int) {
        binding = LifeCoreMobileNewPageProgress(Int64(firstPage))!
    }

    private init(binding: LifeCoreMobilePageProgress) {
        self.binding = binding
    }

    public var nextPage: Int { Int(binding.nextPage()) }
    public var receivedCount: Int { Int(binding.receivedCount()) }

    public func afterPage(itemCount: Int, hasMore: Bool) -> PageProgress {
        PageProgress(binding: binding.afterPage(Int64(itemCount), hasMore: hasMore)!)
    }
}
