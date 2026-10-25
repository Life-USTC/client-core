// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "LifeClientCore",
    platforms: [.iOS(.v17), .macOS(.v13)],
    products: [.library(name: "LifeClientCore", targets: ["LifeClientCore"])],
    targets: [
        .binaryTarget(
            name: "LifeCore",
            url: "https://github.com/Life-USTC/client-core/releases/download/v0.4.0/LifeCore.xcframework.zip",
            checksum: "49c463a081850b35f294625a9be423f9e374dda68685a34a1a90ebdf9bde56c9"
        ),
        .target(
            name: "LifeClientCore",
            dependencies: ["LifeCore"],
            path: "apple/Sources/LifeClientCore"
        ),
    ]
)
