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

See [the PRD and roadmap](docs/architecture.md). Cross-client delivery is tracked
in [Task #109](https://github.com/Life-USTC/Life-USTC/issues/109).

Go clients use `github.com/Life-USTC/client-core/calendar` at `v0.2.0`.
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
