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

# Build architectures
for arch in x86 x64 arm64; do
    goarch=""
    case $arch in
        x86) goarch="386" ;;
        x64) goarch="amd64" ;;
        arm64) goarch="arm64" ;;
    esac
    
    mkdir -p ../build-output/cli/linux/$arch
    echo "🔨 Building $arch..."
    
    GOOS=linux GOARCH=$goarch go build -ldflags "-X main.Version=$VERSION_NAME -X main.BuildCode=$VERSION_CODE -X main.BuildDate='$BUILD_DATE'" \
        -o ../build-output/cli/linux/$arch/greenix-$VERSION_NAME-$VERSION_CODE ../src
    
    if [ $? -eq 0 ]; then
        echo "✅ Built $arch"
        ln -sf greenix-$VERSION_NAME-$VERSION_CODE ../build-output/cli/linux/$arch/greenix
        chmod +x ../build-output/cli/linux/$arch/greenix-$VERSION_NAME-$VERSION_CODE
    else
        echo "❌ Failed to build $arch"
    fi
    echo ""
done

echo "✅ Build complete!"

