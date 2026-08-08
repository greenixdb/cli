# Read version from root
$content = Get-Content -Path "../../version.properties" -ErrorAction SilentlyContinue
if (-not $content) {
    Write-Host "⚠️ version.properties not found, using defaults"
    $versionName = "dev"
    $versionCode = "0"
} else {
    $versionName = ($content | Select-String 'VERSION_NAME=(.*)').Matches.Groups[1].Value
    $versionCode = ($content | Select-String 'VERSION_CODE=(.*)').Matches.Groups[1].Value
}

$buildDate = Get-Date -Format "yyyy-MM-dd HH:mm:ss"

Write-Host "📦 Building Greenix CLI v$versionName (Build $versionCode)" -ForegroundColor Cyan
Write-Host ""

$archs = @("x86", "x64", "arm64")

foreach ($arch in $archs) {
    $goarch = switch ($arch) {
        "x86" { "386" }
        "x64" { "amd64" }
        "arm64" { "arm64" }
    }
    
    $outputDir = "../build-output/cli/windows/$arch"
    New-Item -ItemType Directory -Force -Path $outputDir | Out-Null
    
    $filename = "greenix-$versionName-$versionCode.exe"
    $output = "$outputDir/$filename"
    
    Write-Host "🔨 Building $arch..." -ForegroundColor Yellow
    
    $env:GOOS = "windows"
    $env:GOARCH = $goarch
    
    go build -ldflags "-X main.Version=$versionName -X main.BuildCode=$versionCode -X main.BuildDate='$buildDate'" `
             -o $output ../src
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Built: $filename" -ForegroundColor Green
        
        # Create symlink without version
        $link = "$outputDir/greenix.exe"
        Remove-Item -Path $link -Force -ErrorAction SilentlyContinue
        New-Item -ItemType SymbolicLink -Path $link -Target $filename | Out-Null
    } else {
        Write-Host "❌ Failed to build $arch" -ForegroundColor Red
    }
    Write-Host ""
}

Write-Host "✅ Build complete!" -ForegroundColor Green
