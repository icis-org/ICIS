# update-icons.ps1 — Regenerate all icon assets from build/appicon.png
#
# Usage:
#   .\scripts\update-icons.ps1                     # uses existing build/appicon.png
#   .\scripts\update-icons.ps1 path\to\new.png     # replaces source first
#
# Requires: Go (for the png2ico converter)

param(
    [string]$Source
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent (Split-Path -Parent $PSCommandPath)
$BuildDir = Join-Path $Root "build"
$WinDir = Join-Path $BuildDir "windows"
$ToolDir = Join-Path $env:TEMP "opencode"
$ToolSrc = Join-Path $ToolDir "png2ico.go"
$ToolExe = Join-Path $ToolDir "png2ico.exe"

# Step 0: replace source if provided
if ($Source) {
    if (!(Test-Path $Source)) {
        Write-Error "Source PNG not found: $Source"
        exit 1
    }
    Copy-Item $Source (Join-Path $BuildDir "appicon.png") -Force
    Write-Host "Copied $Source -> build/appicon.png" -ForegroundColor Cyan
}

$Master = Join-Path $BuildDir "appicon.png"
if (!(Test-Path $Master)) {
    Write-Error "Master icon not found: $Master"
    exit 1
}

# Step 1: build the converter if needed
if (!(Test-Path $ToolExe)) {
    Write-Host "Building png2ico tool..." -ForegroundColor Yellow
    Push-Location $ToolDir
    try {
        & go build -o png2ico.exe png2ico.go
        if ($LASTEXITCODE -ne 0) { throw "go build failed" }
    } finally {
        Pop-Location
    }
}

# Step 2: generate ICO files
Write-Host "Generating icon.ico..." -ForegroundColor Green
& $ToolExe $Master (Join-Path $WinDir "icon.ico")
if ($LASTEXITCODE -ne 0) { exit 1 }

Write-Host "Generating appicon.ico..." -ForegroundColor Green
& $ToolExe $Master (Join-Path $WinDir "appicon.ico")
if ($LASTEXITCODE -ne 0) { exit 1 }

# Step 3: copy to frontend
$FavDir = Join-Path $Root "frontend\public"
$ImgDir = Join-Path $Root "frontend\src\assets\images"
New-Item -ItemType Directory -Path $FavDir -Force | Out-Null
New-Item -ItemType Directory -Path $ImgDir -Force | Out-Null

Copy-Item $Master (Join-Path $FavDir "favicon.png") -Force
Copy-Item $Master (Join-Path $ImgDir "logo-universal.png") -Force

Write-Host ""
Write-Host "All icons updated from $Master" -ForegroundColor Cyan
Write-Host "  build/windows/icon.ico       (exe + shortcuts + ARP)"
Write-Host "  build/windows/appicon.ico    (.ici file association)"
Write-Host "  frontend/public/favicon.png  (browser tab)"
Write-Host "  frontend/src/assets/images/logo-universal.png"
Write-Host ""
Write-Host "Run 'wails build -nsis -installscope user' to rebuild." -ForegroundColor Yellow
