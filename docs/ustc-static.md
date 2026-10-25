# Static curriculum and graduate app lookup

Core #19, under the shared extraction Task #12, moves the current iOS static
curriculum projection and graduate app course-name matching into `ustc`.
The app endpoint is independent of the CLI YJS1 protocol. This source slice does
not change login, add graduate academic features or complete native adoption.

The existing callers are `USTCCurriculumHelper.updateCourse`,
`updateGraduateCurriculum`, `USTCSchool.updateCurriculumFromUSTC`, and both
additional-course view loaders. The fixture-only `USTCStaticService` is not the
behavioral baseline: the actual helpers and views read SwiftyJSON values.

## Preserved projection

`ParseStaticSemesters`, `ParseStaticCourseList` and `ParseStaticCourse` preserve
array order, duplicate IDs, per-entry defaults and raw epoch seconds. Native
adapters still construct Dates, select the current semester, upsert SwiftData
objects, replace lectures and merge user-selected extra course IDs.

The scalar accessors follow pinned SwiftyJSON 5.0.2 (`af76cf3`):

- `stringValue` keeps strings and scalar number/bool text; missing/null/compound
  fields become empty strings. Optional `.string` preserves absent versus empty.
- Nonoptional `doubleValue` and `intValue` use `numberValue`. They accept ordinary
  numeric strings, including whitespace and the decimal prefix accepted by
  Foundation's NSDecimalNumber; invalid strings keep the zero default.
- Optional `.int` uses `number` and does not parse strings. Lecture index `"0"`
  stays absent, while numeric zero and false are present zero values.
- `arrayValue` accepts arrays only. Object roots produce an empty list, while
  an array course-detail root produces the same default course as the helper.
  Malformed JSON and non-object/non-array fragments remain errors.

The private defaulting-number helper is deliberately small: it handles the
existing field formats without a dependency or general NSDecimalNumber emulator
for pathological decimal ranges. Optional score helpers moved to `scalar_json.go`
unchanged and retain their prior behavior. The static fixtures distinguish the
optional and defaulting paths, including numeric strings, booleans and nil/zero.

`detailText` has no `description` fallback. The static producer currently emits
`description`, while the importing helper reads `detailText`; its absent result
is preserved. Dates are not parsed, localized, regenerated or rescaled.

## Graduate matching

`ParseGraduateAppCourseCodes` reads only the `d.lists` array and `course_name`.
It keeps the first ASCII parenthesis capture from each name, including empty
captures, whitespace and repeated codes. Full-width parentheses are not
normalized, and a parenthesis pair spanning a newline is skipped, matching the
existing Swift regex. `MatchStaticCourses` tests canonical `lessonCode` OR raw
`lesson_code`, preserving static-list order and repeated IDs. An empty captured
code can match a missing/default field. No captured code produces no match.
The additional-course view continues to display the canonical field only.

Native code keeps the two existing CAS/session-establishing requests and the
shared URLSession. The core exposes the app lookup URL and static relative paths;
it performs no requests or credential access. The mobile parser returns an
immutable typed app-course handle before the native adapter fetches the static
list; its matching method then parses that list and returns the existing ID
collection. This preserves the app-validation-before-next-request error boundary
without a public dynamic JSON or string-collection API.

## Evidence and release order

`ustc/testdata/static_cases.json` contains synthetic contracts from the real
helpers/views and pinned scalar accessors. A bounded macOS Foundation/Swift probe
confirmed normal numeric-string conversions, invalid-string defaults and regex
newline behavior. All 38 fixture cases were also compared with the pinned
SwiftyJSON source and original Swift regex on macOS. Go tests cover every
projection and app-code case. Mobile and raw compiled Swift binding tests consume
the same semester, list, detail, matching and request-path fixtures.

Run full root/mobile race tests and vet, then the device/simulator/host framework
build and compiled Swift tests. Ship the source before adding public Swift
wrappers atomically with the matching released archive/checksum. The subsequent
iOS change must replace all actual helper/view parsers through the native service,
retain persistence/date/extra-course semantics, and verify re-import behavior.
These source fixtures do not establish university-account or production OAuth
acceptance, and the parent Task remains open until actual consumers are accepted.
