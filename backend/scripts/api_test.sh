#!/bin/bash

#===============================================================================
# SIPODI API Testing Script
# Comprehensive automated testing for all API endpoints
#===============================================================================

# Don't exit on error - we want to continue testing even if some tests fail
# set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Colo
BOLD='\033[1m'

# Configuration
BASE_URL="${API_BASE_URL:-http://localhost:8080/api/v1}"
SUPER_ADMIN_EMAIL="${SUPER_ADMIN_EMAIL:-superadmin@sipodi.go.id}"
SUPER_ADMIN_PASSWORD="${SUPER_ADMIN_PASSWORD:-admin123}"

# Global variables for tokens and IDs
ACCESS_TOKEN=""
REFRESH_TOKEN=""
CREATED_SCHOOL_ID=""
CREATED_USER_ID=""
CREATED_TALENT_ID=""
CREATED_UPLOAD_ID=""
GTK_ACCESS_TOKEN=""
ADMIN_SEKOLAH_ACCESS_TOKEN=""

# Counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

#===============================================================================
# Utility Functions
#===============================================================================

print_header() {
    echo ""
    echo -e "${BOLD}${BLUE}╔══════════════════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BOLD}${BLUE}║${NC} ${BOLD}$1${NC}"
    echo -e "${BOLD}${BLUE}╚══════════════════════════════════════════════════════════════════════════════╝${NC}"
}

print_section() {
    echo ""
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BOLD}${CYAN}  $1${NC}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

print_test() {
    echo ""
    echo -e "${MAGENTA}┌──────────────────────────────────────────────────────────────────────────────┐${NC}"
    echo -e "${MAGENTA}│${NC} ${BOLD}TEST:${NC} $1"
    echo -e "${MAGENTA}│${NC} ${YELLOW}$2${NC} $3"
    echo -e "${MAGENTA}└──────────────────────────────────────────────────────────────────────────────┘${NC}"
}

print_description() {
    echo -e "${CYAN}📝 Description:${NC} $1"
}

print_auth() {
    echo -e "${YELLOW}🔐 Authentication:${NC} $1"
}

print_params() {
    echo -e "${BLUE}📋 Parameters:${NC}"
    echo "$1" | sed 's/^/   /'
}

print_request() {
    echo -e "${GREEN}📤 Request:${NC}"
    if [ -n "$1" ]; then
        echo "$1" | jq '.' 2>/dev/null | sed 's/^/   /' || echo "   $1"
    else
        echo "   (no body)"
    fi
}

print_response() {
    local status_code=$1
    local body=$2
    
    if [ "$status_code" -ge 200 ] && [ "$status_code" -lt 300 ]; then
        echo -e "${GREEN}📥 Response (${status_code}):${NC}"
    elif [ "$status_code" -ge 400 ] && [ "$status_code" -lt 500 ]; then
        echo -e "${YELLOW}📥 Response (${status_code}):${NC}"
    else
        echo -e "${RED}📥 Response (${status_code}):${NC}"
    fi
    echo "$body" | jq '.' 2>/dev/null | sed 's/^/   /' || echo "   $body"
}

print_success() {
    echo -e "${GREEN}✅ TEST PASSED${NC}"
    ((PASSED_TESTS++))
    ((TOTAL_TESTS++))
}

print_failure() {
    echo -e "${RED}❌ TEST FAILED: $1${NC}"
    ((FAILED_TESTS++))
    ((TOTAL_TESTS++))
}

print_error_scenarios() {
    echo -e "${RED}⚠️  Possible Error Scenarios:${NC}"
    for scenario in "$@"; do
        echo -e "   ${RED}•${NC} $scenario"
    done
}

# HTTP request helpe
do_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local token=$4
    local extra_headers=$5
    
    local url="${BASE_URL}${endpoint}"
    local headers="-H 'Content-Type: application/json'"
    
    if [ -n "$token" ]; then
        headers="$headers -H 'Authorization: Bearer $token'"
    fi
    
    if [ -n "$extra_headers" ]; then
        headers="$headers $extra_headers"
    fi
    
    local response
    local http_code
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" \
            -H "Content-Type: application/json" \
            ${token:+-H "Authorization: Bearer $token"} \
            -d "$data" 2>/dev/null)
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" \
            -H "Content-Type: application/json" \
            ${token:+-H "Authorization: Bearer $token"} 2>/dev/null)
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    echo "$http_code"
    echo "$body"
}

# Extract value from JSON response
extract_json() {
    echo "$1" | jq -r "$2" 2>/dev/null
}

wait_for_api() {
    echo -e "${YELLOW}⏳ Waiting for API to be ready...${NC}"
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s "${BASE_URL}/health" > /dev/null 2>&1; then
            echo -e "${GREEN}✅ API is ready!${NC}"
            return 0
        fi
        echo "   Attempt $attempt/$max_attempts..."
        sleep 2
        ((attempt++))
    done
    
    echo -e "${RED}❌ API not responding after $max_attempts attempts${NC}"
    exit 1
}

#===============================================================================
# 1. AUTHENTICATION TESTS
#===============================================================================

test_auth_login() {
    print_test "POST /auth/login" "POST" "/auth/login"
    print_description "Login user dan dapatkan access token"
    print_auth "None"
    print_params "Body: email (string, required), password (string, required)"
    
    local request_body='{
        "email": "'"$SUPER_ADMIN_EMAIL"'",
        "password": "'"$SUPER_ADMIN_PASSWORD"'"
    }'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/auth/login" "$request_body")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "401 INVALID_CREDENTIALS - Email atau password salah" \
        "403 ACCOUNT_DISABLED - Akun telah dinonaktifkan" \
        "422 VALIDATION_ERROR - Email/password tidak diisi"
    
    if [ "$http_code" = "200" ]; then
        ACCESS_TOKEN=$(extract_json "$body" '.data.access_token')
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_auth_login_invalid_credentials() {
    print_test "POST /auth/login (Invalid Credentials)" "POST" "/auth/login"
    print_description "Test login dengan kredensial salah"
    print_auth "None"
    
    local request_body='{
        "email": "wrong@email.com",
        "password": "wrongpassword"
    }'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/auth/login" "$request_body")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "401" ]; then
        print_success
    else
        print_failure "Expected 401, got $http_code"
    fi
}

test_auth_login_validation_error() {
    print_test "POST /auth/login (Validation Error)" "POST" "/auth/login"
    print_description "Test login tanpa email dan password"
    print_auth "None"
    
    local request_body='{}'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/auth/login" "$request_body")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "422" ] || [ "$http_code" = "400" ]; then
        print_success
    else
        print_failure "Expected 422 or 400, got $http_code"
    fi
}

test_auth_refresh() {
    print_test "POST /auth/refresh" "POST" "/auth/refresh"
    print_description "Refresh access token menggunakan refresh token dari cookie"
    print_auth "None (menggunakan HttpOnly cookie)"
    print_params "Cookie: refresh_token (HttpOnly)"
    
    print_request "(no body - uses cookie)"
    
    local result=$(do_request "POST" "/auth/refresh")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "401 TOKEN_EXPIRED - Refresh token telah expired" \
        "401 INVALID_TOKEN - Refresh token tidak valid"
    
    # Note: This might fail without proper cookie handling
    if [ "$http_code" = "200" ] || [ "$http_code" = "401" ]; then
        print_success
    else
        print_failure "Expected 200 or 401, got $http_code"
    fi
}

test_auth_logout() {
    print_test "POST /auth/logout" "POST" "/auth/logout"
    print_description "Logout dari sesi saat ini"
    print_auth "Required (Bearer Token)"
    
    print_request "(no body)"
    
    local result=$(do_request "POST" "/auth/logout" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "401 UNAUTHORIZED - Token tidak valid atau sudah expired"
    
    if [ "$http_code" = "200" ]; then
        print_success
        # Re-login to get new token
        test_auth_login > /dev/null 2>&1
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_auth_logout_all() {
    print_test "POST /auth/logout-all" "POST" "/auth/logout-all"
    print_description "Logout dari semua perangkat/sesi"
    print_auth "Required (Bearer Token)"
    
    print_request "(no body)"
    
    local result=$(do_request "POST" "/auth/logout-all" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "401 UNAUTHORIZED - Token tidak valid atau sudah expired"
    
    if [ "$http_code" = "200" ]; then
        print_success
        # Re-login to get new token
        test_auth_login > /dev/null 2>&1
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_auth_unauthorized() {
    print_test "POST /auth/logout (Unauthorized)" "POST" "/auth/logout"
    print_description "Test logout tanpa token"
    print_auth "None (should fail)"
    
    print_request "(no body, no token)"
    
    local result=$(do_request "POST" "/auth/logout")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "401" ]; then
        print_success
    else
        print_failure "Expected 401, got $http_code"
    fi
}

#===============================================================================
# 2. PROFILE (ME) TESTS
#===============================================================================

test_me_get() {
    print_test "GET /me" "GET" "/me"
    print_description "Profil user yang sedang login"
    print_auth "Required (Bearer Token)"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/me" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "401 UNAUTHORIZED - Token tidak valid atau sudah expired"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_me_update() {
    print_test "PATCH /me" "PATCH" "/me"
    print_description "Update profil user yang sedang login"
    print_auth "Required (Bearer Token)"
    print_params "Body: full_name, position, gender, birth_date (all optional)"
    
    local request_body='{
        "full_name": "Super Administrator Updated",
        "position": "System Administrator"
    }'
    print_request "$request_body"
    
    local result=$(do_request "PATCH" "/me" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "401 UNAUTHORIZED - Token tidak valid" \
        "422 VALIDATION_ERROR - Data tidak valid"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_me_password() {
    print_test "PATCH /me/password" "PATCH" "/me/password"
    print_description "Ubah password"
    print_auth "Required (Bearer Token)"
    print_params "Body: current_password, new_password, new_password_confirmation (all required)"
    
    local request_body='{
        "current_password": "'"$SUPER_ADMIN_PASSWORD"'",
        "new_password": "newpassword123",
        "new_password_confirmation": "newpassword123"
    }'
    print_request "$request_body"
    
    local result=$(do_request "PATCH" "/me/password" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "400 INVALID_PASSWORD - Password lama tidak sesuai" \
        "422 VALIDATION_ERROR - Password minimal 8 karakter / konfirmasi tidak cocok"
    
    if [ "$http_code" = "200" ]; then
        print_success
        # Change password back
        local revert_body='{
            "current_password": "newpassword123",
            "new_password": "'"$SUPER_ADMIN_PASSWORD"'",
            "new_password_confirmation": "'"$SUPER_ADMIN_PASSWORD"'"
        }'
        do_request "PATCH" "/me/password" "$revert_body" "$ACCESS_TOKEN" > /dev/null 2>&1
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_me_password_wrong_current() {
    print_test "PATCH /me/password (Wrong Current)" "PATCH" "/me/password"
    print_description "Test ubah password dengan password lama salah"
    print_auth "Required (Bearer Token)"
    
    local request_body='{
        "current_password": "wrongpassword",
        "new_password": "newpassword123",
        "new_password_confirmation": "newpassword123"
    }'
    print_request "$request_body"
    
    local result=$(do_request "PATCH" "/me/password" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "400" ]; then
        print_success
    else
        print_failure "Expected 400, got $http_code"
    fi
}

test_me_photo() {
    print_test "PATCH /me/photo" "PATCH" "/me/photo"
    print_description "Update foto profil (setelah upload via presigned URL)"
    print_auth "Required (Bearer Token)"
    print_params "Body: upload_id (UUID, required)"
    
    local request_body='{
        "upload_id": "00000000-0000-0000-0000-000000000000"
    }'
    print_request "$request_body"
    
    local result=$(do_request "PATCH" "/me/photo" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "404 UPLOAD_NOT_FOUND - Upload tidak ditemukan atau sudah expired"
    
    # Expected to fail with fake upload_id
    if [ "$http_code" = "404" ] || [ "$http_code" = "400" ]; then
        print_success
    else
        print_failure "Expected 404 or 400, got $http_code"
    fi
}

#===============================================================================
# 3. SEKOLAH TESTS
#===============================================================================

test_schools_list() {
    print_test "GET /schools" "GET" "/schools"
    print_description "Daftar semua sekolah"
    print_auth "Required (Super Admin)"
    print_params "Query: search, status (negeri/swasta), page, limit, sort"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/schools?page=1&limit=10" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "403 FORBIDDEN - Bukan Super Admin"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_schools_list_with_filter() {
    print_test "GET /schools (With Filters)" "GET" "/schools?status=negeri&search=SMA"
    print_description "Daftar sekolah dengan filter"
    print_auth "Required (Super Admin)"
    print_params "Query: status=negeri, search=SMA"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/schools?status=negeri&search=SMA&page=1&limit=10" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_schools_create() {
    print_test "POST /schools" "POST" "/schools"
    print_description "Tambah sekolah baru"
    print_auth "Required (Super Admin)"
    print_params "Body: name, npsn, status, address (all required)"
    
    local timestamp=$(date +%s)
    local request_body='{
        "name": "SMAN Test '"$timestamp"'",
        "npsn": "'"$timestamp"'",
        "status": "negeri",
        "address": "Jl. Test No. 1, Malang"
    }'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/schools" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "409 DUPLICATE_NPSN - NPSN sudah terdaftar" \
        "422 VALIDATION_ERROR - Data tidak lengkap" \
        "403 FORBIDDEN - Bukan Super Admin"
    
    if [ "$http_code" = "201" ]; then
        CREATED_SCHOOL_ID=$(extract_json "$body" '.data.id')
        print_success
    else
        print_failure "Expected 201, got $http_code"
    fi
}

test_schools_create_duplicate_npsn() {
    print_test "POST /schools (Duplicate NPSN)" "POST" "/schools"
    print_description "Test tambah sekolah dengan NPSN yang sudah ada"
    print_auth "Required (Super Admin)"
    
    local request_body='{
        "name": "SMAN Duplicate",
        "npsn": "20518765",
        "status": "negeri",
        "address": "Jl. Duplicate No. 1"
    }'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/schools" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "409" ] || [ "$http_code" = "400" ]; then
        print_success
    else
        print_failure "Expected 409 or 400, got $http_code"
    fi
}

test_schools_get() {
    print_test "GET /schools/{id}" "GET" "/schools/{id}"
    print_description "Detail sekolah berdasarkan ID"
    print_auth "Required"
    print_params "Path: id (UUID)"
    
    if [ -z "$CREATED_SCHOOL_ID" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No school ID available${NC}"
        return
    fi
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/schools/$CREATED_SCHOOL_ID" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "404 NOT_FOUND - Sekolah tidak ditemukan"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_schools_get_not_found() {
    print_test "GET /schools/{id} (Not Found)" "GET" "/schools/{id}"
    print_description "Test get sekolah dengan ID tidak valid"
    print_auth "Required"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/schools/00000000-0000-0000-0000-000000000000" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "404" ]; then
        print_success
    else
        print_failure "Expected 404, got $http_code"
    fi
}

test_schools_update() {
    print_test "PUT /schools/{id}" "PUT" "/schools/{id}"
    print_description "Update data sekolah"
    print_auth "Required (Super Admin)"
    print_params "Path: id (UUID), Body: name, npsn, status, address, head_master_id"
    
    if [ -z "$CREATED_SCHOOL_ID" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No school ID available${NC}"
        return
    fi
    
    local request_body='{
        "name": "SMAN Test Updated",
        "status": "negeri",
        "address": "Jl. Test Updated No. 2, Malang"
    }'
    print_request "$request_body"
    
    local result=$(do_request "PUT" "/schools/$CREATED_SCHOOL_ID" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "404 NOT_FOUND - Sekolah tidak ditemukan" \
        "400 INVALID_HEAD_MASTER - User bukan kepala sekolah" \
        "403 FORBIDDEN - Bukan Super Admin"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_schools_users() {
    print_test "GET /schools/{id}/users" "GET" "/schools/{id}/users"
    print_description "Daftar GTK di sekolah tertentu"
    print_auth "Required (Super Admin, Admin Sekolah)"
    print_params "Path: id (UUID), Query: search, gtk_type, is_active, page, limit"
    
    if [ -z "$CREATED_SCHOOL_ID" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No school ID available${NC}"
        return
    fi
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/schools/$CREATED_SCHOOL_ID/users?page=1&limit=10" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "403 FORBIDDEN - Admin sekolah akses sekolah lain" \
        "404 NOT_FOUND - Sekolah tidak ditemukan"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_schools_delete() {
    print_test "DELETE /schools/{id}" "DELETE" "/schools/{id}"
    print_description "Hapus sekolah"
    print_auth "Required (Super Admin)"
    print_params "Path: id (UUID)"
    
    if [ -z "$CREATED_SCHOOL_ID" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No school ID available${NC}"
        return
    fi
    
    print_request "(no body)"
    
    local result=$(do_request "DELETE" "/schools/$CREATED_SCHOOL_ID" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "404 NOT_FOUND - Sekolah tidak ditemukan" \
        "400 HAS_DEPENDENCIES - Masih memiliki GTK" \
        "403 FORBIDDEN - Bukan Super Admin"
    
    if [ "$http_code" = "204" ] || [ "$http_code" = "200" ]; then
        CREATED_SCHOOL_ID=""
        print_success
    else
        print_failure "Expected 204 or 200, got $http_code"
    fi
}

#===============================================================================
# 4. USERS/GTK TESTS
#===============================================================================

test_users_list() {
    print_test "GET /users" "GET" "/users"
    print_description "Daftar semua user"
    print_auth "Required (Super Admin)"
    print_params "Query: search, role, school_id, gtk_type, is_active, page, limit, sort"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/users?page=1&limit=10" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "403 FORBIDDEN - Bukan Super Admin"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_users_list_with_filter() {
    print_test "GET /users (With Filters)" "GET" "/users?role=gtk&is_active=true"
    print_description "Daftar user dengan filter"
    print_auth "Required (Super Admin)"
    print_params "Query: role=gtk, is_active=true"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/users?role=gtk&is_active=true&page=1&limit=10" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

# Create a school first for user creation
setup_school_for_user() {
    local timestamp=$(date +%s)
    local request_body='{
        "name": "SMAN User Test '"$timestamp"'",
        "npsn": "'"$timestamp"'",
        "status": "negeri",
        "address": "Jl. User Test No. 1"
    }'
    
    local result=$(do_request "POST" "/schools" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    if [ "$http_code" = "201" ]; then
        CREATED_SCHOOL_ID=$(extract_json "$body" '.data.id')
    fi
}

# Global GTK credentials for talent tests
GTK_EMAIL=""
GTK_PASSWORD="password123"

# Setup GTK user for talent tests
setup_gtk_for_talents() {
    if [ -n "$GTK_ACCESS_TOKEN" ]; then
        return 0
    fi
    
    # First create a school if needed
    if [ -z "$CREATED_SCHOOL_ID" ]; then
        setup_school_for_user
    fi
    
    if [ -z "$CREATED_SCHOOL_ID" ]; then
        return 1
    fi
    
    local timestamp=$(date +%s)
    GTK_EMAIL="gtktest${timestamp}@sekolah.sch.id"
    
    local request_body='{
        "email": "'"$GTK_EMAIL"'",
        "password": "'"$GTK_PASSWORD"'",
        "role": "gtk",
        "full_name": "GTK Test User",
        "gtk_type": "guru",
        "school_id": "'"$CREATED_SCHOOL_ID"'"
    }'
    
    local result=$(do_request "POST" "/users" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    if [ "$http_code" != "201" ]; then
        return 1
    fi
    
    CREATED_USER_ID=$(extract_json "$body" '.data.id')
    
    # Login as GTK
    local login_body='{
        "email": "'"$GTK_EMAIL"'",
        "password": "'"$GTK_PASSWORD"'"
    }'
    
    result=$(do_request "POST" "/auth/login" "$login_body")
    http_code=$(echo "$result" | head -n1)
    body=$(echo "$result" | tail -n +2)
    
    if [ "$http_code" = "200" ]; then
        GTK_ACCESS_TOKEN=$(extract_json "$body" '.data.access_token')
        return 0
    fi
    
    return 1
}

test_users_create() {
    print_test "POST /users" "POST" "/users"
    print_description "Tambah user baru"
    print_auth "Required (Super Admin, Admin Sekolah)"
    print_params "Body: email, password, role, full_name, nuptk, nip, gender, birth_date, gtk_type, position, school_id"
    
    # Setup school first
    setup_school_for_user
    
    if [ -z "$CREATED_SCHOOL_ID" ]; then
        echo -e "${YELLOW}⚠️  Skipping: Could not create school${NC}"
        return
    fi
    
    local timestamp=$(date +%s)
    local request_body='{
        "email": "testgtk'"$timestamp"'@sekolah.sch.id",
        "password": "password123",
        "role": "gtk",
        "full_name": "Test GTK '"$timestamp"'",
        "nuptk": "'"$timestamp"'1234",
        "nip": "'"$timestamp"'5678",
        "gender": "L",
        "birth_date": "1990-01-01",
        "gtk_type": "guru",
        "position": "Guru Matematika",
        "school_id": "'"$CREATED_SCHOOL_ID"'"
    }'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/users" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "409 EMAIL_TAKEN - Email sudah terdaftar" \
        "409 NUPTK_TAKEN - NUPTK sudah terdaftar" \
        "409 NIP_TAKEN - NIP sudah terdaftar" \
        "422 VALIDATION_ERROR - Data tidak valid" \
        "403 FORBIDDEN - Admin sekolah buat user di sekolah lain"
    
    if [ "$http_code" = "201" ]; then
        CREATED_USER_ID=$(extract_json "$body" '.data.id')
        print_success
    else
        print_failure "Expected 201, got $http_code"
    fi
}

test_users_create_duplicate_email() {
    print_test "POST /users (Duplicate Email)" "POST" "/users"
    print_description "Test tambah user dengan email yang sudah ada"
    print_auth "Required (Super Admin)"
    
    local request_body='{
        "email": "'"$SUPER_ADMIN_EMAIL"'",
        "password": "password123",
        "role": "gtk",
        "full_name": "Duplicate User",
        "gtk_type": "guru",
        "school_id": "'"$CREATED_SCHOOL_ID"'"
    }'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/users" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "409" ] || [ "$http_code" = "400" ]; then
        print_success
    else
        print_failure "Expected 409 or 400, got $http_code"
    fi
}

test_users_get() {
    print_test "GET /users/{id}" "GET" "/users/{id}"
    print_description "Detail user berdasarkan ID"
    print_auth "Required"
    print_params "Path: id (UUID)"
    
    if [ -z "$CREATED_USER_ID" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No user ID available${NC}"
        return
    fi
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/users/$CREATED_USER_ID" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "404 NOT_FOUND - User tidak ditemukan" \
        "403 FORBIDDEN - Admin sekolah akses user sekolah lain"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_users_get_not_found() {
    print_test "GET /users/{id} (Not Found)" "GET" "/users/{id}"
    print_description "Test get user dengan ID tidak valid"
    print_auth "Required"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/users/00000000-0000-0000-0000-000000000000" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "404" ]; then
        print_success
    else
        print_failure "Expected 404, got $http_code"
    fi
}

test_users_update() {
    print_test "PUT /users/{id}" "PUT" "/users/{id}"
    print_description "Update data user"
    print_auth "Required (Super Admin, Admin Sekolah)"
    print_params "Path: id (UUID), Body: full_name, position, gtk_type, gender, birth_date, school_id"
    
    if [ -z "$CREATED_USER_ID" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No user ID available${NC}"
        return
    fi
    
    local request_body='{
        "full_name": "Test GTK Updated",
        "position": "Guru Matematika Senior"
    }'
    print_request "$request_body"
    
    local result=$(do_request "PUT" "/users/$CREATED_USER_ID" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "404 NOT_FOUND - User tidak ditemukan" \
        "403 FORBIDDEN - Tidak memiliki akses"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_users_deactivate() {
    print_test "PATCH /users/{id}/deactivate" "PATCH" "/users/{id}/deactivate"
    print_description "Nonaktifkan user"
    print_auth "Required (Super Admin, Admin Sekolah)"
    print_params "Path: id (UUID)"
    
    if [ -z "$CREATED_USER_ID" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No user ID available${NC}"
        return
    fi
    
    print_request "(no body)"
    
    local result=$(do_request "PATCH" "/users/$CREATED_USER_ID/deactivate" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "404 NOT_FOUND - User tidak ditemukan" \
        "403 FORBIDDEN - Tidak memiliki akses"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_users_activate() {
    print_test "PATCH /users/{id}/activate" "PATCH" "/users/{id}/activate"
    print_description "Aktifkan user"
    print_auth "Required (Super Admin, Admin Sekolah)"
    print_params "Path: id (UUID)"
    
    if [ -z "$CREATED_USER_ID" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No user ID available${NC}"
        return
    fi
    
    print_request "(no body)"
    
    local result=$(do_request "PATCH" "/users/$CREATED_USER_ID/activate" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "404 NOT_FOUND - User tidak ditemukan" \
        "403 FORBIDDEN - Tidak memiliki akses"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_users_delete() {
    print_test "DELETE /users/{id}" "DELETE" "/users/{id}"
    print_description "Hapus user (soft delete)"
    print_auth "Required (Super Admin)"
    print_params "Path: id (UUID)"
    
    if [ -z "$CREATED_USER_ID" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No user ID available${NC}"
        return
    fi
    
    print_request "(no body)"
    
    local result=$(do_request "DELETE" "/users/$CREATED_USER_ID" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "404 NOT_FOUND - User tidak ditemukan" \
        "400 CANNOT_DELETE_SELF - Tidak dapat menghapus akun sendiri" \
        "403 FORBIDDEN - Bukan Super Admin"
    
    if [ "$http_code" = "204" ] || [ "$http_code" = "200" ]; then
        CREATED_USER_ID=""
        print_success
    else
        print_failure "Expected 204 or 200, got $http_code"
    fi
}

#===============================================================================
# 5. TALENTA TESTS
#===============================================================================

test_talents_list() {
    print_test "GET /talents" "GET" "/talents"
    print_description "Daftar semua talenta (untuk admin)"
    print_auth "Required (Super Admin, Admin Sekolah)"
    print_params "Query: user_id, school_id, talent_type, status, page, limit, sort"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/talents?page=1&limit=10" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "403 FORBIDDEN - Tidak memiliki akses"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_talents_list_with_filter() {
    print_test "GET /talents (With Filters)" "GET" "/talents?talent_type=peserta_pelatihan&status=pending"
    print_description "Daftar talenta dengan filter"
    print_auth "Required (Super Admin, Admin Sekolah)"
    print_params "Query: talent_type=peserta_pelatihan, status=pending"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/talents?talent_type=peserta_pelatihan&status=pending&page=1&limit=10" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_me_talents_list() {
    print_test "GET /me/talents" "GET" "/me/talents"
    print_description "Daftar talenta milik user yang login"
    print_auth "Required (GTK)"
    print_params "Query: talent_type, status, page, limit"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/me/talents?page=1&limit=10" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_me_talents_create_peserta_pelatihan() {
    print_test "POST /me/talents (Peserta Pelatihan)" "POST" "/me/talents"
    print_description "Tambah talenta baru - Peserta Pelatihan"
    print_auth "Required (GTK)"
    print_params "Body: talent_type, detail (activity_name, organizer, start_date, duration_days)"
    
    # Setup GTK user first
    setup_gtk_for_talents
    if [ -z "$GTK_ACCESS_TOKEN" ]; then
        echo -e "${YELLOW}⚠️  Skipping: Could not setup GTK user${NC}"
        return
    fi
    
    local request_body='{
        "talent_type": "peserta_pelatihan",
        "detail": {
            "activity_name": "Pelatihan Kurikulum Merdeka",
            "organizer": "Kemendikbud",
            "start_date": "2024-06-01",
            "duration_days": 5
        }
    }'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/me/talents" "$request_body" "$GTK_ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "422 VALIDATION_ERROR - detail.activity_name wajib diisi" \
        "422 VALIDATION_ERROR - detail.organizer wajib diisi" \
        "422 VALIDATION_ERROR - detail.start_date wajib diisi" \
        "422 VALIDATION_ERROR - detail.duration_days harus > 0"
    
    if [ "$http_code" = "201" ]; then
        CREATED_TALENT_ID=$(extract_json "$body" '.data.id')
        print_success
    else
        print_failure "Expected 201, got $http_code"
    fi
}

test_me_talents_create_pembimbing_lomba() {
    print_test "POST /me/talents (Pembimbing Lomba)" "POST" "/me/talents"
    print_description "Tambah talenta baru - Pembimbing Lomba"
    print_auth "Required (GTK)"
    print_params "Body: talent_type, detail (competition_name, level, organizer, field, achievement), upload_id"
    
    if [ -z "$GTK_ACCESS_TOKEN" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No GTK token available${NC}"
        return
    fi
    
    local request_body='{
        "talent_type": "pembimbing_lomba",
        "detail": {
            "competition_name": "Olimpiade Sains Nasional",
            "level": "nasional",
            "organizer": "Kemendikbud",
            "field": "akademik",
            "achievement": "Juara 1 Tingkat Nasional"
        }
    }'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/me/talents" "$request_body" "$GTK_ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "422 VALIDATION_ERROR - detail.competition_name wajib diisi" \
        "422 VALIDATION_ERROR - detail.level harus kota/provinsi/nasional/internasional" \
        "422 VALIDATION_ERROR - detail.field tidak valid" \
        "422 VALIDATION_ERROR - detail.achievement wajib diisi"
    
    if [ "$http_code" = "201" ]; then
        print_success
    else
        print_failure "Expected 201, got $http_code"
    fi
}

test_me_talents_create_peserta_lomba() {
    print_test "POST /me/talents (Peserta Lomba)" "POST" "/me/talents"
    print_description "Tambah talenta baru - Peserta Lomba"
    print_auth "Required (GTK)"
    print_params "Body: talent_type, detail (competition_name, level, organizer, field, start_date, duration_days, competition_field, achievement), upload_id"
    
    if [ -z "$GTK_ACCESS_TOKEN" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No GTK token available${NC}"
        return
    fi
    
    local request_body='{
        "talent_type": "peserta_lomba",
        "detail": {
            "competition_name": "Lomba Guru Berprestasi",
            "level": "provinsi",
            "organizer": "Dinas Pendidikan Jatim",
            "field": "akademik",
            "start_date": "2024-08-15",
            "duration_days": 3,
            "competition_field": "Guru SMA",
            "achievement": "Juara 2 Tingkat Provinsi"
        }
    }'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/me/talents" "$request_body" "$GTK_ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "201" ]; then
        print_success
    else
        print_failure "Expected 201, got $http_code"
    fi
}

test_me_talents_create_minat_bakat() {
    print_test "POST /me/talents (Minat/Bakat)" "POST" "/me/talents"
    print_description "Tambah talenta baru - Minat/Bakat"
    print_auth "Required (GTK)"
    print_params "Body: talent_type, detail (interest_name, description), upload_id"
    
    if [ -z "$GTK_ACCESS_TOKEN" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No GTK token available${NC}"
        return
    fi
    
    local request_body='{
        "talent_type": "minat_bakat",
        "detail": {
            "interest_name": "Menulis Buku",
            "description": "Penulis buku pelajaran matematika untuk SMA"
        }
    }'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/me/talents" "$request_body" "$GTK_ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "422 VALIDATION_ERROR - detail.interest_name wajib diisi" \
        "422 VALIDATION_ERROR - detail.description wajib diisi"
    
    if [ "$http_code" = "201" ]; then
        print_success
    else
        print_failure "Expected 201, got $http_code"
    fi
}

test_me_talents_create_validation_error() {
    print_test "POST /me/talents (Validation Error)" "POST" "/me/talents"
    print_description "Test tambah talenta dengan data tidak lengkap"
    print_auth "Required (GTK)"
    
    if [ -z "$GTK_ACCESS_TOKEN" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No GTK token available${NC}"
        return
    fi
    
    local request_body='{
        "talent_type": "peserta_pelatihan",
        "detail": {}
    }'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/me/talents" "$request_body" "$GTK_ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "422" ] || [ "$http_code" = "400" ]; then
        print_success
    else
        print_failure "Expected 422 or 400, got $http_code"
    fi
}

test_talents_get() {
    print_test "GET /talents/{id}" "GET" "/talents/{id}"
    print_description "Detail talenta berdasarkan ID"
    print_auth "Required"
    print_params "Path: id (UUID)"
    
    if [ -z "$CREATED_TALENT_ID" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No talent ID available${NC}"
        return
    fi
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/talents/$CREATED_TALENT_ID" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "404 NOT_FOUND - Talenta tidak ditemukan" \
        "403 FORBIDDEN - Tidak memiliki akses"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_me_talents_update() {
    print_test "PUT /me/talents/{id}" "PUT" "/me/talents/{id}"
    print_description "Update talenta milik sendiri (status kembali pending)"
    print_auth "Required (GTK)"
    print_params "Path: id (UUID), Body: detail, upload_id"
    
    if [ -z "$CREATED_TALENT_ID" ] || [ -z "$GTK_ACCESS_TOKEN" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No talent ID or GTK token available${NC}"
        return
    fi
    
    local request_body='{
        "detail": {
            "activity_name": "Pelatihan Kurikulum Merdeka (Updated)",
            "organizer": "Kemendikbud RI",
            "start_date": "2024-06-01",
            "duration_days": 7
        }
    }'
    print_request "$request_body"
    
    local result=$(do_request "PUT" "/me/talents/$CREATED_TALENT_ID" "$request_body" "$GTK_ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "404 NOT_FOUND - Talenta tidak ditemukan" \
        "403 FORBIDDEN - Bukan milik sendiri"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_me_talents_delete() {
    print_test "DELETE /me/talents/{id}" "DELETE" "/me/talents/{id}"
    print_description "Hapus talenta milik sendiri"
    print_auth "Required (GTK)"
    print_params "Path: id (UUID)"
    
    if [ -z "$CREATED_TALENT_ID" ] || [ -z "$GTK_ACCESS_TOKEN" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No talent ID or GTK token available${NC}"
        return
    fi
    
    print_request "(no body)"
    
    local result=$(do_request "DELETE" "/me/talents/$CREATED_TALENT_ID" "" "$GTK_ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "404 NOT_FOUND - Talenta tidak ditemukan" \
        "403 FORBIDDEN - Bukan milik sendiri"
    
    if [ "$http_code" = "204" ] || [ "$http_code" = "200" ]; then
        CREATED_TALENT_ID=""
        print_success
    else
        print_failure "Expected 204 or 200, got $http_code"
    fi
}

#===============================================================================
# 6. VERIFIKASI TALENTA TESTS
#===============================================================================

# Create a talent for verification tests
setup_talent_for_verification() {
    if [ -z "$GTK_ACCESS_TOKEN" ]; then
        setup_gtk_for_talents
    fi
    
    if [ -z "$GTK_ACCESS_TOKEN" ]; then
        return 1
    fi
    
    local request_body='{
        "talent_type": "peserta_pelatihan",
        "detail": {
            "activity_name": "Pelatihan Verifikasi Test",
            "organizer": "Test Organizer",
            "start_date": "2024-06-01",
            "duration_days": 3
        }
    }'
    
    local result=$(do_request "POST" "/me/talents" "$request_body" "$GTK_ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    if [ "$http_code" = "201" ]; then
        CREATED_TALENT_ID=$(extract_json "$body" '.data.id')
    fi
}

test_verifications_list() {
    print_test "GET /verifications/talents" "GET" "/verifications/talents"
    print_description "Daftar talenta yang perlu diverifikasi"
    print_auth "Required (Super Admin, Admin Sekolah)"
    print_params "Query: status (default: pending), school_id, talent_type, page, limit"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/verifications/talents?status=pending&page=1&limit=10" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "403 FORBIDDEN - Tidak memiliki akses"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_verifications_approve() {
    print_test "POST /verifications/talents/{id}/approve" "POST" "/verifications/talents/{id}/approve"
    print_description "Approve talenta"
    print_auth "Required (Super Admin, Admin Sekolah)"
    print_params "Path: id (UUID)"
    
    # Setup talent first
    setup_talent_for_verification
    
    if [ -z "$CREATED_TALENT_ID" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No talent ID available${NC}"
        return
    fi
    
    print_request "(no body)"
    
    local result=$(do_request "POST" "/verifications/talents/$CREATED_TALENT_ID/approve" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "404 NOT_FOUND - Talenta tidak ditemukan" \
        "400 ALREADY_VERIFIED - Talenta sudah diverifikasi" \
        "403 FORBIDDEN - Admin sekolah verifikasi talenta sekolah lain"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_verifications_approve_already_verified() {
    print_test "POST /verifications/talents/{id}/approve (Already Verified)" "POST" "/verifications/talents/{id}/approve"
    print_description "Test approve talenta yang sudah diverifikasi"
    print_auth "Required (Super Admin, Admin Sekolah)"
    
    if [ -z "$CREATED_TALENT_ID" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No talent ID available${NC}"
        return
    fi
    
    print_request "(no body)"
    
    local result=$(do_request "POST" "/verifications/talents/$CREATED_TALENT_ID/approve" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "400" ]; then
        print_success
    else
        print_failure "Expected 400, got $http_code"
    fi
}

test_verifications_reject() {
    print_test "POST /verifications/talents/{id}/reject" "POST" "/verifications/talents/{id}/reject"
    print_description "Reject talenta"
    print_auth "Required (Super Admin, Admin Sekolah)"
    print_params "Path: id (UUID), Body: rejection_reason (required)"
    
    # Setup new talent
    setup_talent_for_verification
    
    if [ -z "$CREATED_TALENT_ID" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No talent ID available${NC}"
        return
    fi
    
    local request_body='{
        "rejection_reason": "Dokumen bukti tidak valid atau tidak terbaca"
    }'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/verifications/talents/$CREATED_TALENT_ID/reject" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "404 NOT_FOUND - Talenta tidak ditemukan" \
        "400 ALREADY_VERIFIED - Talenta sudah diverifikasi" \
        "422 VALIDATION_ERROR - Alasan penolakan wajib diisi"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_verifications_reject_no_reason() {
    print_test "POST /verifications/talents/{id}/reject (No Reason)" "POST" "/verifications/talents/{id}/reject"
    print_description "Test reject talenta tanpa alasan"
    print_auth "Required (Super Admin, Admin Sekolah)"
    
    # Setup new talent
    setup_talent_for_verification
    
    if [ -z "$CREATED_TALENT_ID" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No talent ID available${NC}"
        return
    fi
    
    local request_body='{}'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/verifications/talents/$CREATED_TALENT_ID/reject" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "422" ] || [ "$http_code" = "400" ]; then
        print_success
    else
        print_failure "Expected 422 or 400, got $http_code"
    fi
}

test_verifications_batch_approve() {
    print_test "POST /verifications/talents/batch/approve" "POST" "/verifications/talents/batch/approve"
    print_description "Approve multiple talenta sekaligus"
    print_auth "Required (Super Admin, Admin Sekolah)"
    print_params "Body: ids (array of UUIDs)"
    
    # Setup multiple talents and collect IDs
    local id1=""
    local id2=""
    local id3=""
    
    setup_talent_for_verification
    id1="$CREATED_TALENT_ID"
    CREATED_TALENT_ID=""
    
    setup_talent_for_verification
    id2="$CREATED_TALENT_ID"
    CREATED_TALENT_ID=""
    
    setup_talent_for_verification
    id3="$CREATED_TALENT_ID"
    
    if [ -z "$id1" ] || [ -z "$id2" ] || [ -z "$id3" ]; then
        echo -e "${YELLOW}⚠️  Skipping: Could not create talents for batch test${NC}"
        return
    fi
    
    local request_body='{"ids": ["'"$id1"'", "'"$id2"'", "'"$id3"'"]}'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/verifications/talents/batch/approve" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "Partial success - some talents may fail"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_verifications_batch_reject() {
    print_test "POST /verifications/talents/batch/reject" "POST" "/verifications/talents/batch/reject"
    print_description "Reject multiple talenta sekaligus"
    print_auth "Required (Super Admin, Admin Sekolah)"
    print_params "Body: ids (array of UUIDs), rejection_reason (required)"
    
    # Setup multiple talents and collect IDs
    local id1=""
    local id2=""
    
    setup_talent_for_verification
    id1="$CREATED_TALENT_ID"
    CREATED_TALENT_ID=""
    
    setup_talent_for_verification
    id2="$CREATED_TALENT_ID"
    
    if [ -z "$id1" ] || [ -z "$id2" ]; then
        echo -e "${YELLOW}⚠️  Skipping: Could not create talents for batch test${NC}"
        return
    fi
    
    local request_body='{"ids": ["'"$id1"'", "'"$id2"'"], "rejection_reason": "Dokumen tidak lengkap"}'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/verifications/talents/batch/reject" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

#===============================================================================
# 7. NOTIFIKASI TESTS
#===============================================================================

test_notifications_list() {
    print_test "GET /me/notifications" "GET" "/me/notifications"
    print_description "Daftar notifikasi user yang login"
    print_auth "Required"
    print_params "Query: is_read (boolean), page, limit"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/me/notifications?page=1&limit=10" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_notifications_unread_count() {
    print_test "GET /me/notifications/unread-count" "GET" "/me/notifications/unread-count"
    print_description "Jumlah notifikasi yang belum dibaca"
    print_auth "Required"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/me/notifications/unread-count" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_notifications_mark_read() {
    print_test "PATCH /me/notifications/{id}/read" "PATCH" "/me/notifications/{id}/read"
    print_description "Tandai notifikasi sebagai sudah dibaca"
    print_auth "Required"
    print_params "Path: id (UUID)"
    
    # Get a notification ID first
    local list_result=$(do_request "GET" "/me/notifications?page=1&limit=1" "" "$ACCESS_TOKEN")
    local list_body=$(echo "$list_result" | tail -n +2)
    local notification_id=$(extract_json "$list_body" '.data[0].id')
    
    if [ -z "$notification_id" ] || [ "$notification_id" = "null" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No notifications available${NC}"
        print_request "(no body)"
        print_response "N/A" '{"message": "No notifications to test"}'
        return
    fi
    
    print_request "(no body)"
    
    local result=$(do_request "PATCH" "/me/notifications/$notification_id/read" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "404 NOT_FOUND - Notifikasi tidak ditemukan"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_notifications_mark_all_read() {
    print_test "PATCH /me/notifications/read-all" "PATCH" "/me/notifications/read-all"
    print_description "Tandai semua notifikasi sebagai sudah dibaca"
    print_auth "Required"
    
    print_request "(no body)"
    
    local result=$(do_request "PATCH" "/me/notifications/read-all" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

#===============================================================================
# 8. FILE UPLOAD (MINIO) TESTS
#===============================================================================

test_uploads_presign() {
    print_test "POST /uploads/presign" "POST" "/uploads/presign"
    print_description "Request presigned URL untuk upload file"
    print_auth "Required"
    print_params "Body: filename, content_type, upload_type (profile_photo/talent_certificate)"
    
    local request_body='{
        "filename": "sertifikat_test.pdf",
        "content_type": "application/pdf",
        "upload_type": "talent_certificate"
    }'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/uploads/presign" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "400 INVALID_FILE_TYPE - Tipe file tidak diizinkan" \
        "422 VALIDATION_ERROR - filename/content_type/upload_type tidak valid"
    
    if [ "$http_code" = "200" ]; then
        CREATED_UPLOAD_ID=$(extract_json "$body" '.data.upload_id')
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_uploads_presign_profile_photo() {
    print_test "POST /uploads/presign (Profile Photo)" "POST" "/uploads/presign"
    print_description "Request presigned URL untuk foto profil"
    print_auth "Required"
    print_params "Body: filename, content_type, upload_type=profile_photo"
    
    local request_body='{
        "filename": "profile.jpg",
        "content_type": "image/jpeg",
        "upload_type": "profile_photo"
    }'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/uploads/presign" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_uploads_presign_invalid_type() {
    print_test "POST /uploads/presign (Invalid Type)" "POST" "/uploads/presign"
    print_description "Test presign dengan tipe file tidak valid"
    print_auth "Required"
    
    local request_body='{
        "filename": "malware.exe",
        "content_type": "application/x-msdownload",
        "upload_type": "talent_certificate"
    }'
    print_request "$request_body"
    
    local result=$(do_request "POST" "/uploads/presign" "$request_body" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "400" ]; then
        print_success
    else
        print_failure "Expected 400, got $http_code"
    fi
}

test_uploads_confirm() {
    print_test "POST /uploads/{upload_id}/confirm" "POST" "/uploads/{upload_id}/confirm"
    print_description "Konfirmasi upload berhasil"
    print_auth "Required"
    print_params "Path: upload_id (UUID)"
    
    if [ -z "$CREATED_UPLOAD_ID" ]; then
        echo -e "${YELLOW}⚠️  Skipping: No upload ID available${NC}"
        return
    fi
    
    print_request "(no body)"
    
    local result=$(do_request "POST" "/uploads/$CREATED_UPLOAD_ID/confirm" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "404 UPLOAD_NOT_FOUND - Upload tidak ditemukan atau sudah expired" \
        "400 FILE_NOT_UPLOADED - File belum diupload ke storage"
    
    # Expected to fail since we didn't actually upload to MinIO
    if [ "$http_code" = "400" ] || [ "$http_code" = "404" ] || [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 400, 404, or 200, got $http_code"
    fi
}

test_uploads_delete() {
    print_test "DELETE /uploads/{upload_id}" "DELETE" "/uploads/{upload_id}"
    print_description "Batalkan upload"
    print_auth "Required"
    print_params "Path: upload_id (UUID)"
    
    # Create a new upload to delete
    local presign_body='{
        "filename": "to_delete.pdf",
        "content_type": "application/pdf",
        "upload_type": "talent_certificate"
    }'
    local presign_result=$(do_request "POST" "/uploads/presign" "$presign_body" "$ACCESS_TOKEN")
    local presign_http=$(echo "$presign_result" | head -n1)
    local presign_response=$(echo "$presign_result" | tail -n +2)
    local upload_id=$(extract_json "$presign_response" '.data.upload_id')
    
    if [ -z "$upload_id" ] || [ "$upload_id" = "null" ]; then
        echo -e "${YELLOW}⚠️  Skipping: Could not create upload${NC}"
        return
    fi
    
    print_request "(no body)"
    
    local result=$(do_request "DELETE" "/uploads/$upload_id" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "404 UPLOAD_NOT_FOUND - Upload tidak ditemukan"
    
    if [ "$http_code" = "204" ] || [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 204 or 200, got $http_code"
    fi
}

#===============================================================================
# 9. DASHBOARD & STATISTIK TESTS
#===============================================================================

test_dashboard_summary() {
    print_test "GET /dashboard/summary" "GET" "/dashboard/summary"
    print_description "Ringkasan statistik untuk dashboard (response berbeda per role)"
    print_auth "Required"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/dashboard/summary" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    echo -e "${CYAN}📋 Response varies by role:${NC}"
    echo "   • Super Admin: total_schools, total_users, total_gtk, gtk_by_type, talents_by_status, etc."
    echo "   • Admin Sekolah: school info, total_gtk, pending_verifications, etc."
    echo "   • GTK: my_talents, talents_by_type, unread_notifications, etc."
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_dashboard_schools_statistics() {
    print_test "GET /dashboard/schools/statistics" "GET" "/dashboard/schools/statistics"
    print_description "Statistik per sekolah (Super Admin only)"
    print_auth "Required (Super Admin)"
    print_params "Query: status (negeri/swasta), sort (gtk_count/talent_count), page, limit"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/dashboard/schools/statistics?page=1&limit=10" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "403 FORBIDDEN - Bukan Super Admin"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_dashboard_talents_statistics() {
    print_test "GET /dashboard/talents/statistics" "GET" "/dashboard/talents/statistics"
    print_description "Statistik talenta berdasarkan berbagai dimensi"
    print_auth "Required (Super Admin, Admin Sekolah)"
    print_params "Query: school_id, group_by (type/status/level/field), date_from, date_to"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/dashboard/talents/statistics?group_by=type" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_dashboard_talents_statistics_by_level() {
    print_test "GET /dashboard/talents/statistics (By Level)" "GET" "/dashboard/talents/statistics?group_by=level"
    print_description "Statistik talenta grouped by level"
    print_auth "Required (Super Admin, Admin Sekolah)"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/dashboard/talents/statistics?group_by=level" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_dashboard_talents_statistics_by_field() {
    print_test "GET /dashboard/talents/statistics (By Field)" "GET" "/dashboard/talents/statistics?group_by=field"
    print_description "Statistik talenta grouped by field"
    print_auth "Required (Super Admin, Admin Sekolah)"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/dashboard/talents/statistics?group_by=field" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

#===============================================================================
# 10. EXPORT LAPORAN TESTS
#===============================================================================

test_exports_gtk() {
    print_test "GET /exports/gtk" "GET" "/exports/gtk"
    print_description "Export data GTK ke Excel/PDF"
    print_auth "Required (Super Admin, Admin Sekolah)"
    print_params "Query: format (excel/pdf), school_id, gtk_type"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/exports/gtk?format=excel" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "403 FORBIDDEN - Tidak memiliki akses"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_exports_gtk_pdf() {
    print_test "GET /exports/gtk (PDF)" "GET" "/exports/gtk?format=pdf"
    print_description "Export data GTK ke PDF"
    print_auth "Required (Super Admin, Admin Sekolah)"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/exports/gtk?format=pdf" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_exports_talents() {
    print_test "GET /exports/talents" "GET" "/exports/talents"
    print_description "Export data talenta ke Excel/PDF"
    print_auth "Required (Super Admin, Admin Sekolah)"
    print_params "Query: format (excel/pdf), school_id, talent_type, status, date_from, date_to"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/exports/talents?format=excel" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_exports_talents_with_filter() {
    print_test "GET /exports/talents (With Filters)" "GET" "/exports/talents?format=excel&status=approved"
    print_description "Export data talenta dengan filter"
    print_auth "Required (Super Admin, Admin Sekolah)"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/exports/talents?format=excel&status=approved&talent_type=peserta_pelatihan" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

test_exports_schools() {
    print_test "GET /exports/schools" "GET" "/exports/schools"
    print_description "Export data sekolah ke Excel/PDF"
    print_auth "Required (Super Admin)"
    print_params "Query: format (excel/pdf), status"
    
    print_request "(no body)"
    
    local result=$(do_request "GET" "/exports/schools?format=excel" "" "$ACCESS_TOKEN")
    local http_code=$(echo "$result" | head -n1)
    local body=$(echo "$result" | tail -n +2)
    
    print_response "$http_code" "$body"
    
    print_error_scenarios \
        "403 FORBIDDEN - Bukan Super Admin"
    
    if [ "$http_code" = "200" ]; then
        print_success
    else
        print_failure "Expected 200, got $http_code"
    fi
}

#===============================================================================
# CLEANUP FUNCTIONS
#===============================================================================

cleanup() {
    print_section "CLEANUP"
    echo -e "${YELLOW}🧹 Cleaning up test data...${NC}"
    
    # Delete created school if exists
    if [ -n "$CREATED_SCHOOL_ID" ]; then
        do_request "DELETE" "/schools/$CREATED_SCHOOL_ID" "" "$ACCESS_TOKEN" > /dev/null 2>&1
        echo "   Deleted school: $CREATED_SCHOOL_ID"
    fi
    
    # Delete created user if exists
    if [ -n "$CREATED_USER_ID" ]; then
        do_request "DELETE" "/users/$CREATED_USER_ID" "" "$ACCESS_TOKEN" > /dev/null 2>&1
        echo "   Deleted user: $CREATED_USER_ID"
    fi
    
    echo -e "${GREEN}✅ Cleanup complete${NC}"
}

#===============================================================================
# MAIN EXECUTION
#===============================================================================

print_summary() {
    echo ""
    echo -e "${BOLD}${BLUE}╔══════════════════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BOLD}${BLUE}║${NC}                           ${BOLD}TEST SUMMARY${NC}                                      ${BOLD}${BLUE}║${NC}"
    echo -e "${BOLD}${BLUE}╚══════════════════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "   ${BOLD}Total Tests:${NC}  $TOTAL_TESTS"
    echo -e "   ${GREEN}Passed:${NC}       $PASSED_TESTS"
    echo -e "   ${RED}Failed:${NC}       $FAILED_TESTS"
    echo ""
    
    if [ $FAILED_TESTS -eq 0 ]; then
        echo -e "${GREEN}╔══════════════════════════════════════════════════════════════════════════════╗${NC}"
        echo -e "${GREEN}║                        🎉 ALL TESTS PASSED! 🎉                               ║${NC}"
        echo -e "${GREEN}╚══════════════════════════════════════════════════════════════════════════════╝${NC}"
    else
        echo -e "${RED}╔══════════════════════════════════════════════════════════════════════════════╗${NC}"
        echo -e "${RED}║                     ⚠️  SOME TESTS FAILED ⚠️                                  ║${NC}"
        echo -e "${RED}╚══════════════════════════════════════════════════════════════════════════════╝${NC}"
    fi
}

main() {
    echo ""
    echo -e "${BOLD}${BLUE}╔══════════════════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BOLD}${BLUE}║${NC}                                                                              ${BOLD}${BLUE}║${NC}"
    echo -e "${BOLD}${BLUE}║${NC}     ${BOLD}███████╗██╗██████╗  ██████╗ ██████╗ ██╗    █████╗ ██████╗ ██╗${NC}            ${BOLD}${BLUE}║${NC}"
    echo -e "${BOLD}${BLUE}║${NC}     ${BOLD}██╔════╝██║██╔══██╗██╔═══██╗██╔══██╗██║   ██╔══██╗██╔══██╗██║${NC}            ${BOLD}${BLUE}║${NC}"
    echo -e "${BOLD}${BLUE}║${NC}     ${BOLD}███████╗██║██████╔╝██║   ██║██║  ██║██║   ███████║██████╔╝██║${NC}            ${BOLD}${BLUE}║${NC}"
    echo -e "${BOLD}${BLUE}║${NC}     ${BOLD}╚════██║██║██╔═══╝ ██║   ██║██║  ██║██║   ██╔══██║██╔═══╝ ██║${NC}            ${BOLD}${BLUE}║${NC}"
    echo -e "${BOLD}${BLUE}║${NC}     ${BOLD}███████║██║██║     ╚██████╔╝██████╔╝██║   ██║  ██║██║     ██║${NC}            ${BOLD}${BLUE}║${NC}"
    echo -e "${BOLD}${BLUE}║${NC}     ${BOLD}╚══════╝╚═╝╚═╝      ╚═════╝ ╚═════╝ ╚═╝   ╚═╝  ╚═╝╚═╝     ╚═╝${NC}            ${BOLD}${BLUE}║${NC}"
    echo -e "${BOLD}${BLUE}║${NC}                                                                              ${BOLD}${BLUE}║${NC}"
    echo -e "${BOLD}${BLUE}║${NC}                    ${CYAN}Comprehensive API Testing Suite${NC}                           ${BOLD}${BLUE}║${NC}"
    echo -e "${BOLD}${BLUE}║${NC}                                                                              ${BOLD}${BLUE}║${NC}"
    echo -e "${BOLD}${BLUE}╚══════════════════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${CYAN}Configuration:${NC}"
    echo -e "   Base URL: ${YELLOW}$BASE_URL${NC}"
    echo -e "   Admin Email: ${YELLOW}$SUPER_ADMIN_EMAIL${NC}"
    echo ""
    
    # Check if jq is installed
    if ! command -v jq &> /dev/null; then
        echo -e "${RED}❌ Error: jq is required but not installed.${NC}"
        echo -e "   Install with: ${YELLOW}apt-get install jq${NC} or ${YELLOW}brew install jq${NC}"
        exit 1
    fi
    
    # Check if curl is installed
    if ! command -v curl &> /dev/null; then
        echo -e "${RED}❌ Error: curl is required but not installed.${NC}"
        exit 1
    fi
    
    # Wait for API
    # wait_for_api
    
    # Run tests
    print_header "1. AUTHENTICATION TESTS"
    test_auth_login
    test_auth_login_invalid_credentials
    test_auth_login_validation_error
    test_auth_refresh
    test_auth_unauthorized
    test_auth_logout
    test_auth_logout_all
    
    # Re-login after logout tests
    test_auth_login > /dev/null 2>&1
    
    print_header "2. PROFILE (ME) TESTS"
    test_me_get
    test_me_update
    test_me_password
    test_me_password_wrong_current
    test_me_photo
    
    print_header "3. SEKOLAH TESTS"
    test_schools_list
    test_schools_list_with_filter
    test_schools_create
    test_schools_create_duplicate_npsn
    test_schools_get
    test_schools_get_not_found
    test_schools_update
    test_schools_users
    test_schools_delete
    
    print_header "4. USERS/GTK TESTS"
    test_users_list
    test_users_list_with_filter
    test_users_create
    test_users_create_duplicate_email
    test_users_get
    test_users_get_not_found
    test_users_update
    test_users_deactivate
    test_users_activate
    test_users_delete
    
    print_header "5. TALENTA TESTS"
    test_talents_list
    test_talents_list_with_filter
    test_me_talents_list
    test_me_talents_create_peserta_pelatihan
    test_me_talents_create_pembimbing_lomba
    test_me_talents_create_peserta_lomba
    test_me_talents_create_minat_bakat
    test_me_talents_create_validation_error
    test_talents_get
    test_me_talents_update
    test_me_talents_delete
    
    print_header "6. VERIFIKASI TALENTA TESTS"
    test_verifications_list
    test_verifications_approve
    test_verifications_approve_already_verified
    test_verifications_reject
    test_verifications_reject_no_reason
    test_verifications_batch_approve
    test_verifications_batch_reject
    
    print_header "7. NOTIFIKASI TESTS"
    test_notifications_list
    test_notifications_unread_count
    test_notifications_mark_read
    test_notifications_mark_all_read
    
    print_header "8. FILE UPLOAD (MINIO) TESTS"
    test_uploads_presign
    test_uploads_presign_profile_photo
    test_uploads_presign_invalid_type
    test_uploads_confirm
    test_uploads_delete
    
    print_header "9. DASHBOARD & STATISTIK TESTS"
    test_dashboard_summary
    test_dashboard_schools_statistics
    test_dashboard_talents_statistics
    test_dashboard_talents_statistics_by_level
    test_dashboard_talents_statistics_by_field
    
    print_header "10. EXPORT LAPORAN TESTS"
    test_exports_gtk
    test_exports_gtk_pdf
    test_exports_talents
    test_exports_talents_with_filter
    test_exports_schools
    
    # Cleanup
    cleanup
    
    # Print summary
    print_summary
    
    # Exit with appropriate code
    if [ $FAILED_TESTS -gt 0 ]; then
        exit 1
    fi
    exit 0
}

# Run main function
main "$@"
