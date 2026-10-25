# Graduate YJS1 protocol extraction

This source layer belongs to [core #19](https://github.com/Life-USTC/client-core/issues/19)
under [Task #12](https://github.com/Life-USTC/client-core/issues/12). It extracts the
existing CLI YJS1 protocol into `ustc`; it adds no graduate iOS capability.
The iOS graduate course lookup uses `app.ustc.edu.cn` and static course details,
so that protocol remains a separate profile in the same package.

## API and consumer boundary

| Existing CLI path | Shared API |
| --- | --- |
| `graduateScheduleSession` / `fetchGraduateSemesters` | Schedule endpoint constants, `YJS1ScheduleTermsForm`, `ParseYJS1Terms` |
| `fetchGraduateCurriculum` | `YJS1ScheduleRowsForm`, `ParseYJS1Curriculum` |
| `fetchGraduateExams` | Exam endpoint constants, `YJS1ExamRowsForm`, `ParseYJS1Exams` |
| `fetchGraduateScores` | Score endpoint constants, `ParseYJS1ScoreInfo`, `ParseYJS1ScoreProfile`, `YJS1ScoreRowsForm`, `ParseYJS1Scores` |
| `fetchGraduateHomework` | Homework endpoint constants, `YJS1HomeworkTermsForm`, `YJS1HomeworkRowsForm`, `ParseYJS1UserID`, `ParseYJS1Homework` |

The caller owns authentication, credentials, HTTP methods/headers, response-status
checks, response-body reading, session reuse and request order. Score info and
profile parsing are separate so each response is validated before the next
request. The final score parser accepts those parsed values and the score rows.
Clients retain their existing empty form for the exam-term, score-info and
score-profile requests. Row requests retain page 1 and size 999; this extraction
does not add pagination or claim that upstream results are complete beyond it.

`ParseYJS1Terms` returns every term in source order with its original code, dates,
name fallback and current flag. CLI keeps `pickGraduateYJS1Term`: explicit ID,
otherwise current term, otherwise first term. This remains distinct from the
undergraduate catalog policy. Course projections take the already-selected
semester ID rather than selecting a semester themselves.

No YJS1 mobile bindings are added: the current iOS application does not consume
these endpoints. An unused bridge would not establish shared production use.
The CLI adapter must replace the actual import paths and pin the released module
before this layer is accepted; the parent Task remains open for both clients'
remaining protocol profiles and their integration acceptance.

## Preserved behavior

- JSON decoding accepts the first JSON value, keeps strict string fields and
  decodes mixed score values with the existing default `encoding/json` number
  behavior. An empty or `"0"` response code succeeds; other codes retain the
  current error text. The first map-iterated `datas` entry supplies rows, with no
  inferred key, merged sets or ordering guarantee across multiple data sets.
- Course IDs use the leading numeric part of `WID`; names and codes retain
  existing whitespace handling. Teachers retain delimiter splitting, sorting
  and deduplication. Schedule text pairs semicolon-separated times and weeks
  when counts match, otherwise preserves the existing combined display.
- Scores preserve numeric conversion, zero defaults, semester-name fallback,
  grade-detail deduplication and semester/course/lesson ordering. Summary fields
  prefer the detailed profile, then the first info row. Empty fields stay
  omitted; nonnumeric average-GPA text remains text.
- Exams preserve display-field precedence, raw start/end composition, optional
  seat/location/status fields and date/course/lesson ordering.
- Homework preserves submitted-ID precedence, raw date strings, deadline/course/
  title ordering and status precedence: graded, submitted, overdue, not started,
  pending. Numeric zero counts as a grade. The caller passes `time.Now()` and
  `time.Local` explicitly; the package does not read a clock or choose a timezone.
  Dates retain the three existing local-time layouts and strict before/after
  comparisons at the start/deadline boundaries.
- User-ID extraction retains the exact embedded `"USERID":"..."` pattern and
  original missing-ID error with the caller's app URL. CAS login is unchanged.

## Fixtures and verification

`ustc/testdata/yjs1_cases.json` contains 44 parser cases, seven form cases, six
user-ID cases and the 13 endpoint values. These are synthetic compatibility
examples from `CLI/internal/school/graduate.go` and its four existing projection
tests, not captured university account data or live-login evidence.
The cases cover the original semester, curriculum, score and homework examples;
field precedence, sorting, summary, empty results, response errors, exact query
JSON, first-value framing and explicit timezone/boundary behavior.

Run root and mobile uncached race-enabled Go tests and vet. Consumer acceptance
also requires CLI request-chain tests with the released module, including keeping
request error timing and session behavior at the adapter. The package has no
network or credential access and does not alter the separate server OAuth flow.
