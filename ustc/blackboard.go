package ustc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	BlackboardCalendarPageURL = "https://www.bb.ustc.edu.cn/webapps/calendar/viewPersonal"
	BlackboardCalendarsURL    = "https://www.bb.ustc.edu.cn/webapps/calendar/calendarData/calendars?mode=personal&course_id="
	BlackboardCoursesURL      = "https://www.bb.ustc.edu.cn/webapps/portal/execute/tabs/tabAction?tab_tab_group_id=_1_1"

	// BlackboardSelectedCalendarEventsURL retains the existing iOS request,
	// including its repeated query fragment. It does not select CLI calendars.
	BlackboardSelectedCalendarEventsURL = "https://www.bb.ustc.edu.cn/webapps/calendar/calendarData/selectedCalendarEvents?start=&end=&course_id=&mode=personal?start=&end=&course_id=&mode=personal"

	BlackboardCalendarIDsPerRequestMax = 1800
)

// BlackboardCalendarEventsURL preserves the CLI's complete calendar window.
func BlackboardCalendarEventsURL(calendarIDs []string) string {
	return "https://www.bb.ustc.edu.cn/webapps/calendar/calendarData/events?start=0&end=4102444800000&course_id=&calendarIds=" + strings.Join(calendarIDs, ",")
}

func BlackboardGradesURL(courseID string) string {
	return "https://www.bb.ustc.edu.cn/webapps/gradebook/do/student/viewGrades?course_id=" + url.QueryEscape(courseID)
}

// ParseBlackboardCalendarIDs removes personal/institution calendars and returns
// sorted unique identifiers, preserving the CLI calendar-list response profile.
func ParseBlackboardCalendarIDs(data []byte) ([]string, error) {
	var payload struct {
		Calendars []struct {
			ID string `json:"id"`
		} `json:"calendars"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(payload.Calendars))
	for _, calendar := range payload.Calendars {
		if calendar.ID != "" && calendar.ID != "PERSONAL" && calendar.ID != "INSTITUTION" {
			ids = append(ids, calendar.ID)
		}
	}
	slices.Sort(ids)
	return slices.Compact(ids), nil
}

// BlackboardCalendarEvent contains the fields read from the explicit-calendar
// events endpoint. Start and End remain unmodified strings for the CLI adapter.
type BlackboardCalendarEvent struct {
	ID           string `json:"id"`
	CalendarID   string `json:"calendarId"`
	ItemSourceID string `json:"itemSourceId"`
	Title        string `json:"title"`
	Start        string `json:"start"`
	End          string `json:"end"`
	EventType    string `json:"eventType"`
	CalendarName string `json:"calendarName"`
}

func ParseBlackboardCalendarEvents(data []byte) ([]BlackboardCalendarEvent, error) {
	var events []BlackboardCalendarEvent
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&events); err != nil {
		return nil, err
	}
	return events, nil
}

// BlackboardSelectedCalendarEvent contains only the existing iOS homework
// projection. The native adapter converts EndDate with its DateFormatter.
type BlackboardSelectedCalendarEvent struct {
	Title      string `json:"title"`
	CourseName string `json:"courseName"`
	EndDate    string `json:"endDate"`
}

// ParseBlackboardSelectedCalendarEvents preserves the iOS helper's arrayValue
// and stringValue behavior. It retains order, duplicates and raw dates, and does
// not substitute fields from the CLI event profile or infer completion status.
func ParseBlackboardSelectedCalendarEvents(data []byte) ([]BlackboardSelectedCalendarEvent, error) {
	value, err := decodeScalarJSON(data)
	if err != nil {
		return nil, fmt.Errorf("decode Blackboard selected calendar events: %w", err)
	}
	array, _ := value.([]any)
	events := make([]BlackboardSelectedCalendarEvent, 0, len(array))
	for _, value := range array {
		event := scalarJSONObject(value)
		calendar := scalarJSONObject(event["calendarNameLocalizable"])
		events = append(events, BlackboardSelectedCalendarEvent{
			Title:      scalarJSONText(event["title"]),
			CourseName: scalarJSONText(calendar["rawValue"]),
			EndDate:    scalarJSONText(event["endDate"]),
		})
	}
	return events, nil
}

// ParseBlackboardCourseIDs associates portal course labels with launcher IDs.
// Fetching the portal and deciding which grade pages to request stay in clients.
func ParseBlackboardCourseIDs(reader io.Reader) (map[string]string, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	out := map[string]string{}
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		rawCourseID := strings.TrimSpace(strings.Split(strings.TrimSpace(s.Text()), ":")[0])
		if rawCourseID == "" {
			return
		}
		parsed, err := url.Parse(strings.TrimSpace(href))
		if err != nil || parsed.Path != "/webapps/blackboard/execute/launcher" || parsed.Query().Get("type") != "Course" {
			return
		}
		if internalID := parsed.Query().Get("id"); internalID != "" {
			out[rawCourseID] = internalID
		}
	})
	return out, nil
}

// ParseBlackboardGradeStatuses retains the activity-status precedence and item
// ID projection used by the CLI. An unreadable document retains its nil result.
func ParseBlackboardGradeStatuses(reader io.Reader) map[string]string {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil
	}
	out := map[string]string{}
	doc.Find("#grades_wrapper .sortable_item_row").Each(func(_ int, s *goquery.Selection) {
		id := strings.TrimSpace(s.AttrOr("id", ""))
		if id == "" {
			return
		}
		status := strings.TrimSpace(s.Find(".cell.activity .activityType").First().Text())
		if status == "" {
			status = strings.TrimSpace(s.Find(".cell.gradeStatus").First().Text())
		}
		if status != "" {
			out["_"+id+"_1"] = strings.Join(strings.Fields(status), " ")
		}
	})
	return out
}

func SplitBlackboardCalendarID(calendarID string) (lessonCode, semesterCode string) {
	parts := strings.Split(strings.TrimSpace(calendarID), ".")
	if len(parts) < 3 {
		return strings.TrimSpace(calendarID), ""
	}
	semesterCode = parts[len(parts)-1]
	if len(semesterCode) != 6 {
		return strings.TrimSpace(calendarID), ""
	}
	for _, r := range semesterCode[:4] {
		if r < '0' || r > '9' {
			return strings.TrimSpace(calendarID), ""
		}
	}
	switch semesterCode[4:] {
	case "SP", "SU", "FA":
		return strings.Join(parts[:len(parts)-1], "."), semesterCode
	default:
		return strings.TrimSpace(calendarID), ""
	}
}

// BatchBlackboardCalendarIDs retains the CLI's comma-separated query budget.
// Network execution and per-course grade request caching remain in its adapter.
func BatchBlackboardCalendarIDs(calendarIDs []string, maxQueryLength int) [][]string {
	if len(calendarIDs) == 0 {
		return nil
	}
	if maxQueryLength <= 0 {
		return [][]string{append([]string(nil), calendarIDs...)}
	}
	batches := make([][]string, 0, 1)
	current := make([]string, 0, len(calendarIDs))
	currentLen := 0
	for _, calendarID := range calendarIDs {
		extraLen := len(calendarID)
		if len(current) > 0 {
			extraLen++
		}
		if currentLen > 0 && currentLen+extraLen > maxQueryLength {
			batches = append(batches, current)
			current = make([]string, 0, len(calendarIDs))
			currentLen = 0
		}
		current = append(current, calendarID)
		currentLen += extraLen
	}
	if len(current) > 0 {
		batches = append(batches, current)
	}
	return batches
}
