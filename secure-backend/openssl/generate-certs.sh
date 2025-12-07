#!/bin/bash
set -e

echo "🔐 Generando certificados mTLS..."

mkdir -p /certs

# Crear CA
openssl genrsa -out /certs/ca.key 4096
openssl req -x509 -new -nodes -key /certs/ca.key -sha256 -days 3650 -out /certs/ca.crt \
  -subj "/C=US/ST=State/L=City/O=SecureBackend/CN=SecureCA"

# Servicio Auth
openssl genrsa -out /certs/auth.key 2048
openssl req -new -key /certs/auth.key -out /certs/auth.csr \
  -subj "/C=US/ST=State/L=City/O=SecureBackend/CN=auth-service"
openssl x509 -req -in /certs/auth.csr -CA /certs/ca.crt -CAkey /certs/ca.key \
  -CAcreateserial -out /certs/auth.crt -days 365 -sha256 \
  -extensions v3_req -extfile <(
    cat <<-EOF
[ v3_req ]
subjectAltName = @alt_names
[ alt_names ]
DNS.1 = auth-service
DNS.2 = auth-service.docker
DNS.3 = localhost
IP.1 = 127.0.0.1
EOF
  )

# Servicio Ledger
openssl genrsa -out /certs/ledger.key 2048
openssl req -new -key /certs/ledger.key -out /certs/ledger.csr \
  -subj "/C=US/ST=State/L=City/O=SecureBackend/CN=ledger-service"
openssl x509 -req -in /certs/ledger.csr -CA /certs/ca.crt -CAkey /certs/ca.key \
  -CAcreateserial -out /certs/ledger.crt -days 365 -sha256 \
  -extensions v3_req -extfile <(
    cat <<-EOF
[ v3_req ]
subjectAltName = @alt_names
[ alt_names ]
DNS.1 = ledger-service
DNS.2 = ledger-service.docker
DNS.3 = localhost
IP.1 = 127.0.0.1
EOF
  )

# Limpiar CSRs
rm /certs/*.csr

echo "✅ Certificados generados en /certs/"
echo "📜 Lista de archivos:"
ls -la /certs/