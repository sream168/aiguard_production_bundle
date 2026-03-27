param(
  [switch]$Nsis,
  [switch]$Sign,
  [string]$SignPfxPath = "",
  [string]$SignPfxPassword = "",
  [string]$TimestampUrl = "http://ts.ssl.com",
  [string]$WorkspaceRoot = ""
)

$ErrorActionPreference = 'Stop'

function Add-UniqueItem([System.Collections.Generic.List[object]]$list, [string]$key, [object]$item) {
  foreach ($existing in $list) {
    if ($existing.Key -eq $key) {
      return
    }
  }
  $list.Add($item)
}

function Get-SdkRoots() {
  $roots = [System.Collections.Generic.List[object]]::new()

  $candidates = @(
    @{ Kit = '8.1'; Root = $env:WINDOWS_KITS_81_PATH },
    @{ Kit = '8.1'; Root = 'C:\Program Files (x86)\Windows Kits\8.1' },
    @{ Kit = '10';  Root = $env:WINDOWS_KITS_10_PATH },
    @{ Kit = '10';  Root = 'C:\Program Files (x86)\Windows Kits\10' },
    @{ Kit = '11';  Root = $env:WINDOWS_KITS_11_PATH },
    @{ Kit = '11';  Root = 'C:\Program Files (x86)\Windows Kits\11' }
  )

  foreach ($candidate in $candidates) {
    $root = [string]$candidate.Root
    if ([string]::IsNullOrWhiteSpace($root)) {
      continue
    }
    $root = $root.TrimEnd('\\')
    if (-not (Test-Path $root)) {
      continue
    }
    Add-UniqueItem $roots "$($candidate.Kit)|$root" ([pscustomobject]@{
      Key  = "$($candidate.Kit)|$root"
      Kit  = $candidate.Kit
      Root = $root
    })
  }

  return $roots.ToArray()
}

function Resolve-KitFromVersion([string]$versionKey, [string]$fallbackKit) {
  if ([string]::IsNullOrWhiteSpace($versionKey)) {
    return $fallbackKit
  }

  try {
    $ver = [version]$versionKey.TrimEnd('\\')
    if ($ver.Build -ge 22000) {
      return '11'
    }
    return '10'
  }
  catch {
    return $fallbackKit
  }
}

function Get-WindowsSdkCandidates() {
  $items = [System.Collections.Generic.List[object]]::new()
  $roots = Get-SdkRoots

  foreach ($rootInfo in $roots) {
    $binRoot = Join-Path $rootInfo.Root 'bin'
    if (-not (Test-Path $binRoot)) {
      continue
    }

    if ($rootInfo.Kit -eq '8.1') {
      foreach ($arch in @('x64', 'x86')) {
        $candidateDir = Join-Path $binRoot $arch
        $signtool = Join-Path $candidateDir 'signtool.exe'
        if (Test-Path $signtool) {
          Add-UniqueItem $items "8.1|$candidateDir" ([pscustomobject]@{
            Key          = "8.1|$candidateDir"
            Kit          = '8.1'
            Root         = $rootInfo.Root
            VersionKey   = '8.1'
            BinPath      = $candidateDir
            SignToolPath = $signtool
          })
        }
      }
      continue
    }

    $dirs = Get-ChildItem $binRoot -Directory -ErrorAction SilentlyContinue |
      Sort-Object Name -Descending

    foreach ($dir in $dirs) {
      foreach ($arch in @('x64', 'x86')) {
        $candidateDir = Join-Path $dir.FullName $arch
        $signtool = Join-Path $candidateDir 'signtool.exe'
        if (-not (Test-Path $signtool)) {
          continue
        }

        $detectedKit = Resolve-KitFromVersion $dir.Name $rootInfo.Kit
        Add-UniqueItem $items "$detectedKit|$candidateDir" ([pscustomobject]@{
          Key          = "$detectedKit|$candidateDir"
          Kit          = $detectedKit
          Root         = $rootInfo.Root
          VersionKey   = $dir.Name
          BinPath      = $candidateDir
          SignToolPath = $signtool
        })
      }
    }
  }

  return $items.ToArray() | Sort-Object @{ Expression = {
      switch ($_.Kit) {
        '8.1' { 0 }
        '10'  { 1 }
        '11'  { 2 }
        default { 9 }
      }
    }
  }, @{ Expression = { $_.VersionKey }; Descending = $true }
}

function Use-WindowsSdk([object]$candidate) {
  if ($null -eq $candidate) {
    Write-Host 'No explicit Windows SDK candidate detected. Use current environment.' -ForegroundColor Yellow
    return
  }

  $root = "$($candidate.Root)"
  $versionKey = "$($candidate.VersionKey)"

  $env:WindowsSdkDir = $root + '\\'
  $env:WindowsSdkVerBinPath = $candidate.BinPath + '\\'

  if ($candidate.Kit -eq '8.1') {
    $env:WindowsSDKVersion = ''
    $env:UCRTVersion = ''
    $env:UniversalCRTSdkDir = ''
  }
  else {
    $env:WindowsSDKVersion = $versionKey.TrimEnd('\\') + '\\'
    $env:UCRTVersion = $versionKey.TrimEnd('\\') + '\\'
    $env:UniversalCRTSdkDir = $root + '\\'
  }

  Write-Host ("Using Windows Kit {0} | Version {1} | Bin {2}" -f $candidate.Kit, $candidate.VersionKey, $candidate.BinPath) -ForegroundColor Yellow
}

function Invoke-WailsBuildWithFallback([switch]$BuildNsis) {
  $sdkCandidates = Get-WindowsSdkCandidates
  if ($null -eq $sdkCandidates -or $sdkCandidates.Count -eq 0) {
    Write-Host 'No Windows SDK candidate found from Windows Kits 8.1/10/11. Build once with current environment.' -ForegroundColor Yellow
    wails build -platform windows/amd64
    if ($BuildNsis) {
      wails build -platform windows/amd64 -nsis
    }
    return
  }

  $failures = [System.Collections.Generic.List[string]]::new()

  foreach ($candidate in $sdkCandidates) {
    try {
      Use-WindowsSdk $candidate
      Write-Host ("Trying build with Windows Kit {0} ({1})" -f $candidate.Kit, $candidate.VersionKey) -ForegroundColor Cyan
      wails build -platform windows/amd64
      if ($BuildNsis) {
        wails build -platform windows/amd64 -nsis
      }
      Write-Host ("Build succeeded with Windows Kit {0} ({1})" -f $candidate.Kit, $candidate.VersionKey) -ForegroundColor Green
      return
    }
    catch {
      $message = "Windows Kit $($candidate.Kit) [$($candidate.VersionKey)] build failed: $($_.Exception.Message)"
      $failures.Add($message)
      Write-Host $message -ForegroundColor Yellow
    }
  }

  throw ("All Windows Kit build attempts failed.`n" + ($failures -join "`n"))
}

function Resolve-SignTool() {
  $sdkCandidates = Get-WindowsSdkCandidates
  foreach ($candidate in $sdkCandidates) {
    if (Test-Path $candidate.SignToolPath) {
      Write-Host ("Using signtool from Windows Kit {0} ({1})" -f $candidate.Kit, $candidate.VersionKey) -ForegroundColor Yellow
      return $candidate.SignToolPath
    }
  }
  return $null
}

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

  Write-Host "`n[4/5] Build Windows EXE (fallback: Kit 8.1 -> 10 -> 11)" -ForegroundColor Cyan
  Invoke-WailsBuildWithFallback -BuildNsis:$Nsis

  if ($Sign) {
    if ([string]::IsNullOrWhiteSpace($SignPfxPath) -or [string]::IsNullOrWhiteSpace($SignPfxPassword)) {
      throw "-Sign requires both -SignPfxPath and -SignPfxPassword"
    }

    $signtool = Resolve-SignTool
    if ([string]::IsNullOrWhiteSpace($signtool)) {
      throw "signtool.exe was not found in Windows Kits 8.1/10/11. Please install the Windows SDK."
    }

    Get-ChildItem ./build/bin -Filter *.exe | ForEach-Object {
      Write-Host "Sign $($_.FullName)" -ForegroundColor Yellow
      & $signtool sign /fd sha256 /tr $TimestampUrl /f $SignPfxPath /p $SignPfxPassword $_.FullName
    }
  }

  Write-Host "`n[5/5] Archive artifacts" -ForegroundColor Cyan
  ./scripts/package_release.ps1

  Write-Host "`nDone. Artifacts are available in build/bin and release." -ForegroundColor Green
}
finally {
  Pop-Location
}
