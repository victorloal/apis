#!/bin/bash

# Configuración
BASE_URL="http://localhost:8000"
TEST_DATA_DIR="./test_data"
mkdir -p $TEST_DATA_DIR

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variables para almacenar IDs creados durante las pruebas
ELECTION_ID=""
AUTHORITY_ID=""
VOTER_ID=""
CANDIDATE_ID=""
BALLOT_ID="test-ballot-$(date +%s)"
TALLY_RESULT_ID=""
AUDIT_CONFIG_ID=""

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

test_endpoint() {
    local method=$1
    local url=$2
    local data=$3
    local expected_status=$4
    local description=$5
    
    log "Testing: $description"
    echo "URL: $method $url"
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "%{http_code}" -X GET "$url" -o response_body.json)
    elif [ "$method" = "DELETE" ]; then
        response=$(curl -s -w "%{http_code}" -X DELETE "$url" -o response_body.json)
    else
        response=$(curl -s -w "%{http_code}" -X "$method" "$url" -H "Content-Type: application/json" -d "$data" -o response_body.json)
    fi
    
    http_code=${response: -3}
    response_body=$(cat response_body.json)
    rm -f response_body.json
    
    if [ "$http_code" -eq "$expected_status" ]; then
        success "✓ $description - Status: $http_code"
        echo "Response: $response_body"
        echo
        return 0
    else
        error "✗ $description - Expected: $expected_status, Got: $http_code"
        echo "Response: $response_body"
        echo
        return 1
    fi
}

# Generar datos de prueba
generate_test_data() {
    # Datos para election
    cat > $TEST_DATA_DIR/election.json << EOF
{
    "name": "Test Election $(date +%Y%m%d%H%M%S)",
    "description": "This is a test election created by API test script",
    "encrypted": true,
    "status": 1,
    "start_date": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
    "end_date": "$(date -u -d "+30 days" +"%Y-%m-%dT%H:%M:%SZ")"
}
EOF

    # Datos para authority
    cat > $TEST_DATA_DIR/authority.json << EOF
{
    "cc": 123456789,
    "name": "Test Authority $(date +%s)",
    "email": "authority$(date +%s)@test.com",
    "password": "testpassword123",
    "s_key": "dGVzdC1zZWNyZXQta2V5"
}
EOF

    # Datos para voter
    cat > $TEST_DATA_DIR/voter.json << EOF
{
    "token": "voter-token-$(date +%s)",
    "vote_status": false,
    "verification_hash": "dGVzdC1oYXNo",
    "is_active": true
}
EOF

    # Datos para candidate
    cat > $TEST_DATA_DIR/candidate.json << EOF
{
    "name": "Test Candidate $(date +%s)",
    "description": "Test candidate description",
    "photo_url": "https://example.com/photo$(date +%s).jpg",
    "candidate_order": 1
}
EOF

    # Datos para ballot
    cat > $TEST_DATA_DIR/ballot.json << EOF
{
    "id": "$BALLOT_ID",
    "vote": {
        "candidate_id": 1,
        "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
        "encrypted_vote": "encrypted-data-here"
    },
    "voting_device_fingerprint": "device-fingerprint-123",
    "ip_address": "192.168.1.100"
}
EOF

    # Datos para homomorphic key
    cat > $TEST_DATA_DIR/homomorphic_key.json << EOF
{
    "p_key": "cHVibGljLWtleS1kYXRh",
    "params": {
        "scheme": "BFV",
        "poly_modulus_degree": 4096,
        "coeff_modulus_bits": [40, 40, 40],
        "plain_modulus": 1024
    }
}
EOF

    # Datos para tally result
    cat > $TEST_DATA_DIR/tally_result.json << EOF
{
    "results": {
        "candidate_1": 150,
        "candidate_2": 200
    },
    "total_votes": 350,
    "computed_by": "test-script",
    "proof": {
        "verification": "success",
        "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
    }
}
EOF

    # Datos para status
    cat > $TEST_DATA_DIR/status.json << EOF
{
    "id": 999,
    "name": "test_status_$(date +%s)"
}
EOF

    # Datos para audit config
    cat > $TEST_DATA_DIR/audit_config.json << EOF
{
    "enable_ballot_audit": true,
    "enable_access_logs": true
}
EOF

    # Datos para audit log
    cat > $TEST_DATA_DIR/audit_log.json << EOF
{
    "action": "test_action",
    "user_type": "test_user",
    "user_id": "test_user_123",
    "ip_address": "192.168.1.100",
    "user_agent": "Test-Script/1.0",
    "details": {
        "test_field": "test_value",
        "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
    }
}
EOF
}

# Función principal de pruebas
run_tests() {
    log "Starting API tests..."
    echo "Base URL: $BASE_URL"
    echo
    
    # Generar datos de prueba
    generate_test_data
    
    # 1. Health Check
    test_endpoint "GET" "$BASE_URL/health" "" 200 "Health Check"
    
    # 2. Status endpoints
    log "=== TESTING STATUS ENDPOINTS ==="
    test_endpoint "POST" "$BASE_URL/status" "$(cat $TEST_DATA_DIR/status.json)" 201 "Create Status"
    test_endpoint "GET" "$BASE_URL/status" "" 200 "Get All Status"
    test_endpoint "GET" "$BASE_URL/status/999" "" 200 "Get Status by ID"
    test_endpoint "PUT" "$BASE_URL/status/999" '{"name":"updated_status"}' 200 "Update Status"
    test_endpoint "GET" "$BASE_URL/status/name/updated_status" "" 200 "Get Status by Name"
    
    # 3. Election endpoints
    log "=== TESTING ELECTION ENDPOINTS ==="
    test_endpoint "POST" "$BASE_URL/elections" "$(cat $TEST_DATA_DIR/election.json)" 201 "Create Election"
    
    # Obtener el ID de la elección creada
    ELECTION_RESPONSE=$(curl -s -X GET "$BASE_URL/elections")
    ELECTION_ID=$(echo $ELECTION_RESPONSE | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
    if [ -n "$ELECTION_ID" ]; then
        success "Election ID: $ELECTION_ID"
    fi
    
    test_endpoint "GET" "$BASE_URL/elections" "" 200 "Get All Elections"
    test_endpoint "GET" "$BASE_URL/elections/$ELECTION_ID" "" 200 "Get Election by ID"
    test_endpoint "PUT" "$BASE_URL/elections/$ELECTION_ID" '{"name":"Updated Election Name"}' 200 "Update Election"
    
    # 4. Authorities endpoints
    log "=== TESTING AUTHORITIES ENDPOINTS ==="
    # Actualizar el election_id en los datos de authority
    sed -i "s/\"election\": \"[^\"]*\"/\"election\": \"$ELECTION_ID\"/" $TEST_DATA_DIR/authority.json
    
    test_endpoint "POST" "$BASE_URL/authorities" "$(cat $TEST_DATA_DIR/authority.json)" 201 "Create Authority"
    
    # Obtener el ID de la autoridad creada
    AUTHORITIES_RESPONSE=$(curl -s -X GET "$BASE_URL/authorities/election/$ELECTION_ID")
    AUTHORITY_ID=$(echo $AUTHORITIES_RESPONSE | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
    if [ -n "$AUTHORITY_ID" ]; then
        success "Authority ID: $AUTHORITY_ID"
    fi
    
    test_endpoint "GET" "$BASE_URL/authorities" "" 200 "Get All Authorities"
    test_endpoint "GET" "$BASE_URL/authorities/$AUTHORITY_ID" "" 200 "Get Authority by ID"
    test_endpoint "GET" "$BASE_URL/authorities/election/$ELECTION_ID" "" 200 "Get Authorities by Election"
    test_endpoint "PUT" "$BASE_URL/authorities/$AUTHORITY_ID" '{"name":"Updated Authority Name"}' 200 "Update Authority"
    test_endpoint "PATCH" "$BASE_URL/authorities/$AUTHORITY_ID" '{"email":"updated@test.com"}' 200 "Partial Update Authority"
    
    # 5. Voters endpoints
    log "=== TESTING VOTERS ENDPOINTS ==="
    # Actualizar el elections en los datos de voter
    sed -i "s/\"elections\": \"[^\"]*\"/\"elections\": \"$ELECTION_ID\"/" $TEST_DATA_DIR/voter.json
    
    test_endpoint "POST" "$BASE_URL/voters" "$(cat $TEST_DATA_DIR/voter.json)" 201 "Create Voter"
    
    # Obtener el ID del voter creado
    VOTERS_RESPONSE=$(curl -s -X GET "$BASE_URL/voters/election/$ELECTION_ID")
    VOTER_ID=$(echo $VOTERS_RESPONSE | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
    if [ -n "$VOTER_ID" ]; then
        success "Voter ID: $VOTER_ID"
    fi
    
    test_endpoint "GET" "$BASE_URL/voters" "" 200 "Get All Voters"
    test_endpoint "GET" "$BASE_URL/voters/$VOTER_ID" "" 200 "Get Voter by ID"
    test_endpoint "GET" "$BASE_URL/voters/election/$ELECTION_ID" "" 200 "Get Voters by Election"
    test_endpoint "PUT" "$BASE_URL/voters/$VOTER_ID" '{"vote_status":true}' 200 "Update Voter"
    test_endpoint "PUT" "$BASE_URL/voters/$VOTER_ID/vote-status" '{"status":true}' 200 "Update Vote Status"
    
    # 6. Candidates endpoints
    log "=== TESTING CANDIDATES ENDPOINTS ==="
    # Actualizar el elections en los datos de candidate
    sed -i "s/\"elections\": \"[^\"]*\"/\"elections\": \"$ELECTION_ID\"/" $TEST_DATA_DIR/candidate.json
    
    test_endpoint "POST" "$BASE_URL/candidates" "$(cat $TEST_DATA_DIR/candidate.json)" 201 "Create Candidate"
    
    # Obtener el ID del candidate creado
    CANDIDATES_RESPONSE=$(curl -s -X GET "$BASE_URL/candidates/election/$ELECTION_ID")
    CANDIDATE_ID=$(echo $CANDIDATES_RESPONSE | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
    if [ -n "$CANDIDATE_ID" ]; then
        success "Candidate ID: $CANDIDATE_ID"
    fi
    
    test_endpoint "GET" "$BASE_URL/candidates" "" 200 "Get All Candidates"
    test_endpoint "GET" "$BASE_URL/candidates/$CANDIDATE_ID" "" 200 "Get Candidate by ID"
    test_endpoint "GET" "$BASE_URL/candidates/election/$ELECTION_ID" "" 200 "Get Candidates by Election"
    test_endpoint "GET" "$BASE_URL/candidates/election/$ELECTION_ID/order" "" 200 "Get Candidates by Order"
    test_endpoint "PUT" "$BASE_URL/candidates/$CANDIDATE_ID" '{"name":"Updated Candidate"}' 200 "Update Candidate"
    
    # 7. Ballots endpoints
    log "=== TESTING BALLOTS ENDPOINTS ==="
    # Actualizar elections y voter en los datos de ballot
    sed -i "s/\"elections\": \"[^\"]*\"/\"elections\": \"$ELECTION_ID\"/" $TEST_DATA_DIR/ballot.json
    sed -i "s/\"voter\": [0-9]*/\"voter\": $VOTER_ID/" $TEST_DATA_DIR/ballot.json
    
    test_endpoint "POST" "$BASE_URL/ballots" "$(cat $TEST_DATA_DIR/ballot.json)" 201 "Create Ballot"
    
    test_endpoint "GET" "$BASE_URL/ballots/election/$ELECTION_ID" "" 200 "Get Ballots by Election"
    test_endpoint "GET" "$BASE_URL/ballots/voter/$VOTER_ID" "" 200 "Get Ballots by Voter"
    test_endpoint "GET" "$BASE_URL/ballots/election/$ELECTION_ID/with-details" "" 200 "Get Ballots with Voter Details"
    test_endpoint "GET" "$BASE_URL/ballots/election/$ELECTION_ID/voter/$VOTER_ID/id/$BALLOT_ID" "" 200 "Get Ballot by ID"
    test_endpoint "PUT" "$BASE_URL/ballots/election/$ELECTION_ID/voter/$VOTER_ID/id/$BALLOT_ID" '{"vote":{"updated":true}}' 200 "Update Ballot"
    
    # 8. Homomorphic Keys endpoints
    log "=== TESTING HOMOMORPHIC KEYS ENDPOINTS ==="
    # Actualizar elections en los datos de homomorphic key
    sed -i "s/\"elections\": \"[^\"]*\"/\"elections\": \"$ELECTION_ID\"/" $TEST_DATA_DIR/homomorphic_key.json
    
    test_endpoint "POST" "$BASE_URL/keys" "$(cat $TEST_DATA_DIR/homomorphic_key.json)" 201 "Create Homomorphic Key"
    
    test_endpoint "GET" "$BASE_URL/keys/election/$ELECTION_ID" "" 200 "Get Key by Election"
    test_endpoint "PUT" "$BASE_URL/keys/election/$ELECTION_ID/params" '{"params":{"updated":true}}' 200 "Update Key Params"
    
    # 9. Tally Results endpoints
    log "=== TESTING TALLY RESULTS ENDPOINTS ==="
    # Actualizar election en los datos de tally result
    sed -i "s/\"election\": \"[^\"]*\"/\"election\": \"$ELECTION_ID\"/" $TEST_DATA_DIR/tally_result.json
    
    test_endpoint "POST" "$BASE_URL/tally-results" "$(cat $TEST_DATA_DIR/tally_result.json)" 201 "Create Tally Result"
    
    test_endpoint "GET" "$BASE_URL/tally-results/election/$ELECTION_ID" "" 200 "Get Tally Result by Election"
    test_endpoint "GET" "$BASE_URL/tally-results/with-details" "" 200 "Get Tally Results with Election Details"
    test_endpoint "POST" "$BASE_URL/tally-results/election/$ELECTION_ID/compute" '{"computed_by":"test-script"}' 200 "Compute Tally Result"
    
    # 10. Audit Config endpoints
    log "=== TESTING AUDIT CONFIG ENDPOINTS ==="
    # Actualizar election en los datos de audit config
    sed -i "s/\"election\": \"[^\"]*\"/\"election\": \"$ELECTION_ID\"/" $TEST_DATA_DIR/audit_config.json
    
    test_endpoint "POST" "$BASE_URL/audit-config" "$(cat $TEST_DATA_DIR/audit_config.json)" 201 "Create Audit Config"
    
    test_endpoint "GET" "$BASE_URL/audit-config/election/$ELECTION_ID" "" 200 "Get Audit Config by Election"
    test_endpoint "PUT" "$BASE_URL/audit-config/election/$ELECTION_ID/ballot-audit" '{"enable":false}' 200 "Disable Ballot Audit"
    test_endpoint "PUT" "$BASE_URL/audit-config/election/$ELECTION_ID/access-logs" '{"enable":false}' 200 "Disable Access Logs"
    
    # 11. Audit Logs endpoints
    log "=== TESTING AUDIT LOGS ENDPOINTS ==="
    # Actualizar election en los datos de audit log
    sed -i "s/\"election\": \"[^\"]*\"/\"election\": \"$ELECTION_ID\"/" $TEST_DATA_DIR/audit_log.json
    
    test_endpoint "POST" "$BASE_URL/audit-logs" "$(cat $TEST_DATA_DIR/audit_log.json)" 201 "Create Audit Log"
    
    test_endpoint "GET" "$BASE_URL/audit-logs/election/$ELECTION_ID" "" 200 "Get Audit Logs by Election"
    test_endpoint "POST" "$BASE_URL/audit-logs/election/$ELECTION_ID/vote" '{"voter_id":"test_voter","ip_address":"192.168.1.100","user_agent":"Test-Agent"}' 200 "Log Vote Action"
    test_endpoint "POST" "$BASE_URL/audit-logs/election/$ELECTION_ID/authority" '{"authority_id":"test_auth","action":"test_action","details":{"test":true}}' 200 "Log Authority Action"
    
    # Limpieza (opcional - comentar si quieres mantener los datos)
    log "=== CLEANUP ==="
    test_endpoint "DELETE" "$BASE_URL/ballots/election/$ELECTION_ID/voter/$VOTER_ID/id/$BALLOT_ID" "" 204 "Delete Ballot"
    test_endpoint "DELETE" "$BASE_URL/candidates/$CANDIDATE_ID" "" 204 "Delete Candidate"
    test_endpoint "DELETE" "$BASE_URL/voters/$VOTER_ID" "" 204 "Delete Voter"
    test_endpoint "DELETE" "$BASE_URL/authorities/$AUTHORITY_ID" "" 204 "Delete Authority"
    test_endpoint "DELETE" "$BASE_URL/elections/$ELECTION_ID" "" 204 "Delete Election"
    test_endpoint "DELETE" "$BASE_URL/status/999" "" 204 "Delete Status"
    
    log "=== TEST COMPLETED ==="
}

# Verificar dependencias
check_dependencies() {
    if ! command -v curl &> /dev/null; then
        error "curl is required but not installed. Please install curl."
        exit 1
    fi
    
    if ! command -v jq &> /dev/null; then
        warning "jq is not installed. Some features might not work optimally."
    fi
}

# Función para probar endpoints específicos
test_specific_endpoints() {
    local endpoints=(
        "GET:/health"
        "GET:/elections"
        "GET:/status"
        "GET:/authorities"
        "GET:/voters"
        "GET:/candidates"
        "GET:/keys"
        "GET:/tally-results"
        "GET:/audit-config"
    )
    
    log "Testing specific endpoints..."
    for endpoint in "${endpoints[@]}"; do
        IFS=':' read -r method path <<< "$endpoint"
        test_endpoint "$method" "$BASE_URL$path" "" 200 "Specific: $method $path"
    done
}

# Menú principal
case "${1:-}" in
    "specific")
        check_dependencies
        test_specific_endpoints
        ;;
    "cleanup")
        # Función de limpieza específica
        log "Running cleanup..."
        ;;
    *)
        check_dependencies
        run_tests
        ;;
esac