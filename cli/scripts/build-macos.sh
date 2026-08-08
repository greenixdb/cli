#!/bin/bash

# Read version from root
if [ -f "../../version.properties" ]; then
    source ../../version.properties
else
    VERSION_NAME="dev"
    VERSION_CODE="0"
fi

BUILD_DATE=$(date "+%Y-%m-%d %H:%M:%S")

echo "📦 Building Greenix CLI v$VERSION_NAME (Build $VERSION_CODE)"
echo ""

# Build x64 (Intel)
mkdir -p ../build-output/cli/macos/x64
echo "🔨 Building x64 (Intel)..."
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$VERSION_NAME -X main.BuildCode=$VERSION_CODE -X main.BuildDate='$BUILD_DATE'" \
    -o ../build-output/cli/macos/x64/greenix-$VERSION_NAME-$VERSION_CODE ../src

if [ $? -eq 0 ]; then
    echo "✅ Built x64"
    ln -sf greenix-$VERSION_NAME-$VERSION_CODE ../build-output/cli/macos/x64/greenix
else
    echo "❌ Failed to build x64"
fi

# Build arm64 (Apple Silicon)
mkdir -p ../build-output/cli/macos/arm64
echo "🔨 Building arm64 (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.Version=$VERSION_NAME -X main.BuildCode=$VERSION_CODE -X main.BuildDate='$BUILD_DATE'" \
    -o ../build-output/cli/macos/arm64/greenix-$VERSION_NAME-$VERSION_CODE ../src

if [ $? -eq 0 ]; then
    echo "✅ Built arm64"
    ln -sf greenix-$VERSION_NAME-$VERSION_CODE ../build-output/cli/macos/arm64/greenix
else
    echo "❌ Failed to build arm64"
fi

# Create Universal Binary
echo ""
echo "🔨 Creating Universal Binary..."
mkdir -p ../build-output/cli/macos/universal

lipo -create -output ../build-output/cli/macos/universal/greenix-$VERSION_NAME-$VERSION_CODE \
    ../build-output/cli/macos/x64/greenix-$VERSION_NAME-$VERSION_CODE \
    ../build-output/cli/macos/arm64/greenix-$VERSION_NAME-$VERSION_CODE

if [ $? -eq 0 ]; then
    echo "✅ Universal Binary created"
    ln -sf greenix-$VERSION_NAME-$VERSION_CODE ../build-output/cli/macos/universal/greenix
    chmod +x ../build-output/cli/macos/universal/greenix-$VERSION_NAME-$VERSION_CODE
else
    echo "❌ Failed to create Universal Binary"
fi

echo ""
echo "✅ Build complete!"

