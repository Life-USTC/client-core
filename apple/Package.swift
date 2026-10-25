// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "LifeClientCore",
    platforms: [.iOS(.v17), .macOS(.v13)],
    products: [.library(name: "LifeClientCore", targets: ["LifeClientCore"])],
    targets: [
        .binaryTarget(name: "LifeCore", path: "Frameworks/LifeCore.xcframework"),
        .target(name: "LifeClientCore", dependencies: ["LifeCore"]),
        .testTarget(
            name: "LifeClientCoreTests",
            dependencies: ["LifeClientCore", "LifeCore"],
            resources: [.copy("Fixtures")]
        ),
    ]
)
