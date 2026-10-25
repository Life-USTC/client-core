package mobile

import (
	"net/url"
	"sort"

	"github.com/Life-USTC/client-core/young"
)

// YoungEventQuery converts the public event filters to gomobile-supported
// values. Presence flags keep an omitted boolean distinct from false.
type YoungEventQuery struct {
	Search         string
	Category       string
	Module         string
	ActivityLevel  string
	OrganizerID    string
	HasActive      bool
	Active         bool
	HasDateUnknown bool
	DateUnknown    bool
	DateFrom       string
	DateTo         string
	TimeBasis      string
}

func NewYoungEventQuery() *YoungEventQuery {
	return &YoungEventQuery{}
}

// YoungEventParameters is the concrete parameter list for one Young request.
type YoungEventParameters struct {
	names  []string
	values url.Values
}

func (p *YoungEventParameters) Count() int64 {
	return int64(len(p.names))
}

func (p *YoungEventParameters) Name(index int64) string {
	return p.names[index]
}

func (p *YoungEventParameters) Value(index int64) string {
	return p.values.Get(p.names[index])
}

// Parameters preserves pagination presence without choosing a page or size.
func (q *YoungEventQuery) Parameters(page int64, hasPage bool, pageSize int64, hasPageSize bool) *YoungEventParameters {
	var pageValue, sizeValue *int
	if hasPage {
		value := int(page)
		pageValue = &value
	}
	if hasPageSize {
		value := int(pageSize)
		sizeValue = &value
	}
	values := q.query().Parameters(pageValue, sizeValue)
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	sort.Strings(names)
	return &YoungEventParameters{names: names, values: values}
}

// ForUnknownDates delegates the filter transition to the shared Young core.
func (q *YoungEventQuery) ForUnknownDates() *YoungEventQuery {
	value := q.query().ForUnknownDates()
	result := &YoungEventQuery{
		Search: value.Search, Category: value.Category, Module: value.Module,
		ActivityLevel: value.ActivityLevel, OrganizerID: value.OrganizerID,
		DateFrom: value.DateFrom, DateTo: value.DateTo, TimeBasis: string(value.TimeBasis),
	}
	if value.Active != nil {
		result.HasActive, result.Active = true, *value.Active
	}
	if value.DateUnknown != nil {
		result.HasDateUnknown, result.DateUnknown = true, *value.DateUnknown
	}
	return result
}

func (q *YoungEventQuery) query() young.EventQuery {
	result := young.EventQuery{
		Search: q.Search, Category: q.Category, Module: q.Module,
		ActivityLevel: q.ActivityLevel, OrganizerID: q.OrganizerID,
		DateFrom: q.DateFrom, DateTo: q.DateTo, TimeBasis: young.TimeBasis(q.TimeBasis),
	}
	if q.HasActive {
		value := q.Active
		result.Active = &value
	}
	if q.HasDateUnknown {
		value := q.DateUnknown
		result.DateUnknown = &value
	}
	return result
}
