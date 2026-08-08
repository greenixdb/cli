#!/bin/bash

echo "🚀 Building Greenix CLI for all platforms..."
echo ""

# Build Windows
echo "📦 Windows..."
powershell.exe -File ./build-windows.ps1

# Build macOS
echo ""
echo "📦 macOS..."
./build-macos.sh

# Build Linux
echo ""
echo "📦 Linux..."
./build-linux.sh

echo ""
echo "✅ All builds complete!"

