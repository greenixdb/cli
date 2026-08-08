#!/bin/bash
# Builds macOS binaries into the same layout the CI workflow produces:
#   build-output/macos/<x64|arm64>/greenix-cli-v<VERSION_NAME>-<VERSION_CODE>
set -u

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"

if [ -f "$ROOT_DIR/version.properties" ]; then
    source "$ROOT_DIR/version.properties"
else
    VERSION_NAME="dev"
    VERSION_CODE="0"
fi

BUILD_DATE=$(date -u "+%Y-%m-%d %H:%M:%S UTC")
BIN="greenix-cli-v${VERSION_NAME}-${VERSION_CODE}"

echo "📦 Building Greenix CLI v$VERSION_NAME (Build $VERSION_CODE) - macos"
echo ""

for arch in x64 arm64; do
    case $arch in
        x64) goarch="amd64" ;;
        arm64) goarch="arm64" ;;
    esac

    OUT_DIR="$ROOT_DIR/build-output/macos/$arch"
    mkdir -p "$OUT_DIR"
    echo "🔨 Building $arch..."

    (cd "$ROOT_DIR/cli/src" && GOOS=darwin GOARCH=$goarch CGO_ENABLED=0 go build -trimpath \
        -ldflags "-s -w -X main.Version=$VERSION_NAME -X main.BuildCode=$VERSION_CODE -X 'main.BuildDate=$BUILD_DATE'" \
        -o "$OUT_DIR/$BIN" .)

    if [ $? -eq 0 ]; then
        chmod +x "$OUT_DIR/$BIN"
        echo "✅ Built $OUT_DIR/$BIN"
    else
        echo "❌ Failed to build $arch"
    fi
    echo ""
done
