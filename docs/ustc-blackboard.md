# Shared Blackboard homework parsing

[Core #16](https://github.com/Life-USTC/client-core/issues/16),
[CLI #100](https://github.com/Life-USTC/cli/issues/100) and
[iOS #132](https://github.com/Life-USTC/Life-USTC/issues/132) belong to the
[shared USTC Task #12](https://github.com/Life-USTC/client-core/issues/12).
This slice extracts existing homework request and response rules. Completion
requires both production consumers and the matching released Go/Apple package.

## Two existing response profiles

| Consumer | Source | Preserved result |
| --- | --- | --- |
| CLI | Explicit calendars and calendar events | Event ID, calendar ID, item source ID, title, start, end, event type and calendar name |
| CLI | Portal course links and grade pages | Calendar-to-internal-course association and item-to-status association |
| iOS | Selected calendar events | Title, `calendarNameLocalizable.rawValue` and raw `endDate` |

`ParseBlackboardCalendarIDs` retains sorted unique course calendars, excluding
empty, `PERSONAL` and `INSTITUTION` IDs. `BlackboardCalendarEventsURL` retains the
CLI window from 0 through 4102444800000 milliseconds. The existing batching
budget and calendar identity split are exposed as pure helpers. Event parsing
keeps the eight CLI fields, input order, duplicate IDs and unmodified time
strings. It retains the original decoder behavior, including missing/null
string fields and rejection of non-string values in those fields.

`ParseBlackboardCourseIDs` reads course-launcher links and the label before the
first colon. `ParseBlackboardGradeStatuses` prefers the first activity type,
then the first grade-status cell, collapses whitespace and maps row IDs to
`_<id>_1`. Duplicate associations keep the last value, as before. These functions
read supplied HTML; they do not request pages or interpret statuses as complete.
The grade parser retains its existing nil result for a reader failure.

`ParseBlackboardSelectedCalendarEvents` retains the iOS helper's SwiftyJSON
profile through the shared scalar conversion helpers. A top-level object has
an empty `arrayValue`; fragments, malformed JSON and trailing JSON are errors.
Missing or non-scalar fields project to empty strings, numbers and booleans use
the existing scalar text conversion, and string contents remain unchanged.
Array order, duplicates and invalid date strings survive parsing. This parser
does not fall back to CLI `end` or `calendarName`, filter dates, or infer status.
The selected-events URL keeps the existing repeated query fragment.

## Consumer boundaries

CLI retains HTTP clients, authentication, response-status handling, request
batch execution and the per-course grade-page cache. Its existing `FetchHomework`
path uses the parsed associations to override an event type only when a matching
nonempty grade status exists, deduplicates by event ID and sorts the final
`HomeworkItem` list. Command output and Life@USTC synchronization remain CLI
responsibilities, including their current interpretation of completion status.

iOS connects `USTCSchool.updateHomework` to the native `USTCBlackboardService`
and the compiled shared parser. The service keeps URLSession/cookies; the
helper keeps login, demo handling and SwiftData. Foundation retains the same
UTC DateFormatter and skips dates it cannot parse. The helper retains replacement
of old homework and title-plus-course-name upsert order. Stored values remain
title, course name and due date; the current time-derived `isFinished` behavior
does not become Blackboard submission status.

The mobile binding exposes only the selected-events URL and a concrete
collection with `Count`, `Title`, `CourseName` and `EndDate` accessors. It performs
conversion only. The Swift wrapper converts that collection to value types
without another JSON parser or a new date parser.

## Evidence and delivery

`ustc/testdata/blackboard_cases.json` is a synthetic compatibility fixture,
derived from the two existing implementations and their tests. It covers request
URLs, filtering, identity splitting, batching, field preservation, HTML selectors,
duplicates, scalar conversion, missing values and invalid input. Go core tests
cover both profiles; mobile and raw Swift binding tests consume the same selected
event cases. Native consumer tests must also verify date filtering and the
persisted three-field result. No fixture implies successful university login.

Ship source functions and raw compiled-binding tests first. Build and publish
the matching framework, then introduce the public Swift wrapper together with
the root package archive/checksum update. Consumers pin the resulting version.
The core source alone does not complete the parent Task or replace native
integration and production-path acceptance.
