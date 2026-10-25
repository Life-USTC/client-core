// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "LifeClientCore",
    platforms: [.iOS(.v17), .macOS(.v13)],
    products: [.library(name: "LifeClientCore", targets: ["LifeClientCore"])],
    targets: [
        .binaryTarget(
            name: "LifeCore",
            url: "https://github.com/Life-USTC/client-core/releases/download/v0.6.0/LifeCore.xcframework.zip",
            checksum: "f11e0746471f9d4b3619847ea0a389bc4886c23199bd5199b0eb29f358a6d9ce"
        ),
        .target(
            name: "LifeClientCore",
            dependencies: ["LifeCore"],
            path: "apple/Sources/LifeClientCore"
        ),
    ]
)
