$ErrorActionPreference = "Stop"
Set-Location (Join-Path $PSScriptRoot "..\..")
go run ./cmd/seeder seed @args
