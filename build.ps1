# Build script for WatchUp Agent (Windows)
# Builds binaries for all supported platforms

param(
    [string]$Version = "dev"
)

$ErrorActionPreference = "Stop"

$OUTPUT_DIR = "dist"

Write-Host "Building WatchUp Agent v$Version" -ForegroundColor Blue
Write-Host ""

# Create output directory
New-Item -ItemType Directory -Force -Path $OUTPUT_DIR | Out-Null

# Build matrix
$platforms = @(
    @{OS="linux"; ARCH="amd64"},
    @{OS="linux"; ARCH="arm64"},
    @{OS="linux"; ARCH="arm"; ARM="7"},
    @{OS="darwin"; ARCH="amd64"},
    @{OS="darwin"; ARCH="arm64"},
    @{OS="windows"; ARCH="amd64"},
    @{OS="windows"; ARCH="arm64"}
)

foreach ($platform in $platforms) {
    $GOOS = $platform.OS
    $GOARCH = $platform.ARCH
    $GOARM = $platform.ARM
    
    $outputName = "watchup-agent-$GOOS-$GOARCH"
    if ($GOARM) {
        $outputName = "watchup-agent-$GOOS-armv$GOARM"
    }
    
    if ($GOOS -eq "windows") {
        $outputName = "$outputName.exe"
    }
    
    $platformStr = "$GOOS/$GOARCH"
    if ($GOARM) {
        $platformStr = "$platformStr/v$GOARM"
    }
    
    Write-Host "Building for $platformStr..." -ForegroundColor Blue
    
    $env:GOOS = $GOOS
    $env:GOARCH = $GOARCH
    $env:GOARM = $GOARM
    $env:CGO_ENABLED = "0"
    
    go build -ldflags="-s -w -X main.Version=$Version" `
        -o "$OUTPUT_DIR\$outputName" `
        cmd/agent/main.go cmd/agent/setup.go
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Built $outputName" -ForegroundColor Green
    } else {
        Write-Host "✗ Failed to build $outputName" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "Build complete! Binaries are in $OUTPUT_DIR\" -ForegroundColor Green
Write-Host ""
Write-Host "To create a release:"
Write-Host "  git tag v$Version"
Write-Host "  git push origin v$Version"