// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "LifeClientCore",
    platforms: [.iOS(.v17), .macOS(.v13)],
    products: [.library(name: "LifeClientCore", targets: ["LifeClientCore"])],
    targets: [
        .binaryTarget(
            name: "LifeCore",
            url: "https://github.com/Life-USTC/client-core/releases/download/v0.1.0/LifeCore.xcframework.zip",
            checksum: "ba9f64888d0c5ae197eb47c6d9d2fa489c72f3087c84bf89fbaef03c43ceab56"
        ),
        .target(
            name: "LifeClientCore",
            dependencies: ["LifeCore"],
            path: "apple/Sources/LifeClientCore"
        ),
    ]
)
