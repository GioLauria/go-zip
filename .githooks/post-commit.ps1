#!/usr/bin/env pwsh
Set-StrictMode -Version Latest

$chlog = "CHANGELOG.md"

# If last commit already modified CHANGELOG.md, skip
if ((git show -1 --name-only) -match [regex]::Escape($chlog)) { exit 0 }
if (-not (Test-Path $chlog)) { exit 0 }

$msg = git log -1 --pretty=%B

$lines = Get-Content $chlog
$out = New-Object System.Collections.Generic.List[string]
$inserted = $false

for ($i = 0; $i -lt $lines.Length; $i++) {
    $out.Add($lines[$i])
    if (-not $inserted -and $lines[$i] -match '^## \[Unreleased\]') {
        $out.Add("")
        $out.Add("- $msg")
        $inserted = $true
    }
}

$out | Set-Content $chlog -Encoding UTF8

git add $chlog
try {
    git commit --amend --no-edit -q
} catch {
    # leave staged if amend fails
}

exit 0
