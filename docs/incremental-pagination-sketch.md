# Accepted-page progress sketch

Draft consumer wiring; this document is not evidence of released integration.

The shared value holds only the next requested page and the number of received
items. `pagination.NewProgress(firstPage)` starts at an explicit positive page.
`AfterPage(itemCount, hasMore)` returns a new value after the adapter accepts
that page. `NextPage() == 0` means the adapter reported no further page.
`ReceivedCount()` counts supplied rows, not distinct IDs or the server's total.
Calls use counts from decoded arrays and already validated page numbers.

```go
progress := pagination.NewProgress(1)
for progress.NextPage() != 0 {
    page, err := fetch(progress.NextPage())
    if err != nil { return err } // Preserve adapter's current error contract.
    // Existing adapter validation, record append and has-more calculation.
    progress = progress.AfterPage(len(page.Data), hasMore)
}
// Existing final collection validation still runs before output or persistence.
```

The has-more input preserves three actual Young consumers: iOS compares with
the first response's totalPages; CLI keeps its fixed first bound and strict
validation; Bot keeps `youngPageHasNext` on the current response. Their requested
page advances by one in each case. This is a shared progress primitive, not
unification of their termination or record-merge policies.

The CLI already shares `FetchAllPages` across Young, workspace, calendar,
exam, homework, upload, school and comment readers. The progress value belongs
in `pagination`, so those existing callers use it without a Young dependency.
Each retains its complete-output and response-validation boundary.

For iOS, the native model stores progress beside received rows. The current
all-pages `fetchYoungEventsPage` needs a one-page transport operation reused by
the full reader and the incremental display. The model's existing generation
and cancellation checks run before publishing each batch. In schematic Swift:

```swift
var progress = PageProgress(firstPage: 1)
while progress.nextPage != 0 {
    let requestedPage = progress.nextPage
    let batch = try await fetchOnePage(query, page: requestedPage)
    guard requestStillCurrent else { return }
    // Keep the existing first-response page bound and append/meta rules.
    rows.append(contentsOf: batch.data)
    progress = progress.afterPage(itemCount: batch.data.count, hasMore: hasMore)
    // Observable rows are now visible before awaiting the next request.
}
```

The Swift value wraps the generated mobile binding through `LifeClientCore`.
It ships with the matching framework distribution update. Networking,
generation checks and screen state remain native. If a later request throws,
keep rows/progress and show the adapter's error. Continue re-enters this loop;
it does not advance a page or clear rows. Query changes create fresh progress.
Bot retains its existing next-page command; CLI/Bot complete readers keep
their current buffered result and failure behavior.

Prove display timing with a controlled second response, not a sleep: release
page one, hold page two, and assert first-page rows plus loading state in both
list and calendar. Fail page two, assert retained rows and next page two, then
Continue and verify the captured request remains page two. Release a successful
terminal response and verify combined content. Core/binding tests only prove
the value transitions; these native tests prove the requested experience.
