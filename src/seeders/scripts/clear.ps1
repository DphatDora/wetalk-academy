$ErrorActionPreference = "Stop"
Set-Location (Join-Path $PSScriptRoot "..\..")
go run ./cmd/seeder clear @args
