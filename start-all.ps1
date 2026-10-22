# start-all.ps1
# Multi-process launcher for Smart Pharmacy on Windows

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "  Starting Smart Pharmacy and Prescription Assistant Stack " -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

$scriptRoot = $PSScriptRoot

# 1. Start F# Validation Microservice
Write-Host "[1/4] Starting F# Validation Service on port 5005..." -ForegroundColor Green
$fsharpProcess = Start-Process -FilePath "dotnet" -ArgumentList "run" -WorkingDirectory "$scriptRoot\validation-fsharp" -PassThru

# 2. Start Python OCR Microservice
Write-Host "[2/4] Starting Python OCR Service on port 8000..." -ForegroundColor Green
$pythonProcess = Start-Process -FilePath "python" -ArgumentList "main.py" -WorkingDirectory "$scriptRoot\ocr-python" -PassThru

# 3. Start Go Backend
Write-Host "[3/4] Starting Go Backend on port 8080..." -ForegroundColor Green
$goProcess = Start-Process -FilePath "go" -ArgumentList "run", ".\cmd\server" -WorkingDirectory "$scriptRoot\backend-go" -PassThru

# 4. Start Frontend
Write-Host "[4/4] Starting TypeScript/Vite Frontend on port 5173..." -ForegroundColor Green
$frontendProcess = Start-Process -FilePath "npm" -ArgumentList "run", "dev" -WorkingDirectory "$scriptRoot\frontend" -PassThru

Start-Sleep -Seconds 3

Write-Host ""
Write-Host "All services have been initiated!" -ForegroundColor Cyan
Write-Host "----------------------------------------------------------"
Write-Host "  Web Frontend:        http://localhost:5173" -ForegroundColor Yellow
Write-Host "  Go REST Backend:     http://localhost:8080" -ForegroundColor White
Write-Host "  F# Validation Svc:   http://localhost:5005" -ForegroundColor White
Write-Host "  Python OCR Svc:      http://localhost:8000" -ForegroundColor White
Write-Host "----------------------------------------------------------"
Write-Host "Demo Credentials:" -ForegroundColor Yellow
Write-Host "  Customer:   customer@pharmacy.local / CustomerPass123!"
Write-Host "  Pharmacist: staff@pharmacy.local    / StaffPass123!"
Write-Host "  Admin:      admin@pharmacy.local    / AdminPass123!"
Write-Host ""
Write-Host "Press Ctrl+C or close the terminal to stop."

# Wait for exit
try {
    Wait-Process -Id $frontendProcess.Id
} finally {
    Write-Host "Stopping all background processes..." -ForegroundColor Yellow
    Stop-Process -Id $fsharpProcess.Id -ErrorAction SilentlyContinue
    Stop-Process -Id $pythonProcess.Id -ErrorAction SilentlyContinue
    Stop-Process -Id $goProcess.Id -ErrorAction SilentlyContinue
    Stop-Process -Id $frontendProcess.Id -ErrorAction SilentlyContinue
}
