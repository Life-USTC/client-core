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
cp calendar/testdata/teaching_week_cases.json apple/Tests/LifeClientCoreTests/Fixtures/teaching_week_cases.json
cp young/testdata/query_cases.json apple/Tests/LifeClientCoreTests/Fixtures/young_query_cases.json
cp pagination/testdata/progress_cases.json apple/Tests/LifeClientCoreTests/Fixtures/page_progress_cases.json
cp ustc/testdata/course_table_cases.json apple/Tests/LifeClientCoreTests/Fixtures/ustc_course_table_cases.json
cp ustc/testdata/academic_cases.json apple/Tests/LifeClientCoreTests/Fixtures/ustc_academic_cases.json
cp ustc/testdata/blackboard_cases.json apple/Tests/LifeClientCoreTests/Fixtures/ustc_blackboard_cases.json
