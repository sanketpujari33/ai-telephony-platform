$ErrorActionPreference = "Stop"

# Get the root directory (parent of deploy folder)
$deployDir = Split-Path -Parent $PSScriptRoot
$rootDir = Split-Path -Parent $deployDir

$services = @(
    "api-gateway",
    "call-router",
    "campaign-manager",
    "media-proxy",
    "orchestrator",
    "sip-ingress",
    "transfer-manager"
)

Write-Host "Starting all AI Telephony Platform services..." -ForegroundColor Cyan
Write-Host "Root directory: $rootDir" -ForegroundColor DarkGray

foreach ($service in $services) {
    Write-Host "Starting $service in a new window..." -ForegroundColor Yellow
    # Start each Go service in a new terminal window
    Start-Process -FilePath "go" -ArgumentList "run ./cmd/$service" -WorkingDirectory $rootDir -WindowStyle Normal
}

Write-Host "All services have been started in separate windows!" -ForegroundColor Green
Write-Host "To stop them, you can close the individual terminal windows." -ForegroundColor DarkGray
