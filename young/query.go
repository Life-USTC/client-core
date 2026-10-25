// Package young contains shared query projections for the public Young catalog.
package young

import (
	"net/url"
	"strconv"
	"strings"
)

// TimeBasis selects the event timestamps used by the server's date filters.
type TimeBasis string

const (
	Activity     TimeBasis = "activity"
	Registration TimeBasis = "registration"
)

// EventQuery contains the filters for the public Young event list. Empty text
// is omitted; nil booleans are distinct from an explicit false. Adapters retain
// input validation, pagination policy and HTTP execution.
type EventQuery struct {
	Search        string
	Category      string
	Module        string
	ActivityLevel string
	OrganizerID   string
	Active        *bool
	DateUnknown   *bool
	DateFrom      string
	DateTo        string
	TimeBasis     TimeBasis
}

// Parameters projects one request without changing the filters. Nil page or
// pageSize values are omitted; supplied values are neither defaulted nor capped.
func (q EventQuery) Parameters(page, pageSize *int) url.Values {
	values := url.Values{}
	for key, value := range map[string]string{
		"search":        q.Search,
		"category":      q.Category,
		"module":        q.Module,
		"activityLevel": q.ActivityLevel,
		"organizerId":   q.OrganizerID,
		"dateFrom":      q.DateFrom,
		"dateTo":        q.DateTo,
		"timeBasis":     string(q.TimeBasis),
	} {
		if value = strings.TrimSpace(value); value != "" {
			values.Set(key, value)
		}
	}
	if q.Active != nil {
		values.Set("active", strconv.FormatBool(*q.Active))
	}
	if q.DateUnknown != nil {
		values.Set("dateUnknown", strconv.FormatBool(*q.DateUnknown))
	}
	if page != nil {
		values.Set("page", strconv.Itoa(*page))
	}
	if pageSize != nil {
		values.Set("pageSize", strconv.Itoa(*pageSize))
	}
	return values
}

// ForUnknownDates keeps the selected basis and scalar filters while requesting
// events with an unknown start. It removes both date bounds from the copy.
func (q EventQuery) ForUnknownDates() EventQuery {
	unknown := true
	q.DateUnknown = &unknown
	q.DateFrom = ""
	q.DateTo = ""
	return q
}
