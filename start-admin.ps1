$host.UI.RawUI.WindowTitle = "C2F Admin Web App"
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
Set-Location "$scriptDir\src\c2f_admin_webapp"
Write-Host "Starting C2F Admin Web App on http://localhost:5005..." -ForegroundColor Green
Write-Host "Connected to AWS RDS MySQL (C2F_NEW)" -ForegroundColor Cyan
& ".\c2f_admin_webapp.exe"

