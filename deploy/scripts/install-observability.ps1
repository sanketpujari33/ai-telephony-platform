$ErrorActionPreference = "Stop"

Write-Host "Setting up Observability Stack (Prometheus & Grafana)..." -ForegroundColor Cyan

# Add the Prometheus Community Helm repository
Write-Host "Adding prometheus-community Helm repo..." -ForegroundColor DarkGray
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

$namespace = "monitoring"

# Install kube-prometheus-stack
Write-Host "Installing kube-prometheus-stack..." -ForegroundColor Yellow
# We disable some components that aren't needed for local development to save resources
helm upgrade --install prometheus-stack prometheus-community/kube-prometheus-stack `
    --namespace $namespace --create-namespace `
    --set alertmanager.enabled=false `
    --set prometheus-node-exporter.enabled=false `
    --set kube-state-metrics.enabled=false `
    --set grafana.adminPassword="admin" `
    --wait

# Apply custom Prometheus alerts
$alertsFile = Join-Path $PSScriptRoot "..\prometheus\alerts.yml"
if (Test-Path $alertsFile) {
    Write-Host "Applying custom Prometheus alerts from alerts.yml..." -ForegroundColor DarkGray
    # We create a PrometheusRule custom resource from the alerts.yml
    $alertsYaml = (Get-Content $alertsFile) -join "`n"
    $alertsYamlIndented = $alertsYaml -replace '(?m)^', '  '
    $prometheusRule = @"
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: ai-telephony-alerts
  namespace: monitoring
  labels:
    app: kube-prometheus-stack
    release: prometheus-stack
spec:
$alertsYamlIndented
"@
    $prometheusRule | kubectl apply -f -
}

# Apply custom Grafana dashboards
$dashboardsDir = Join-Path $PSScriptRoot "..\grafana\dashboards"
if (Test-Path $dashboardsDir) {
    Write-Host "Creating ConfigMap for custom Grafana dashboards..." -ForegroundColor DarkGray
    $dashboards = Get-ChildItem -Path $dashboardsDir -Filter "*.json"
    foreach ($dashboard in $dashboards) {
        $name = "grafana-dashboard-" + $dashboard.BaseName.ToLower() -replace '[^a-z0-9-]', '-'
        kubectl create configmap $name --from-file="$($dashboard.FullName)" -n $namespace --dry-run=client -o yaml | kubectl apply -f -
        # Label the configmap so the Grafana sidecar automatically loads it
        kubectl label configmap $name grafana_dashboard=1 -n $namespace --overwrite
    }
}

Write-Host "Observability Stack installed successfully!" -ForegroundColor Green
Write-Host "To access Grafana, run: kubectl port-forward svc/prometheus-stack-grafana -n $namespace 3000:80" -ForegroundColor Green
Write-Host "Default login: admin / admin" -ForegroundColor Green
