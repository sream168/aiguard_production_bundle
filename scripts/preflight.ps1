param()

$ErrorActionPreference = 'Stop'

function Test-Command($name) {
  $cmd = Get-Command $name -ErrorAction SilentlyContinue
  if ($null -eq $cmd) {
    Write-Host "[FAIL] Missing command: $name" -ForegroundColor Red
    return $false
  }
  Write-Host "[ OK ] $name -> $($cmd.Source)" -ForegroundColor Green
  return $true
}

$all = $true
$all = (Test-Command go) -and $all
$all = (Test-Command git) -and $all
$all = (Test-Command npm) -and $all
$all = (Test-Command wails) -and $all

try {
  $wv2 = Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\EdgeUpdate\Clients\{F1E7E5D0-0B2B-4B15-9A0B-FA6FE7B9D28B}' -ErrorAction Stop
  Write-Host "[ OK ] WebView2 Runtime detected" -ForegroundColor Green
} catch {
  Write-Host "[WARN] WebView2 Runtime was not detected. Please install it before building." -ForegroundColor Yellow
}

if ($all) {
  Write-Host "`nRunning wails doctor..." -ForegroundColor Cyan
  wails doctor
} else {
  Write-Host "`nDependencies are incomplete. Please install them and try again." -ForegroundColor Red
  exit 1
}
