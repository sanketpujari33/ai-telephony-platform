$ErrorActionPreference = "Stop"

# Get the directory where the helm charts are located
$deployDir = Split-Path -Parent $PSScriptRoot
$helmDir = Join-Path -Path $deployDir -ChildPath "helm"
$namespace = "ai-telephony"

$services = @(
    "api-gateway",
    "call-router",
    "campaign-manager",
    "media-proxy",
    "orchestrator",
    "sip-ingress",
    "transfer-manager"
)

Write-Host "Starting full local deployment..." -ForegroundColor Cyan

# 1. Build Docker images and load into Minikube
& (Join-Path $PSScriptRoot "build-images.ps1")

# 2. Install Observability Stack (Prometheus & Grafana)
& (Join-Path $PSScriptRoot "install-observability.ps1")

Write-Host "Installing all AI Telephony Platform Helm charts to namespace: $namespace..." -ForegroundColor Cyan

foreach ($service in $services) {
    Write-Host "Installing/Upgrading helm release: $service..." -ForegroundColor Yellow
    $chartPath = Join-Path -Path $helmDir -ChildPath $service
    
    # Run the helm install/upgrade command with CRD-dependent features disabled for local minikube compatibility
    helm upgrade --install $service $chartPath --namespace $namespace --create-namespace `
        --set istio.destinationRule.enabled=false `
        --set istio.virtualService.enabled=false `
        --set monitoring.serviceMonitor.enabled=false
}

Write-Host "All Helm charts have been successfully installed!" -ForegroundColor Green
Write-Host "Run 'kubectl get pods -n $namespace' to check the status of your deployments." -ForegroundColor DarkGray
