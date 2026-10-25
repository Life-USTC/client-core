# Undergraduate catalog and score metadata

This extraction completes the remaining CLI catalog protocol slice of core #19
under Task #12. It adds no iOS catalog capability or mobile bindings.

`CatalogSemestersURL` and `CatalogLessonsURL(semesterID, page, pageSize)` retain
the existing paths. The CLI continues to request 500 rows per page, starting at
page 1. Existing `ScoreSemestersURL` and `ScoreTypesURL` remain unchanged.

`ParseUndergraduateSemesters` handles the exact seven-field array consumed by
both catalog semester listing and JW score metadata. `ParseScoreSheetTypes`
retains the integer ID array. Both consume the first JSON value and preserve
source order, duplicates, zero values, strict field types and nil versus empty
arrays. Neither chooses a semester, display name or transcript type.

`ParseCatalogLessonPage` first tries a direct `[]Lesson`, then a
`{rows,totalPage}` object, using strict whole-document JSON decoding for both
attempts. It reuses the existing `Lesson` fields and localized-text decoding.
The `Direct` flag also covers JSON null, which succeeds as a nil direct array.
The CLI preserves its existing behavior: a direct response returns immediately
and replaces any previously accumulated pages; a page appends rows, stopping
when totalPage is zero or the requested page reaches it. The parser neither
merges pages nor changes that policy.

Authentication, HTTP headers/status/body handling, error prefixes, catalog-to-JW
fallback and curriculum/semester selection remain in CLI. Parser errors preserve
their JSON error types; the native adapter supplies endpoint context.

`ustc/testdata/catalog_cases.json` contains synthetic fixtures derived from the
current CLI code. Tests cover both lesson shapes, nil/empty/order/duplicates,
strict versus first-value framing, typed fields and exact URL construction.
Acceptance also requires the real CLI consumer and its request-chain tests to
use the released Go module; source tests alone do not complete Task #12.
