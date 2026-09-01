$ErrorActionPreference = "Stop"

# Get the root directory
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

Write-Host "Configuring Docker client to use Minikube's Docker daemon..." -ForegroundColor Cyan
& minikube docker-env | Invoke-Expression

Write-Host "Building Docker images for AI Telephony Platform..." -ForegroundColor Cyan

Set-Location $rootDir

foreach ($service in $services) {
    Write-Host "Building image for $service..." -ForegroundColor Yellow
    $imageName = "ai-telephony-platform/${service}:latest"
    
    # Build the image using the multi-stage Dockerfile
    docker build -f deploy/docker/Dockerfile --build-arg SERVICE_NAME=$service -t $imageName .
    
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to build image for $service"
        exit 1
    }
    
    Write-Host "Successfully built $imageName" -ForegroundColor Green
}

Write-Host "All images built and loaded into Minikube successfully!" -ForegroundColor Green
