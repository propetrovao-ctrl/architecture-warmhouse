# PowerShell script to run the application with database migrations

Write-Host "Starting Smart Home API with migrations..." -ForegroundColor Green

# Wait for database to be ready
Write-Host "Waiting for database to be ready..." -ForegroundColor Yellow
do {
    try {
        $connection = New-Object System.Data.SqlClient.SqlConnection
        $connection.ConnectionString = "Server=postgres,1433;Database=smarthome;User Id=postgres;Password=postgres;"
        $connection.Open()
        $connection.Close()
        $dbReady = $true
    }
    catch {
        Write-Host "Database is unavailable - sleeping" -ForegroundColor Yellow
        Start-Sleep -Seconds 2
        $dbReady = $false
    }
} while (-not $dbReady)

Write-Host "Database is ready!" -ForegroundColor Green

# Run migrations
Write-Host "Running database migrations..." -ForegroundColor Yellow
.\migrate.exe -database-url="$env:DATABASE_URL" -command=up

Write-Host "Migrations completed successfully!" -ForegroundColor Green

# Start the application
Write-Host "Starting Smart Home API..." -ForegroundColor Green
.\smarthome.exe
