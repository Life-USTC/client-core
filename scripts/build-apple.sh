#!/usr/bin/env bash
set -euo pipefail

repo_dir="$(cd "$(dirname "$0")/.." && pwd)"
cd "$repo_dir"
mkdir -p .build/bin apple/Frameworks
export GOBIN="$repo_dir/.build/bin"
export PATH="$GOBIN:$PATH"
cd mobile
go install golang.org/x/mobile/cmd/gomobile golang.org/x/mobile/cmd/gobind

# macOS runs the same Swift bridge fixtures without starting a simulator.
# Device and both simulator architectures remain part of the delivered artifact.
gomobile bind -target=ios,iossimulator,macos -iosversion=17.0 -macosversion=13.0 -prefix LifeCore \
    -o "$repo_dir/apple/Frameworks/LifeCore.xcframework" .

cd "$repo_dir"
mkdir -p apple/Tests/LifeClientCoreTests/Fixtures
cp calendar/testdata/calendar_cases.json apple/Tests/LifeClientCoreTests/Fixtures/calendar_cases.json
