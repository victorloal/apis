#!/bin/bash

# Configuración
BASE_URL="http://localhost:8001"
TEST_DATA_DIR="./test_data"
mkdir -p $TEST_DATA_DIR

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Funciones de utilidad
log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Verificar si el servidor está ejecutándose
check_server() {
    log "Checking if server is running at $BASE_URL..."
    
    if curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/health" | grep -q "200"; then
        success "Server is running!"
        return 0
    else
        error "Server is not running at $BASE_URL"
        echo
        echo "Please start the server first:"
        echo "cd ~/Documentos/Proyecto_Votaciones/apis/Templates/api_db"
        echo "go run cmd/api/main.go"
        echo
        return 1
    fi
}

test_endpoint() {
    local method=$1
    local url=$2
    local data=$3
    local expected_status=$4
    local description=$5
    
    log "Testing: $description"
    echo "URL: $method $url"
    
    # Crear archivo temporal para la respuesta
    local temp_file=$(mktemp)
    
    if [ "$method" = "GET" ]; then
        http_code=$(curl -s -o "$temp_file" -w "%{http_code}" -X GET "$url")
    elif [ "$method" = "DELETE" ]; then
        http_code=$(curl -s -o "$temp_file" -w "%{http_code}" -X DELETE "$url")
    else
        http_code=$(curl -s -o "$temp_file" -w "%{http_code}" -X "$method" "$url" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi
    
    response_body=$(cat "$temp_file")
    rm -f "$temp_file"
    
    if [ "$http_code" -eq "$expected_status" ]; then
        success "✓ $description - Status: $http_code"
        if [ -n "$response_body" ]; then
            echo "Response: $response_body"
        fi
        echo
        return 0
    else
        error "✗ $description - Expected: $expected_status, Got: $http_code"
        if [ -n "$response_body" ]; then
            echo "Response: $response_body"
        fi
        echo
        return 1
    fi
}

# Pruebas básicas de conectividad
run_basic_tests() {
    log "Running basic connectivity tests..."
    echo
    
    # Test de health check
    test_endpoint "GET" "$BASE_URL/health" "" 200 "Health Check"
    
    # Test de endpoints básicos
    test_endpoint "GET" "$BASE_URL/elections" "" 200 "Get Elections"
    test_endpoint "GET" "$BASE_URL/status" "" 200 "Get Status"
    test_endpoint "GET" "$BASE_URL/authorities" "" 200 "Get Authorities"
}

# Prueba crear una elección simple
test_election_creation() {
    log "Testing election creation..."
    
    local election_data='{
        "name": "Test Election",
        "description": "Test election created by script",
        "encrypted": true,
        "status": 1,
        "start_date": "2024-01-01T00:00:00Z",
        "end_date": "2024-12-31T23:59:59Z"
    }'
    
    test_endpoint "POST" "$BASE_URL/elections" "$election_data" 201 "Create Simple Election"
}

# Prueba crear un status simple
test_status_creation() {
    log "Testing status creation..."
    
    local status_data='{
        "id": 100,
        "name": "test_status"
    }'
    
    test_endpoint "POST" "$BASE_URL/status" "$status_data" 201 "Create Status"
    test_endpoint "GET" "$BASE_URL/status/100" "" 200 "Get Status by ID"
}

# Función principal
main() {
    echo "=== API TEST SCRIPT ==="
    echo
    
    # Verificar dependencias
    if ! command -v curl &> /dev/null; then
        error "curl is required but not installed. Please install curl."
        exit 1
    fi
    
    # Verificar si el servidor está ejecutándose
    if ! check_server; then
        exit 1
    fi
    
    # Ejecutar pruebas básicas
    run_basic_tests
    
    # Ejecutar pruebas de creación
    test_election_creation
    test_status_creation
    
    success "Basic tests completed!"
    echo
    log "If these basic tests pass, you can run the full test suite."
}

# Ejecutar función principal
main