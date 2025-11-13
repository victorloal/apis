#!/usr/bin/env pwsh

Write-Host "🔧 REPARANDO DOCKER COMPOSE Y CERTIFICADOS" -ForegroundColor Cyan

# 1. Limpiar todo
Write-Host "`n1. 🧹 Limpiando contenedores previos..." -ForegroundColor Yellow
docker-compose down --remove-orphans

# 2. Verificar y generar certificados localmente
Write-Host "`n2. 🔐 Generando certificados localmente..." -ForegroundColor Yellow
if (-not (Test-Path "./certs")) {
    New-Item -ItemType Directory -Path "./certs" | Out-Null
}

if (-not (Test-Path "./certs/ca.crt")) {
    Write-Host "   Generando certificados..." -ForegroundColor Cyan
    
    # Generar CA
    openssl genrsa -out ./certs/ca.key 4096
    openssl req -x509 -new -nodes -key ./certs/ca.key -sha256 -days 3650 -out ./certs/ca.crt -subj "/CN=SecureCA"

    # Servicio Auth
    openssl genrsa -out ./certs/auth.key 2048
    openssl req -new -key ./certs/auth.key -out ./certs/auth.csr -subj "/CN=auth-service"
    openssl x509 -req -in ./certs/auth.csr -CA ./certs/ca.crt -CAkey ./certs/ca.key -CAcreateserial -out ./certs/auth.crt -days 365 -sha256

    # Servicio Ledger
    openssl genrsa -out ./certs/ledger.key 2048
    openssl req -new -key ./certs/ledger.key -out ./certs/ledger.csr -subj "/CN=ledger-service"
    openssl x509 -req -in ./certs/ledger.csr -CA ./certs/ca.crt -CAkey ./certs/ca.key -CAcreateserial -out ./certs/ledger.crt -days 365 -sha256

    # Limpiar CSRs
    Remove-Item ./certs/*.csr -ErrorAction SilentlyContinue
    
    Write-Host "   ✅ Certificados generados" -ForegroundColor Green
} else {
    Write-Host "   ✅ Certificados ya existen" -ForegroundColor Green
}

# 3. Reconstruir servicios
Write-Host "`n3. 🚀 Reconstruyendo servicios..." -ForegroundColor Yellow
docker-compose build --no-cache

# 4. Iniciar servicios
Write-Host "`n4. ⬆️ Iniciando servicios..." -ForegroundColor Yellow
docker-compose up -d

# 5. Esperar inicialización
Write-Host "`n5. ⏳ Esperando inicialización (15 segundos)..." -ForegroundColor Yellow
Start-Sleep -Seconds 15

# 6. Verificar estado
Write-Host "`n6. 📊 Estado de servicios:" -ForegroundColor Yellow
docker-compose ps

# 7. Verificar logs
Write-Host "`n7. 📋 Últimos logs:" -ForegroundColor Yellow
docker-compose logs --tail=3 auth-service
docker-compose logs --tail=3 ledger-service

# 8. Probar comunicación
Write-Host "`n8. 🧪 Probando comunicación segura..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# Función para calcular HMAC
function Get-HMACSignature {
    param([string]$Message, [string]$Secret)
    
    $hmac = New-Object System.Security.Cryptography.HMACSHA256
    $hmac.Key = [Text.Encoding]::UTF8.GetBytes($Secret)
    $signature = $hmac.ComputeHash([Text.Encoding]::UTF8.GetBytes($Message))
    return [BitConverter]::ToString($signature).Replace('-', '').ToLower()
}

Write-Host "`n🔐 Probando Auth Service..." -ForegroundColor Cyan
$hmacSig = Get-HMACSignature -Message '{}' -Secret 'hmacsecret'

try {
    $tokenResponse = & curl -s -k `
        --cert ./certs/auth.crt `
        --key ./certs/auth.key `
        -X POST https://localhost:8442/token `
        -H "Content-Type: application/json" `
        -H "X-HMAC: $hmacSig" `
        -d '{}'
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "   ✅ Auth Service responde!" -ForegroundColor Green
        $tokenData = $tokenResponse | ConvertFrom-Json
        Write-Host "   Token: $($tokenData.token)" -ForegroundColor Gray
        
        # Probar Ledger
        Write-Host "`n📋 Probando Ledger Service..." -ForegroundColor Cyan
        $txData = '{"id":"tx_test","amount":100,"from":"A","to":"B"}'
        $txHmac = Get-HMACSignature -Message $txData -Secret 'hmacsecret'
        
        $ledgerResponse = & curl -s -k `
            --cert ./certs/ledger.crt `
            --key ./certs/ledger.key `
            -X POST https://localhost:8443/api/transaction `
            -H "Content-Type: application/json" `
            -H "Authorization: $($tokenData.token)" `
            -H "X-HMAC: $txHmac" `
            -d $txData
            
        if ($LASTEXITCODE -eq 0) {
            Write-Host "   ✅ Ledger Service responde!" -ForegroundColor Green
            Write-Host "   Respuesta: $ledgerResponse" -ForegroundColor Gray
        } else {
            Write-Host "   ❌ Ledger Service error" -ForegroundColor Red
        }
    } else {
        Write-Host "   ❌ Auth Service no responde" -ForegroundColor Red
    }
} catch {
    Write-Host "   ❌ Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n✅ Proceso completado" -ForegroundColor Green