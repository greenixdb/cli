# Greenix Studio CLI

A command-line tool for building, deploying, and managing Greenix Studio projects across multiple platforms.

## Installation

### Download Binary

Download the appropriate binary for your platform from the [releases page](https://github.com/greenix-studio/cli/releases):

- **Windows**: `greenix.exe`
- **macOS**: `greenix`
- **Linux**: `greenix`

### Build from Source

```bash
git clone https://github.com/greenix-studio/cli.git
cd cli/src
go build -o greenix
```

Quick Start

```bash
# Initialize a new project
greenix init my-project

# Build for Android
greenix build --target android

# Build for all platforms
greenix build --target all

# Check version
greenix version
```

Commands

Command Description
init Initialize a new Greenix Studio project
build Build for one or more platforms
version Print version information
help Display help information

Build Targets

· android: Android APK/AAB
· ios: iOS IPA
· windows: Windows EXE/MSI
· macos: macOS APP/DMG
· linux: Linux ELF/DEB/RPM
· all: All platforms

Architecture Support

Platform x86 x64 arm64
Windows ✅ ✅ ✅
macOS ❌ ✅ ✅
Linux ✅ ✅ ✅

Development

```bash
# Run tests
go test ./...

# Build locally
go build -o greenix

# Run
./greenix version
```

License

PROPRIETARY License - see LICENSE for details.

```

---

## 5. Git Configuration

### `cli/.gitignore`
```

Binaries

*.exe
*.exe~
*.dll
*.so
.dylib
greenix
greenix-

Test binary

*.test

Output of the go coverage tool

*.out

Dependency directories

vendor/

Go workspace file

go.work

IDE

.idea/
.vscode/
*.swp
*.swo

Build output

build-output/
dist/

OS files

.DS_Store
Thumbs.db

```

---

## 6. Usage Examples

### Build the CLI locally
```bash
cd cli/src
go build -o greenix
./greenix version
```

Output

```
██████╗ ██████╗ ███████╗███████╗███╗   ██╗██╗██╗  ██╗
██╔════╝ ██╔══██╗██╔════╝██╔════╝████╗  ██║██║╚██╗██╔╝
██║  ███╗██████╔╝█████╗  █████╗  ██╔██╗ ██║██║ ╚███╔╝ 
██║   ██║██╔══██╗██╔══╝  ██╔══╝  ██║╚██╗██║██║ ██╔██╗ 
╚██████╔╝██║  ██║███████╗███████╗██║ ╚████║██║██╔╝ ██╗
 ╚═════╝ ╚═╝  ╚═╝╚══════╝╚══════╝╚═╝  ╚═══╝╚═╝╚═╝  ╚═╝

📦 Greenix Studio CLI v1.2.3 (Build 123)
   Build Date: 2026-08-08 14:30:00
   Go Version: go1.21.5
   Platform:   darwin/arm64

📖 Documentation: https://docs.greenix.studio
🐛 Report issues: https://github.com/greenix-studio/cli/issues
```
