<div align="center">

# Greenix CLI

**Build, deploy and manage Greenix Studio projects from your terminal.**

Published by **Edgicode Limited**

[![Build CLI Tool](https://github.com/greenixdb/cli/actions/workflows/build-cli.yml/badge.svg)](https://github.com/greenixdb/cli/actions/workflows/build-cli.yml)
[![Latest release](https://img.shields.io/github/v/release/greenixdb/cli?sort=semver)](https://github.com/greenixdb/cli/releases/latest)

</div>

---

## Contents

- [Install](#install)
- [Quick start](#quick-start)
- [Commands](#commands)
- [Supported platforms](#supported-platforms)
- [Manual download](#manual-download)
- [Code signing and security warnings](#code-signing-and-security-warnings)
- [Verifying a download](#verifying-a-download)
- [Uninstall](#uninstall)
- [Repository layout](#repository-layout)
- [Building from source](#building-from-source)
- [Release process](#release-process)
- [Support](#support)
- [License](#license)

---

## Install

The install scripts pick the right build for your machine, verify its SHA-256
checksum, put the binary on your `PATH` as `greenix`, and print what to do next.

### macOS and Linux

```bash
curl -fsSL https://raw.githubusercontent.com/greenixdb/cli/main/scripts/install.sh | sh
```

Installs to `/usr/local/bin` when writable, otherwise `~/.local/bin`.
Override with environment variables:

```bash
GREENIX_INSTALL="$HOME/bin" GREENIX_VERSION="v0.1.0" \
  sh -c "$(curl -fsSL https://raw.githubusercontent.com/greenixdb/cli/main/scripts/install.sh)"
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/greenixdb/cli/main/scripts/install.ps1 | iex
```

Installs to `%LOCALAPPDATA%\Programs\Greenix` and adds it to your user `PATH`.
**Open a new terminal window** afterwards so the updated `PATH` takes effect.

### Verify the install

```bash
greenix --version
```

---

## Quick start

```bash
greenix login              # authenticate against Greenix Studio
greenix whoami             # show the signed-in account
greenix init my-project    # scaffold a new project
greenix build --target all # build for every platform
greenix logout             # clear stored credentials
```

---

## Commands

| Command | Description |
| --- | --- |
| `greenix init <name>` | Initialize a new Greenix Studio project |
| `greenix build --target <t>` | Build for one or more platforms |
| `greenix login` | Authenticate and store credentials locally |
| `greenix logout` | Remove stored credentials |
| `greenix whoami` | Show the currently authenticated account |
| `greenix version` | Print version, build code, build date and platform |
| `greenix help` | Show help for any command |

Global flags: `--version`, `--help`.

Build targets for `--target`: `android`, `ios`, `windows`, `macos`, `linux`, `all`.

---

## Supported platforms

| Platform | x86 (32-bit) | x64 | arm64 | arm (v7) | universal |
| --- | :---: | :---: | :---: | :---: | :---: |
| Windows | yes | yes | yes | – | – |
| macOS | – | yes | yes | – | yes |
| Linux | yes | yes | yes | yes | – |

All builds are statically linked with `CGO_ENABLED=0`, so there are no runtime
dependencies — no Go toolchain, no libc version requirements.

The macOS **universal** build contains both Intel and Apple Silicon slices and
is what the install script and Homebrew formula use.

---

## Manual download

Every release publishes two shapes of each build on the
[releases page](https://github.com/greenixdb/cli/releases/latest):

| Asset | Use it for |
| --- | --- |
| `greenix-<platform>-<arch>.tar.gz` / `.zip` | Recommended. Contains a binary already named `greenix`, with the executable bit preserved. |
| `greenix-cli-v<version>-<code>-<platform>-<arch>` | Raw, unpacked binary kept for compatibility with existing links. |

The same tree is also committed to [`build-output/`](build-output) in this
repository for direct linking.

> **If you download manually you must do three things yourself**, which the
> install scripts otherwise handle:
>
> 1. **Rename it.** The file is `greenix-cli-v0.1.0-1`; the command is `greenix`.
>    Typing `greenix` will not find a differently named file.
> 2. **Make it executable** (macOS/Linux): `chmod +x greenix`.
> 3. **Put it on your `PATH`**, otherwise you must run it as `./greenix` from the
>    folder it sits in.

<details>
<summary>Manual install, step by step</summary>

**Linux / macOS**

```bash
tar -xzf greenix-linux-x64.tar.gz
chmod +x greenix
sudo mv greenix /usr/local/bin/
greenix --version
```

**Windows (PowerShell)**

```powershell
Expand-Archive greenix-windows-x64.zip -DestinationPath "$env:LOCALAPPDATA\Programs\Greenix"
[Environment]::SetEnvironmentVariable(
  "Path",
  [Environment]::GetEnvironmentVariable("Path","User") + ";$env:LOCALAPPDATA\Programs\Greenix",
  "User")
# open a NEW terminal, then:
greenix --version
```

</details>

---

## Code signing and security warnings

Release binaries are signed automatically **when signing credentials are
configured in the repository secrets**. Until then, operating systems will warn
about them. What you will see, and how to get past it:

### Windows

An unsigned `.exe` triggers Microsoft Defender SmartScreen:
*"Windows protected your PC – unknown publisher."* Choose **More info → Run
anyway**. Browsers may also flag the download itself.

Once an Authenticode certificate is configured, the workflow signs and
timestamps every `.exe` and the publisher shows as **Edgicode Limited**. Note
that a standard **OV** certificate still shows SmartScreen warnings until the
signature builds download reputation; an **EV** certificate is trusted
immediately.

### macOS

macOS is stricter. Anything downloaded by a browser gets a
`com.apple.quarantine` attribute, and Gatekeeper then refuses to run it —
often with the misleading message *"greenix is damaged and can't be opened."*
That message does not mean the file is corrupt; it means it is not notarized.

Work around it for a manual download:

```bash
xattr -d com.apple.quarantine ./greenix
```

…or open **System Settings → Privacy & Security → Open Anyway**.

**You will not hit this at all** if you install with the script above, with
Homebrew, or with `go install` — none of those set the quarantine attribute.
Once an Apple Developer ID certificate is configured, builds are signed with
the hardened runtime, notarized by Apple and the warning disappears entirely.

### Linux

No signing gates. Verify the SHA-256 checksum instead.

---

## Verifying a download

Each release publishes a `SHA256SUMS` file covering every asset.

```bash
curl -fsSLO https://github.com/greenixdb/cli/releases/latest/download/SHA256SUMS
sha256sum -c SHA256SUMS --ignore-missing
```

```powershell
(Get-FileHash .\greenix-windows-x64.zip -Algorithm SHA256).Hash
# compare against the matching line in SHA256SUMS
```

The install scripts perform this check automatically.

---

## Uninstall

```bash
# macOS / Linux
sudo rm -f /usr/local/bin/greenix    # or ~/.local/bin/greenix
rm -rf ~/.greenix                    # stored credentials and config
```

```powershell
# Windows
Remove-Item "$env:LOCALAPPDATA\Programs\Greenix" -Recurse -Force
# then remove that folder from your user PATH in System Environment Variables
```

---

## Repository layout

```text
.
├── .github/workflows/build-cli.yml   CI: build, sign, checksum, release
├── build-output/                     Binaries committed by CI, per platform/arch
│   ├── linux/{x86,x64,arm,arm64}/
│   ├── macos/{x64,arm64,universal}/
│   └── windows/{x86,x64,arm64}/
├── cli/
│   ├── src/                          Go source (main.go, cmd/, internal/, pkg/)
│   ├── scripts/                      Local build scripts, mirroring CI
│   └── tests/
├── packaging/
│   ├── homebrew/greenix.rb           Homebrew formula template
│   └── scoop/greenix.json            Scoop manifest template
├── scripts/
│   ├── install.sh                    macOS + Linux installer
│   └── install.ps1                   Windows installer
└── version.properties                Single source of truth for the version
```

---

## Building from source

Requires Go 1.21 or newer.

```bash
git clone https://github.com/greenixdb/cli.git
cd cli/cli/src
go build -o greenix .
./greenix --version
```

To reproduce the full CI matrix locally into `build-output/`:

```bash
bash cli/scripts/build-all.sh      # all platforms
bash cli/scripts/build-linux.sh    # linux only
bash cli/scripts/build-macos.sh    # macos only
pwsh -File cli/scripts/build-windows.ps1
```

Run the tests:

```bash
cd cli/src && go test ./...
```

---

## Release process

1. Bump `VERSION_NAME` / `VERSION_CODE` in [`version.properties`](version.properties).
2. Commit, then tag and push:

   ```bash
   git tag v0.1.1 && git push origin v0.1.1
   ```

3. The `Build CLI Tool` workflow then:
   - compiles all nine platform/arch targets,
   - creates the macOS universal binary with `lipo`,
   - signs and notarizes (when signing secrets are present),
   - generates archives plus `SHA256SUMS`,
   - commits the refreshed `build-output/` tree back to the default branch,
   - publishes a GitHub Release with every asset attached.

`workflow_dispatch` runs everything except the release step.

---

## Support

- Documentation: <https://docs.greenix.studio>
- Issues: <https://github.com/greenixdb/cli/issues>

---

## License

Proprietary. Copyright © Edgicode Limited. All rights reserved.
