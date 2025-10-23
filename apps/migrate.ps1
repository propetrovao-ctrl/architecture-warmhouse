# Migration management script for Smart Home API (PowerShell)

param(
    [string]$DatabaseUrl = "postgres://postgres:postgres@localhost:5432/smarthome",
    [string]$Command = "up",
    [switch]$Help
)

if ($Help) {
    Write-Host "Usage: .\migrate.ps1 [OPTIONS]" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Options:" -ForegroundColor Yellow
    Write-Host "  -DatabaseUrl URL    Database connection URL" -ForegroundColor Yellow
    Write-Host "  -Command COMMAND    Migration command (up, status)" -ForegroundColor Yellow
    Write-Host "  -Help              Show this help message" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Yellow
    Write-Host "  .\migrate.ps1                                    # Run migrations with default settings" -ForegroundColor Yellow
    Write-Host "  .\migrate.ps1 -Command status                   # Show migration status" -ForegroundColor Yellow
    Write-Host "  .\migrate.ps1 -DatabaseUrl postgres://user:pass@host/db    # Use custom database URL" -ForegroundColor Yellow
    exit 0
}

Write-Host "Smart Home API Migration Tool" -ForegroundColor Yellow
Write-Host "Database URL: $DatabaseUrl" -ForegroundColor Green
Write-Host "Command: $Command" -ForegroundColor Green
Write-Host ""

# Check if we're running in Docker
if (Test-Path "C:\.dockerenv") {
    Write-Host "Running in Docker container" -ForegroundColor Green
    .\migrate.exe -database-url="$DatabaseUrl" -command="$Command"
} else {
    Write-Host "Running locally" -ForegroundColor Green
    Set-Location smart_home
    go run cmd/migrate/main.go -database-url="$DatabaseUrl" -command="$Command"
}

Write-Host "Migration completed successfully!" -ForegroundColor Green
