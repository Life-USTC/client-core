# Life@USTC client core

Shared Go domain logic for the iOS app, Bot and CLI. The server remains the
source of business permissions and persisted state; clients share calculations
and use cases through this module.

The first slice provides inclusive calendar bounds, explicit week conventions,
Shanghai civil dates and bounded request windows. Bot/CLI import `calendar`;
iOS consumes `mobile` through a generated XCFramework.

Teaching-week bounds and week numbers use the caller's explicit first weekday.
Week one contains the semester's first day; week bounds remain complete even
when part of the week falls outside the semester. Adapters retain explicit
meeting-date priority and semester inclusion rules.

The `young` package projects public event filters into request parameters and
switches a query to unknown dates while retaining its filter context. Optional
booleans preserve false versus absence. Callers choose pagination, validate user
input and aggregate responses.

The `pagination` package records accepted row counts and the next request page.
`Progress` advances using the caller's has-more decision; validation, fetching
and collection completeness remain in each client. Its existing consumers
include the CLI's shared reader across catalog, workspace and school endpoints
and the Bot's Young activity reader.

The `ustc` package shares undergraduate course-table, catalog, score, exam,
Blackboard, graduate YJS1 and static/app curriculum request and response rules. Explicit parser profiles preserve each client’s
existing formats. Clients retain authentication, transport, persistence, native
date conversion and empty-table policy. See the
[USTC extraction PRD and roadmap](docs/ustc.md) for its two-client acceptance
and the remaining university capabilities.

See [the PRD and roadmap](docs/architecture.md). Initial shared-core delivery is tracked
in [Task #109](https://github.com/Life-USTC/Life-USTC/issues/109).

Go clients import the `calendar`, `young`, `pagination` and `ustc` packages from
`github.com/Life-USTC/client-core` at `v0.7.0`.
Apple clients add this repository as a Swift package and use the `LifeClientCore`
product. Pin the full distribution commit revision containing the updated root
`Package.swift`; release tags identify the Go and framework source, while the
following distribution commit selects that archive and checksum. Swift
applications do not run the Go build tools.

The release includes `source-commit.txt` and `LifeCore.xcframework.sha256` with
the verified build archive. The Swift adapter source matches the release source;
the later package manifest supplies distribution metadata.

```sh
go test -race -count=1 ./...
go vet ./...
(cd mobile && go test -race -count=1 ./... && go vet ./...)
```

Apple framework generation and bridge verification require macOS and Xcode.
Generated binaries remain build artifacts and are not committed to this repo.

```sh
scripts/build-apple.sh
cd apple
swift test
xcodebuild test -scheme LifeClientCore \
  -destination 'platform=iOS Simulator,name=iPhone 17'
```

The build uses the gomobile tools pinned by `mobile/go.mod`. The framework contains
iOS device, iOS simulator and macOS slices; macOS tests exercise the same Swift
adapter and shared fixtures without booting a simulator. CI publishes the
framework archive, its checksum and the exact source commit as one artifact.

`calendar` has no third-party module dependencies. The `mobile` build module
isolates gomobile and its tools so they cannot upgrade Bot/CLI dependencies.
Its local core replacement binds the checked-out source; applications consume
the released Go module or framework, not this build module.
