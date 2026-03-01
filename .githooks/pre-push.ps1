#!/usr/bin/env pwsh
Set-StrictMode -Version Latest

# pre-push PowerShell hook: detect tags being pushed and move Unreleased -> <tag>

while ($true) {
    $line = [Console]::In.ReadLine()
    if ($null -eq $line) { break }
    $parts = $line -split "\s+"
    if ($parts.Length -lt 1) { continue }
    $localRef = $parts[0]
    if ($localRef -like 'refs/tags/*') {
        $tag = $localRef -replace '^refs/tags/', ''
        $chlog = "CHANGELOG.md"
        if (Test-Path $chlog) {
            $date = (Get-Date).ToString('yyyy-MM-dd')
            $tmp = "$chlog.tmp"
            $new = "$chlog.new"
            # Replace first Unreleased with tag header
            (Get-Content $chlog) | ForEach-Object -Begin {$repl=$false} -Process {
                if (-not $repl -and $_ -eq '## [Unreleased]') { "$([string]::Format('## [{0}] - {1}', $tag, $date))"; $repl=$true } else { $_ }
            } | Set-Content $tmp -Encoding UTF8
            # Insert fresh Unreleased after top line
            $header = Get-Content -Path $tmp -TotalCount 1
            $rest = Get-Content -Path $tmp | Select-Object -Skip 1
            @($header) + @('## [Unreleased]','') + $rest | Set-Content -Path $new -Encoding UTF8
            Move-Item -Force $new $chlog
            Remove-Item -Force $tmp -ErrorAction SilentlyContinue
            git add $chlog
            git commit -m "chore(release): move Unreleased -> $tag" -q -q -q -ErrorAction SilentlyContinue
        }
    }
}

exit 0
