# Undergraduate scores and exams extraction

This slice belongs to [Task #12](https://github.com/Life-USTC/client-core/issues/12),
with [core #15](https://github.com/Life-USTC/client-core/issues/15),
[CLI #99](https://github.com/Life-USTC/cli/issues/99) and
[iOS #131](https://github.com/Life-USTC/Life-USTC/issues/131).

The existing clients interpret different score and exam formats. Their explicit
parsers can share a package and typed results without guessing a common upstream
schema. No live university account is required to extract these existing rules;
the fixtures characterize client behavior, not a successful university login.

## Scores

`ParseNumericScores` retains the CLI's strict numeric JSON fields, fallback
precedence, numeric semester sort, zero for absent grade points and summary JSON.
Its semester-name lookup comes from the caller's existing semester request.
`NumericScoreListURL` keeps the explicit selected semester list.

`ParseScalarScores` retains the iOS helper's response order and scalar conversion.
Strings keep their whitespace. Numbers convert to their scalar text; booleans
convert to `true` or `false`; absent, array and object values produce empty text.
Entry credits use that text with zero on failed numeric conversion; an absent or
unparseable grade point remains nil, distinct from an explicit zero.

Metadata follows the helper's separate optional accessors: numeric values and
booleans can supply GPA/rank numbers, strings cannot; the major name must be a
string. The result keeps metadata presence and raw major/rank/count fields so
Swift retains its existing defaults, including the missing-major-name label.
`ScalarScoreListURL` preserves the existing all-semesters request.

The scalar rules were checked against SwiftyJSON 5.0.2, revision
`af76cf3ef710b6ca5f8c05f3a31307d44a3c5828`, used by native acceptance. A Foundation
probe on the development Mac confirmed that JSON numbers `6.0`, `95`, `1e2`,
`0.0`, `1.25` and `1e-7` produce `6`, `95`, `100`, `0`, `1.25` and `1e-07`.
Shared fixtures exercise both these values and SwiftyJSON's distinct boolean
text/number accessors. The package-private scalar helpers are also reusable by
the separately tracked Blackboard extraction when it adopts the same rules.

The unused CLI rendered-score HTML parser and the unused iOS Codable score DTOs
do not define either live JSON path. This slice does not introduce them as
fallbacks or remove unrelated pre-existing code.

## Exams

Both clients request `ExamArrangementURL`, but parser choice remains explicit:

| Parser | Existing behavior |
| --- | --- |
| `ParseExamStatusTable` | CLI: all `table tbody tr`, descendant cell text, course/lesson/type/date/location/seat/status columns; sort by date, course and lesson. |
| `ParseExamRoomTable` | iOS service: `#exams > tbody > tr`, direct cell text, lesson/type/course/date/room/building/district/detail columns; preserve response order. |

The shared `Exam` keeps these fields separate. Neither parser detects or switches
formats. Room parsing follows the existing typed service's eight-cell minimum;
the CLI's four-cell minimum and optional later columns remain distinct. The iOS
service's UTF-8 requirement and direct-text treatment of nested tags are retained.

Date strings pass through unchanged, including unknown dates. The existing iOS
DateFormatter, local timezone, start/end extraction and date-based import filtering
remain in Swift. The server/static catalog exam JSON is a third format and cannot
establish the column order of either JW HTML table.

## Acceptance

The shared fixture contains 16 score cases and 8 exam cases, covering profile
selection, ordering, metadata presence, nil versus zero GPA, string/number/bool
conversion, malformed responses, separate room fields and direct versus
descendant text. Go tests cover both profiles; mobile and compiled Swift tests
cover the actual iOS profiles without exporting unused CLI formats to Apple.

The source PR tests the generated `LifeCore` symbols directly. The subsequent
distribution PR adds Swift wrappers and their tests with the matching binary
URL/checksum, keeping the root released-package build valid. CLI and iOS must
replace their actual score/exam paths and pass their complete suites before this
slice is complete. Authentication, persistence and server OAuth remain outside it.
