package mobile

import "github.com/Life-USTC/client-core/ustc"

func USTCBlackboardSelectedCalendarEventsURL() string {
	return ustc.BlackboardSelectedCalendarEventsURL
}

// USTCBlackboardSelectedCalendarEvents exposes the fields used by the iOS import.
type USTCBlackboardSelectedCalendarEvents struct {
	values []ustc.BlackboardSelectedCalendarEvent
}

func (events *USTCBlackboardSelectedCalendarEvents) Count() int64 {
	return int64(len(events.values))
}

func (events *USTCBlackboardSelectedCalendarEvents) Title(index int64) string {
	return events.values[index].Title
}

func (events *USTCBlackboardSelectedCalendarEvents) CourseName(index int64) string {
	return events.values[index].CourseName
}

func (events *USTCBlackboardSelectedCalendarEvents) EndDate(index int64) string {
	return events.values[index].EndDate
}

func ParseUSTCBlackboardSelectedCalendarEvents(data []byte) (*USTCBlackboardSelectedCalendarEvents, error) {
	events, err := ustc.ParseBlackboardSelectedCalendarEvents(data)
	if err != nil {
		return nil, err
	}
	return &USTCBlackboardSelectedCalendarEvents{values: events}, nil
}
