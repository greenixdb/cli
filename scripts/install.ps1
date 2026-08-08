<#
.SYNOPSIS
    Greenix CLI installer for Windows.
    Publisher: Edgicode Limited

.DESCRIPTION
    Downloads the correct Greenix CLI build for this machine, installs it to
    %LOCALAPPDATA%\Programs\Greenix and adds that folder to the user PATH.

.EXAMPLE
    irm https://raw.githubusercontent.com/greenixdb/cli/main/scripts/install.ps1 | iex
#>

[CmdletBinding()]
param(
    [string]$Version = $env:GREENIX_VERSION,
    [string]$InstallDir = $env:GREENIX_INSTALL
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

$Repo = "greenixdb/cli"
$BinName = "greenix.exe"

if (-not $Version) { $Version = "latest" }
if (-not $InstallDir) { $InstallDir = Join-Path $env:LOCALAPPDATA "Programs\Greenix" }

function Write-Info { param($m) Write-Host "  $m" }
function Fail { param($m) Write-Host "error: $m" -ForegroundColor Red; exit 1 }

# ------------------------------------------------------------- detect arch
$arch = switch ($env:PROCESSOR_ARCHITECTURE) {
    "AMD64" { "x64" }
    "ARM64" { "arm64" }
    "x86"   { if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") { "x64" } else { "x86" } }
    default { Fail "unsupported architecture: $env:PROCESSOR_ARCHITECTURE" }
}

$asset = "greenix-windows-$arch.zip"
$base = if ($Version -eq "latest") {
    "https://github.com/$Repo/releases/latest/download"
} else {
    "https://github.com/$Repo/releases/download/$Version"
}

Write-Host ""
Write-Host "Installing Greenix CLI" -ForegroundColor Green
Write-Info "platform : windows/$arch"
Write-Info "release  : $Version"
Write-Info "target   : $InstallDir\$BinName"
Write-Host ""

$tmp = Join-Path $env:TEMP ("greenix-" + [Guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Force -Path $tmp | Out-Null

try {
    $zipPath = Join-Path $tmp $asset
    Write-Info "downloading $asset..."
    try {
        Invoke-WebRequest -Uri "$base/$asset" -OutFile $zipPath -UseBasicParsing
    } catch {
        Fail "download failed: $base/$asset"
    }

    # ------------------------------------------------------ verify checksum
    try {
        $sumsPath = Join-Path $tmp "SHA256SUMS"
        Invoke-WebRequest -Uri "$base/SHA256SUMS" -OutFile $sumsPath -UseBasicParsing
        $line = Select-String -Path $sumsPath -Pattern ([Regex]::Escape($asset)) | Select-Object -First 1
        if ($line) {
            $expected = ($line.Line -split '\s+')[0]
            $actual = (Get-FileHash -Path $zipPath -Algorithm SHA256).Hash.ToLower()
            if ($actual -ne $expected.ToLower()) { Fail "checksum mismatch for $asset" }
            Write-Info "checksum verified"
        }
    } catch {
        Write-Info "warning: SHA256SUMS not available, skipping verification"
    }

    # Remove the mark-of-the-web so SmartScreen does not block extraction.
    Unblock-File -Path $zipPath -ErrorAction SilentlyContinue

    Expand-Archive -Path $zipPath -DestinationPath $tmp -Force
    $extracted = Join-Path $tmp $BinName
    if (-not (Test-Path $extracted)) { Fail "archive did not contain $BinName" }

    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    $target = Join-Path $InstallDir $BinName

    # A running greenix.exe cannot be overwritten; retire the old copy first.
    if (Test-Path $target) {
        try { Remove-Item $target -Force } catch {
            Move-Item $target "$target.old" -Force -ErrorAction SilentlyContinue
        }
    }
    Move-Item $extracted $target -Force
    Unblock-File -Path $target -ErrorAction SilentlyContinue

    # --------------------------------------------------------- update PATH
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if (-not $userPath) { $userPath = "" }
    $entries = $userPath -split ';' | Where-Object { $_ -ne "" }
    if ($entries -notcontains $InstallDir) {
        $newPath = (($entries + $InstallDir) -join ';')
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        $pathUpdated = $true
    }
    $env:Path = "$env:Path;$InstallDir"

    Write-Host ""
    Write-Host "Greenix CLI installed to $target" -ForegroundColor Green
    Write-Host ""
    if ($pathUpdated) {
        Write-Info "Added $InstallDir to your user PATH."
        Write-Info "Open a NEW terminal window, then run: greenix --version"
    } else {
        Write-Info "Run: greenix --version"
    }
    Write-Host ""
} finally {
    Remove-Item -Recurse -Force $tmp -ErrorAction SilentlyContinue
}
