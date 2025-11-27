#!/bin/bash

echo "🧪 Probando comunicación segura entre servicios..."

# 1. Probar que los certificados se generaron
echo "📜 Verificando certificados..."
ls -la ./certs/

# 2. Probar auth-service
echo "🔑 Probando Auth Service..."


curl -k --cert ./certs/auth.crt --key ./certs/auth.key \
  -X POST https://localhost:8442/token \
  -H "Content-Type: application/json" \
  -H "X-HMAC: $(echo -n '{}' | openssl dgst -sha256 -hmac 'hmacsecret' | cut -d' ' -f2)" \
  -d '{}'

echo -e "\n"

# 3. Probar ledger-service
echo "📋 Probando Ledger Service..."
TOKEN=$(curl -s -k --cert ./certs/auth.crt --key ./certs/auth.key \
  -X POST https://localhost:8442/token \
  -H "Content-Type: application/json" \
  -H "X-HMAC: $(echo -n '{}' | openssl dgst -sha256 -hmac 'hmacsecret' | cut -d' ' -f2)" \
  -d '{}' | jq -r .token)

echo "Token obtenido: $TOKEN"

curl -k --cert ./certs/ledger.crt --key ./certs/ledger.key \
  -X POST https://localhost:8443/api/transaction \
  -H "Content-Type: application/json" \
  -H "Authorization: $TOKEN" \
  -H "X-HMAC: $(echo -n '{"id":"tx123","amount":100,"from":"A","to":"B"}' | openssl dgst -sha256 -hmac 'hmacsecret' | cut -d' ' -f2)" \
  -d '{"id":"tx123","amount":100,"from":"A","to":"B"}'

echo -e "\n✅ Pruebas completadas"