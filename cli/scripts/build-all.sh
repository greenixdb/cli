#!/bin/bash
# Builds every platform/arch the CI workflow builds, into build-output/ at repo root.
set -u
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🚀 Building Greenix CLI for all platforms..."
echo ""

echo "📦 Windows..."
if command -v pwsh >/dev/null 2>&1; then
    pwsh -File "$SCRIPT_DIR/build-windows.ps1"
elif command -v powershell.exe >/dev/null 2>&1; then
    powershell.exe -File "$SCRIPT_DIR/build-windows.ps1"
else
    echo "⚠️ PowerShell not found, skipping Windows builds"
fi

echo ""
echo "📦 macOS..."
bash "$SCRIPT_DIR/build-macos.sh"

echo ""
echo "📦 Linux..."
bash "$SCRIPT_DIR/build-linux.sh"

echo ""
echo "✅ All builds complete!"
