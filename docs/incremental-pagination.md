# Incremental pagination PRD and roadmap

Status: implementation in progress; consumer integration is not yet verified.
The accepted direction is to show iOS lists and calendars while pages arrive.
Bot and CLI share accepted-page progress while keeping their existing fetching,
output and validation decisions.

Under Life-USTC/Life-USTC#109, start with the public Young list and calendar,
already consumed by all three clients. #100, #101 and #107 retain their separate
acceptance criteria.

## Experience

1. A list contains 300 activities. After two batches arrive, iOS shows those
   200 activities and says it is still loading. If the third request fails, the
   200 remain visible with an interruption message and an explicit Continue
   action. They are not presented as the complete list.
2. A calendar already has downloaded events. New results become visible as
   they arrive. An interruption keeps the available events on screen. Events
   that have not arrived yet are not treated as deleted from the saved complete
   calendar; replacing that saved collection still requires a complete read.

An error before the first batch means no data was loaded, not that there are
no matches. Only a successful complete empty read shows the empty result.
Adapters distinguish previous complete data from the current partial download.
A changed query, date range, backend or account starts a separate read.

Loading and interruption indicators sit alongside the visible content.
Continuation after interruption is explicit. Progress belongs to the active
read; this slice adds no background scheduler or automatic retries.

## Client behavior

| Consumer | Behavior in this slice |
| --- | --- |
| iOS Young list and public calendar | Publish accepted batches during download; retain visible progress on interruption; expose Continue. |
| Bot Young list | Keep its ten-item page and next-page command. |
| Bot Young calendar and organizer activities | Keep collecting before producing the existing command response. |
| CLI Young list and date views | Keep existing terminal output after collection. An explicit page/limit still requests one page. |
| CLI JSON and jq | Preserve one complete result document for a full-collection command. A later failure emits an error without a partial success document. Single-page commands retain their page envelope. |
| Other existing CLI `FetchAllPages` readers | Retain their complete result and failure contracts across workspace, calendar, exam, homework, upload, school and comment endpoints. |

Bot's delivery interface currently sends messages and cannot edit an existing
result. Its direct command response is buffered; the agent SendResponse hook
also does not establish immediate delivery. This slice adds neither message
editing nor a new message for each page. CLI tables measure the collected rows
before printing; interactive updates are a later presentation change.

## Shared boundary

Go's `pagination.Progress` records progress after the adapter accepts a page:
received item count and next page, with no next page meaning traversal has
reached the adapter's end.
The adapter supplies the existing has-more decision. This shares progress and
continuation, not the definition of a complete or valid server collection.

No interrupted/continue state machine is needed. On request failure, the
adapter retains the same progress and its error; Continue requests the same
next page. A successful accepted page returns a new progress value.

Adapters retain page/record validation, existing append or merge rules,
termination policy, HTTP, credentials, account ownership, cancellation,
stale-result checks, observable state, storage, localization and rendering.
The core performs no I/O and needs no planner, policy configuration or streaming
bridge. See the [implementation sketch](incremental-pagination-sketch.md).

## Existing complete operations

Displaying received records does not authorize replacing a complete saved
collection with them. Cache replacement, an operation advertised as a full
export, and selecting a sole catalog match still require complete data.
Interruption retains the previous complete collection; a successfully complete
empty result may clear it. These existing boundaries remain outside the initial
Young migration.

The clients' existing duplicate and total-change checks stay in place. This
slice neither adds new rejection rules nor requires agreement on rare changed
list responses before delivering incremental display. Each adapter computes
has-more as it does today; cursor endpoints remain outside the numbered-page
slice. A shared progress value alone never authorizes replacing saved data or
emitting a complete result through an existing validation gate.

## Delivery and acceptance

1. Specify accepted-page progress with shared fixtures: each batch advances
   count and page, an explicit starting page is retained, and the adapter can
   declare an empty or short page either intermediate or terminal. Existing
   adapter validation tests must keep their current expectations.
2. Implement the synchronous Go value and scalar Apple binding. Migrate the
   Bot Young loop and the CLI's existing shared reader without changing
   validation, append rules, source metadata, final completeness checks or
   output boundaries.
3. Release a matching Go module and Apple framework. Connect iOS list/calendar
   batches through the native adapter. Demonstrate the first batch on screen
   while a later request is deliberately delayed; then interrupt and Continue
   without clearing received content. Verify query changes and existing
   session/cancellation behavior through that adapter.
4. Run shared fixtures, Go consumer contracts, real Swift binding tests and
   native journeys. Review iPhone/iPad evidence and complete each repository's
   full local and protected-PR checks.

An unused core API or a fully buffered iOS call does not satisfy this slice.
Existing traversal and failure tests remain useful; incremental display and
continuation require their own behavioral evidence.
