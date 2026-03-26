param(
  [string]$WorkspaceRoot = ""
)

$ErrorActionPreference = 'Stop'

if ([string]::IsNullOrWhiteSpace($WorkspaceRoot)) {
  $WorkspaceRoot = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
}

$releaseDir = Join-Path $WorkspaceRoot 'release'
$binDir = Join-Path $WorkspaceRoot 'build\bin'
New-Item -ItemType Directory -Force -Path $releaseDir | Out-Null

if (-not (Test-Path $binDir)) {
  throw "build/bin was not found. Please build the application first."
}

$timestamp = Get-Date -Format 'yyyyMMdd-HHmmss'
$zipPath = Join-Path $releaseDir ("AIGuard-windows-$timestamp.zip")
if (Test-Path $zipPath) { Remove-Item $zipPath -Force }
Compress-Archive -Path (Join-Path $binDir '*') -DestinationPath $zipPath
Write-Host "Archived to $zipPath" -ForegroundColor Green
