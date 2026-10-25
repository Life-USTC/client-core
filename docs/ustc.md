# Shared USTC logic PRD and roadmap

[Task #12](https://github.com/Life-USTC/client-core/issues/12) tracks a shared
`ustc` Go package inside the existing client-core module. CLI imports the typed
Go API; iOS consumes the same implementation through the existing mobile
binding, XCFramework and `LifeClientCore` Swift package.

The outcome is one maintained implementation of the university-specific request
and response rules used by both clients. Each slice must replace a real
production call path in both clients before it is accepted. A package, bridge or
passing fixture suite without those consumers does not complete the outcome.

## Existing capabilities and boundaries

CLI already offers `workspace school` for undergraduate and graduate semesters,
curriculum, exams, scores and homework, plus synchronization to Life@USTC.
Moving this code is extraction, not a new CLI feature. iOS has undergraduate
curriculum, exams, scores and Blackboard homework; its graduate path currently
imports curriculum only. Bringing graduate exams, scores or homework to iOS is
an additional product capability and needs its own acceptance.

| Owner | Responsibilities |
| --- | --- |
| `ustc` | University endpoint selection, request construction, response decoding and shared data projections as each slice migrates |
| CLI | Commands, configuration and `.env` credentials, HTTP/session ownership, output, semester-selection policy and server synchronization |
| iOS | WebKit login, Keychain, URLSession/cookies, user-selected extra courses, SwiftData, SwiftUI and Apple lifecycle |
| Server | Life@USTC OAuth, permissions, canonical records and persisted mutations |

University CAS and Life@USTC server OAuth are separate authentication systems.
Server production OAuth acceptance is outside this extraction. Existing CLI CAS
code submits the username/password form but reports an error for an OTP
challenge; accepting a TOTP argument does not establish complete OTP support.
The first slice does not change either client's login behavior.

The iOS `USTC*Service` classes initially had fixture tests but no production
callers. Actual imports lived in `Schools/USTC/Helpers`. Migration must reach
those helpers, reusing the existing service where it provides the needed native
transport, and remove the duplicated request/decoding logic in the same change.

## First slice: undergraduate course table

Tracked by [core #13](https://github.com/Life-USTC/client-core/issues/13),
[CLI #98](https://github.com/Life-USTC/cli/issues/98) and
[iOS #130](https://github.com/Life-USTC/Life-USTC/issues/130).

The shared functions cover:

- The JW course-table page URL and the selected-semester `get-data` URL.
- Extracting the numeric student identifier from `/course-table/<id>` or
  `/course-table/info/<id>` in the final redirect URL.
- Decoding selected lesson identifiers and any embedded curriculum details.

`ParseCourseTable` preserves the CLI's established response selection:
nonempty `studentTableVm` fields take priority, followed by top-level fields.
When neither contains lesson details, `lessonId2Lesson` supplies matching
details in sorted identifier order. Identifiers are returned sorted and unique.
The map lookup retains the existing duplicate detail entries when the source
contains repeated identifiers. List-provided detail order is unchanged.

The result preserves lesson, embedded course, teacher and localized schedule
fields needed by the CLI's existing curriculum projection. It does not replace
the iOS static course/lecture feed with the CLI's shorter schedule text.

An empty or missing identifier list produces an empty table. CLI retains its
existing no-lesson error and semester-search behavior; iOS can still import
user-selected extra courses. Response decoding errors remain errors through the
compiled Apple binding. The package performs no network or credential access.

The two production paths are CLI `fetchStudentID` / `fetchCurrentCourseTable`
and iOS `updateUnderGraduateCurriculum`, through `USTCAASService`. Swift converts
its string semester identifier to an integer at the adapter. The mobile bridge
exposes the identifier collection currently needed by iOS; it does not add an
unused parallel Swift model for the embedded CLI lesson details.

## Remaining extraction sequence

1. **Undergraduate course table in both clients.** Ship this slice and the
   versioned framework, preserving the existing static lecture import and CLI
   curriculum output. Keep the parent Task open.
2. **Undergraduate catalog, scores and exams.** Establish representative,
   redacted response fixtures before unifying differing schemas. Extract their
   request and decoding rules, then replace each actual consumer.
3. **Blackboard homework.** Reconcile the CLI calendars/events and course/grade
   association path with the iOS selected-calendar-events path. Adding statuses
   or new course associations to iOS requires observable acceptance alongside
   the extraction; server synchronization remains a client use case.
   The [Blackboard extraction](ustc-blackboard.md) preserves the existing
   endpoint-specific responses and does not add those iOS capabilities.
4. **Graduate YJS1.** Move the existing CLI requests and projections, and connect
   the agreed iOS graduate capabilities. Its current `app.ustc.edu.cn` course
   lookup differs from CLI `yjs1.ustc.edu.cn`; this is not a mechanical endpoint
   replacement. Retain lecture detail and semester identity requirements.
   The [static/app extraction](ustc-static.md) retains the existing iOS scalar
   projection, dated lectures and canonical-or-legacy course matching.
5. **Remaining university protocol rules.** Share portable CAS/service request
   and response rules where both clients need them, with representative login
   responses and native integration evidence. Keep platform interaction,
   credential storage and session lifecycle at the edge. Any additional OTP
   capability must be specified and verified explicitly.

Known contract differences to resolve before their slices:

| Area | Current difference |
| --- | --- |
| JW exams | CLI maps the first three table columns to course name, lesson code and type; iOS maps lesson code, type and course name. CLI's later columns are location/seat/status; iOS keeps room/building/district/detail. Existing fixtures cannot establish the actual upstream table shape. |
| JW scores | CLI expects numeric credits/GPA and composes ranking text. iOS's legacy parser reads string-valued entries, preserves absent GPA, and needs separate major/rank/count fields. |
| Curriculum | CLI exposes course metadata and schedule text. iOS persists dated lecture occurrences from the static feed, plus user-selected extra courses. |
| Blackboard | CLI combines explicit calendar IDs, events and grade status pages; iOS imports selected events with localized course names and end dates. |
| Graduate | CLI uses YJS1 term codes and schedule/score/exam/homework endpoints. iOS currently uses the mobile app course-name lookup and static course matching. |

These differences require source-backed fixtures and explicit consumer behavior,
not speculative schema fallbacks. Existing fixtures are synthetic contract
examples, not evidence of a successful live university login.

## Verification and completion

- Shared fixtures specify request URLs, redirect identifiers, response
  precedence, empty tables, duplicate identifiers, missing details and malformed
  responses. Go and the real compiled Swift binding consume the same fixtures.
- Go tests additionally verify that course, teacher and localized schedule
  fields survive decoding for the CLI adapter.
- Run uncached race-enabled tests and `go vet` in the root and mobile modules.
  Build device, simulator and host framework slices; run the Swift bridge suite
  on macOS and the iOS simulator.
- Run the CLI's full contract suite and the iOS native checks with the shared
  implementation on the actual production paths. Verify the iOS extra-course
  import and the CLI's no-lesson semester policy at their owning adapters.
- Consumers pin the released Go module and matching framework distribution.
  Local replacements and an unused binary do not satisfy integration acceptance.
- Preserve the existing two-PR release sequence: the source change tests new
  symbols directly through the generated `LifeCore` binding while the root Swift
  package still builds against its current release. Publish the framework from
  the accepted source commit, then add the Swift wrapper and its tests atomically
  with the root package's new archive URL and checksum. This keeps every commit
  buildable and retains the released-package CI check.
- Keep native/university-account acceptance explicit when interaction is still
  outstanding. Completion of the first slice does not complete the whole
  university extraction Task.
