// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "LifeClientCore",
    platforms: [.iOS(.v17), .macOS(.v13)],
    products: [.library(name: "LifeClientCore", targets: ["LifeClientCore"])],
    targets: [
        .binaryTarget(
            name: "LifeCore",
            url: "https://github.com/Life-USTC/client-core/releases/download/v0.2.0/LifeCore.xcframework.zip",
            checksum: "decb26e729488a7a10b389f05a70b624e5710a120c934e679c64eae31cd0fc42"
        ),
        .target(
            name: "LifeClientCore",
            dependencies: ["LifeCore"],
            path: "apple/Sources/LifeClientCore"
        ),
    ]
)
