#!/bin/bash

# Configuración
API_DIR="$HOME/Documentos/Proyecto_Votaciones/apis/Templates/api_db"
TEST_SCRIPT="./test_api_improved.sh"

# Colores
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo "=== API Auto Test ==="
echo

# Verificar si el directorio existe
if [ ! -d "$API_DIR" ]; then
    echo -e "${RED}Error: API directory not found at $API_DIR${NC}"
    exit 1
fi

# Navegar al directorio
cd "$API_DIR"

# Verificar si go.mod existe
if [ ! -f "go.mod" ]; then
    echo -e "${RED}Error: go.mod not found. Are you in the correct directory?${NC}"
    exit 1
fi

# Verificar si la API ya está ejecutándose
if curl -s http://localhost:8001/health > /dev/null; then
    echo -e "${GREEN}API is already running!${NC}"
    echo "Running tests..."
    ./test_api_improved.sh
else
    echo "Starting API server..."
    
    # Ejecutar la API en segundo plano
    go run cmd/api/main.go &
    API_PID=$!
    
    # Esperar a que la API esté lista
    echo "Waiting for API to start..."
    sleep 3
    
    # Verificar si la API se inició correctamente
    if curl -s http://localhost:8001/health > /dev/null; then
        echo -e "${GREEN}API started successfully!${NC}"
        echo "Running tests..."
        ./test_api_improved.sh
        
        # Detener la API después de las pruebas
        echo "Stopping API server..."
        kill $API_PID
    else
        echo -e "${RED}Failed to start API${NC}"
        kill $API_PID
        exit 1
    fi
fi