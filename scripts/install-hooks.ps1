# PowerShell script to install git hooks on Windows
Set-StrictMode -Version Latest
git config core.hooksPath .githooks
Write-Host "Git hooks installed. (core.hooksPath = .githooks)"
