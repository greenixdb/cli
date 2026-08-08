# Builds Windows binaries into the same layout the CI workflow produces:
#   build-output/windows/<x86|x64|arm64>/greenix-cli-v<VERSION_NAME>-<VERSION_CODE>.exe

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$rootDir = (Resolve-Path (Join-Path $scriptDir "../..")).Path
$versionFile = Join-Path $rootDir "version.properties"

$content = Get-Content -Path $versionFile -ErrorAction SilentlyContinue
if (-not $content) {
    Write-Host "⚠️ version.properties not found, using defaults"
    $versionName = "dev"
    $versionCode = "0"
} else {
    $versionName = ($content | Select-String 'VERSION_NAME=(.*)').Matches.Groups[1].Value
    $versionCode = ($content | Select-String 'VERSION_CODE=(.*)').Matches.Groups[1].Value
}

$buildDate = (Get-Date).ToUniversalTime().ToString("yyyy-MM-dd HH:mm:ss") + " UTC"
$binName = "greenix-cli-v$versionName-$versionCode.exe"

Write-Host "📦 Building Greenix CLI v$versionName (Build $versionCode) - windows" -ForegroundColor Cyan
Write-Host ""

$archs = @("x86", "x64", "arm64")

foreach ($arch in $archs) {
    $goarch = switch ($arch) {
        "x86" { "386" }
        "x64" { "amd64" }
        "arm64" { "arm64" }
    }

    $outputDir = Join-Path $rootDir "build-output/windows/$arch"
    New-Item -ItemType Directory -Force -Path $outputDir | Out-Null
    $output = Join-Path $outputDir $binName

    Write-Host "🔨 Building $arch..." -ForegroundColor Yellow

    $env:GOOS = "windows"
    $env:GOARCH = $goarch
    $env:CGO_ENABLED = "0"

    Push-Location (Join-Path $rootDir "cli/src")
    go build -trimpath -ldflags "-s -w -X main.Version=$versionName -X main.BuildCode=$versionCode -X `"main.BuildDate=$buildDate`"" -o $output .
    $code = $LASTEXITCODE
    Pop-Location

    if ($code -eq 0) {
        Write-Host "✅ Built: $output" -ForegroundColor Green
    } else {
        Write-Host "❌ Failed to build $arch" -ForegroundColor Red
    }
    Write-Host ""
}

Write-Host "✅ Build complete!" -ForegroundColor Green
