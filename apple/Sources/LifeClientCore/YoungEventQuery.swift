import Foundation
import LifeCore

public enum YoungTimeBasis: String, Sendable {
    case activity, registration
}

/// Public Young filters with request projection delegated to the Go core.
public struct YoungEventQuery: Equatable, Sendable {
    public var search: String
    public var category: String
    public var module: String
    public var activityLevel: String
    public var organizerID: String
    public var active: Bool?
    public var dateUnknown: Bool?
    public var dateFrom: String
    public var dateTo: String
    public var timeBasis: YoungTimeBasis?

    public init(
        search: String = "", category: String = "", module: String = "",
        activityLevel: String = "", organizerID: String = "", active: Bool? = nil,
        dateUnknown: Bool? = nil, dateFrom: String = "", dateTo: String = "",
        timeBasis: YoungTimeBasis? = nil
    ) {
        self.search = search
        self.category = category
        self.module = module
        self.activityLevel = activityLevel
        self.organizerID = organizerID
        self.active = active
        self.dateUnknown = dateUnknown
        self.dateFrom = dateFrom
        self.dateTo = dateTo
        self.timeBasis = timeBasis
    }

    /// Nil pagination fields are omitted; the caller chooses each requested page.
    public func queryItems(page: Int? = nil, pageSize: Int? = nil) -> [URLQueryItem] {
        let parameters = binding().parameters(
            Int64(page ?? 0), hasPage: page != nil,
            pageSize: Int64(pageSize ?? 0), hasPageSize: pageSize != nil
        )!
        return (0..<parameters.count()).map { index in
            URLQueryItem(name: parameters.name(index), value: parameters.value(index))
        }
    }

    public func forUnknownDates() -> YoungEventQuery {
        let value = binding().forUnknownDates()!
        return YoungEventQuery(
            search: value.search, category: value.category, module: value.module,
            activityLevel: value.activityLevel, organizerID: value.organizerID,
            active: value.hasActive ? value.active : nil,
            dateUnknown: value.hasDateUnknown ? value.dateUnknown : nil,
            dateFrom: value.dateFrom, dateTo: value.dateTo,
            timeBasis: YoungTimeBasis(rawValue: value.timeBasis)
        )
    }

    private func binding() -> LifeCoreMobileYoungEventQuery {
        let value = LifeCoreMobileNewYoungEventQuery()!
        value.search = search
        value.category = category
        value.module = module
        value.activityLevel = activityLevel
        value.organizerID = organizerID
        value.hasActive = active != nil
        value.active = active ?? false
        value.hasDateUnknown = dateUnknown != nil
        value.dateUnknown = dateUnknown ?? false
        value.dateFrom = dateFrom
        value.dateTo = dateTo
        value.timeBasis = timeBasis?.rawValue ?? ""
        return value
    }
}
