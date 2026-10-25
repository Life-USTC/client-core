// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "LifeClientCore",
    platforms: [.iOS(.v17), .macOS(.v13)],
    products: [.library(name: "LifeClientCore", targets: ["LifeClientCore"])],
    targets: [
        .binaryTarget(
            name: "LifeCore",
            url: "https://github.com/Life-USTC/client-core/releases/download/v0.5.0/LifeCore.xcframework.zip",
            checksum: "e82c74ea8c9463e8141bd605829f04a02cf1d65efaa22a5028a518103abbad8a"
        ),
        .target(
            name: "LifeClientCore",
            dependencies: ["LifeCore"],
            path: "apple/Sources/LifeClientCore"
        ),
    ]
)
