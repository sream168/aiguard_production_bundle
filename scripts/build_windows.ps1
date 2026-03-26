param(
  [switch]$Nsis,
  [switch]$Sign,
  [string]$SignPfxPath = "",
  [string]$SignPfxPassword = "",
  [string]$TimestampUrl = "http://ts.ssl.com",
  [string]$WorkspaceRoot = ""
)

$ErrorActionPreference = 'Stop'

if ([string]::IsNullOrWhiteSpace($WorkspaceRoot)) {
  $WorkspaceRoot = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
}

Push-Location $WorkspaceRoot
try {
  ./scripts/preflight.ps1

  Write-Host "`n[1/5] Install frontend dependencies" -ForegroundColor Cyan
  Push-Location frontend
  npm install
  npm run build
  Pop-Location

  Write-Host "`n[2/5] Generate go.sum and tidy modules" -ForegroundColor Cyan
  go mod tidy

  Write-Host "`n[3/5] Run wails doctor" -ForegroundColor Cyan
  wails doctor

  Write-Host "`n[4/5] Build Windows EXE" -ForegroundColor Cyan
  wails build -platform windows/amd64

  if ($Nsis) {
    Write-Host "`n[4.5/5] Build NSIS installer" -ForegroundColor Cyan
    wails build -platform windows/amd64 -nsis
  }

  if ($Sign) {
    if ([string]::IsNullOrWhiteSpace($SignPfxPath) -or [string]::IsNullOrWhiteSpace($SignPfxPassword)) {
      throw "-Sign requires both -SignPfxPath and -SignPfxPassword"
    }

    $signtool = Get-ChildItem 'C:\Program Files (x86)\Windows Kits\10\bin' -Recurse -Filter signtool.exe |
      Sort-Object FullName -Descending |
      Select-Object -First 1

    if ($null -eq $signtool) {
      throw "signtool.exe was not found. Please install the Windows SDK."
    }

    Get-ChildItem ./build/bin -Filter *.exe | ForEach-Object {
      Write-Host "Sign $($_.FullName)" -ForegroundColor Yellow
      & $signtool.FullName sign /fd sha256 /tr $TimestampUrl /f $SignPfxPath /p $SignPfxPassword $_.FullName
    }
  }

  Write-Host "`n[5/5] Archive artifacts" -ForegroundColor Cyan
  ./scripts/package_release.ps1

  Write-Host "`nDone. Artifacts are available in build/bin and release." -ForegroundColor Green
}
finally {
  Pop-Location
}
