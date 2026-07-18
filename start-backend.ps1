# Script de démarrage du backend avec admin par défaut
$env:DATA_DIR = "C:\dev\FamilyBass\drive\data"
$env:PORT = "8080"
$env:JWT_SECRET = "dev_secret_change_me_in_production"
$env:ADMIN_EMAIL = "admin@example.com"
$env:ADMIN_PASSWORD = "admin"

Write-Host "🚀 Starting backend on port 8080..." -ForegroundColor Cyan
Write-Host "📧 Admin email: $env:ADMIN_EMAIL" -ForegroundColor Green
Write-Host "🔐 Admin password: $env:ADMIN_PASSWORD" -ForegroundColor Green
Write-Host ""

cd C:\dev\FamilyBass\drive\backend
.\backend.exe

