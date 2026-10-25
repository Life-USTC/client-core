# Shared client core PRD and roadmap

The shared Go core gives iOS, Bot and CLI one implementation of client-side
calculations and use cases. The same source runs in Linux tests, in Go clients,
and inside the iOS app through a compiled XCFramework. Swift owns SwiftUI and
Apple integration. The server remains authoritative for permissions, canonical
records and mutation semantics.

Tracked by Life-USTC/Life-USTC#109. Existing feature parity delivery #101 and
production OAuth acceptance #100 retain their separate acceptance criteria.

## Why this boundary

Current clients repeat calendar bounds, date conversion and schedule projections.
The iOS logic tests are linked to the application target even when they have no
UI dependency. Sharing an executable domain module allows these algorithms to
be tested directly on Linux and reused by the Go clients.

This does not eliminate native testing. SwiftUI layout, Observation/MainActor,
WebKit, Keychain, ASWebAuthenticationSession, SwiftData, URLSession behavior and
Widgets require their platform integration tests. The existing 348 iOS unit
cases include 21 cases directly importing Apple presentation or persistence
frameworks; other cases still have transitive platform dependencies.

## Product and module boundaries

| Owner | Responsibilities |
| --- | --- |
| Server | Canonical records, scopes, authorization, persisted mutations and upstream integration |
| Go core | Shared client calculations, projections and subsequently migrated client use cases |
| iOS adapter | Swift data conversion, observable state, Apple networking/storage/authentication and lifecycle |
| SwiftUI | Screens, navigation, accessibility, localized formatting and user interaction |
| Bot/CLI adapters | Commands, chat/terminal formatting, transport and client-specific account storage |

The first package is `calendar`. Bot and CLI import its typed Go API. The
`mobile` package exposes only supported scalar values and concrete return types
for the Objective-C binding used by Swift. It contains conversion, not a second
implementation of domain calculations. iOS links the library into its process;
there is no subprocess or local HTTP service.

The `mobile` build module owns gomobile dependencies and points to the core
source in the same checkout. Keeping those dependencies outside the root module
prevents a calendar import from upgrading Bot/CLI networking and tooling modules.
Only the root module and the compiled Apple artifact are consumer dependencies.

## First slice: calendar bounds and request windows

Inputs are explicit timestamps or ISO civil dates, a day/week/month period and
an explicit first weekday. Outputs are inclusive `from`/`to` civil dates and
contiguous request windows with at most 366 days each. Asia/Shanghai timezone
data is bundled as a fallback for environments without host zoneinfo files.

Preserve these existing client conventions during extraction:

| Consumer | First weekday |
| --- | --- |
| iOS personal calendar | Monday |
| iOS public Young calendar | Sunday |
| Bot teaching timetable | Sunday |
| Bot/CLI Young date filters | Monday |

These are explicit caller choices, not a shared default. Any product alignment
of those conventions is a separate behavior change. CLI date-only UTC filters
for todos retain their existing contract and are outside this slice.

The core receives an explicit clock input and performs no I/O. It introduces no
credentials, mutable application store, goroutines or callback lifecycle. Swift
retains account ownership, cancellation and stale-result checks around requests.

## Delivery sequence

1. **Calendar core and bridge.** Deliver typed Go APIs, common behavioral fixtures,
   Linux tests, an iOS device/simulator XCFramework and Swift binding tests. Record
   framework size and clean/incremental build time using the actual artifacts.
2. **Three live consumers.** Replace the corresponding existing Bot/CLI functions
   and iOS calendar calculations with the versioned module. Remove their migrated
   duplicate implementations in the same changes. Preserve command output,
   generated query bounds and the existing native calendar journeys.
3. **Schedule and activity projections.** Migrate meeting dates, teaching-week
   interpretation, overlap/conflict calculations and calendar grouping where
   two clients need the same rule. Resolve existing input-format differences
   against server contracts before changing them.
4. **Client use cases.** Move shared request construction, response projection,
   pagination and mutation planning once their actual consumers and transport
   boundaries are established. Keep platform credentials/storage/network hooks
   at the edge and expose feature-specific APIs instead of a generic dispatcher.

Each phase lands as a working product. A later phase starts from measured
results and concrete duplicated behavior; this roadmap does not authorize
unrelated features or additional robustness machinery.

The next proposed client-use-case slice is specified in the
[incremental pagination PRD](incremental-pagination.md).

## Acceptance

- Shared fixtures independently specify Shanghai day rollover, Sunday/Monday
  weeks, leap February, month/year boundaries, invalid dates, and 366-day windows.
- Those fixtures pass in Go and through the real compiled Swift binding. Tests
  use the same expected data, not one implementation as the other's oracle.
- Linux runs uncached Go tests with the race detector and `go vet`.
- macOS builds supported iPhone and simulator architectures and executes bridge
  tests. Existing SwiftUI/system integration tests remain required.
- Bot and CLI existing command/HTTP contract tests pass using the shared module.
- Each consumer pins a released module/framework version; the binary is built
  from the same commit as the Go source. No permanent local `replace` directive
  or checked-in generated framework is the released integration.
- Record signed app size/build-time impact after actual iOS integration. A
  successfully built but unused XCFramework is not completion of this task.

## Cost accepted for the shared implementation

The Go runtime and generated Objective-C interface add artifact size, a second
compiler/toolchain, conversion code and language-boundary tests. Swift structs,
optional values and async operations do not transparently become Go APIs.
The first synchronous scalar/value boundary keeps that cost observable before
moving stateful work. No speculative race or security machinery is included.

Initial measurements on arm64 macOS with Go 1.27.0 and Xcode 27.0:

| Measurement | Observed result |
| --- | --- |
| Empty Go compiler cache, warm module/tool downloads | 24.434 seconds |
| Same-source incremental framework build | 8.668 seconds |
| Device, simulator and macOS framework ZIP | 12,251,056 bytes |

These are individual local build observations. Network download failures are
recorded separately. Signed app size and app build time will be measured after
the actual iOS integration. The archive contains iOS arm64, simulator
arm64/x86_64 and macOS arm64/x86_64 slices.
