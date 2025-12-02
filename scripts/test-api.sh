#!/bin/bash

# API Test Script for Base App Service
# Usage: ./scripts/test-api.sh

set -e

API_URL=${API_URL:-http://localhost:8080/v1}
BASE_URL=${BASE_URL:-http://localhost:8080}

# Check if server is running
check_server() {
    if curl -s "$BASE_URL/health" &> /dev/null; then
        return 0
    fi
    # Try alternative port
    if curl -s "http://localhost:8081/health" &> /dev/null; then
        API_URL="http://localhost:8081/v1"
        BASE_URL="http://localhost:8081"
        return 0
    fi
    return 1
}

echo "Testing Base App API at $API_URL"
echo "=================================="
echo ""

# Check if server is running
echo "Checking if server is running..."
if ! check_server; then
    echo "❌ Server is not running!"
    echo ""
    echo "Please start the server first:"
    echo "  make dev   # (sets up database + runs server)"
    echo "  # OR"
    echo "  make run   # (if database is already set up)"
    echo ""
    echo "You can also check prerequisites with:"
    echo "  ./scripts/check-prerequisites.sh"
    exit 1
fi

echo "✓ Server is running"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if jq is available for JSON parsing
if command -v jq &> /dev/null; then
    USE_JQ=true
else
    USE_JQ=false
    echo -e "${YELLOW}Note: jq not found. Using basic parsing. Install jq for better output.${NC}"
    echo ""
fi

# Test 1: Health Check
echo -e "${YELLOW}Test 1: Health Check${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/health" || echo "000")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✓ Health check passed${NC}"
    if [ "$USE_JQ" = true ]; then
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
    else
        echo "Response: $body"
    fi
else
    echo -e "${RED}✗ Health check failed (HTTP $http_code)${NC}"
    echo "Response: $body"
    exit 1
fi
echo ""

# Test 2: Signup
echo -e "${YELLOW}Test 2: User Signup${NC}"
signup_response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/auth/signup" \
    -H "Content-Type: application/json" \
    -H "X-Product-Name: test-product" \
    -d '{
        "email": "test@example.com",
        "password": "TestPass123!",
        "name": "Test User",
        "terms_accepted": true,
        "terms_version": "1.0"
    }' || echo "000")

signup_http_code=$(echo "$signup_response" | tail -n1)
signup_body=$(echo "$signup_response" | sed '$d')

if [ "$signup_http_code" = "201" ]; then
    echo -e "${GREEN}✓ Signup successful${NC}"
    if [ "$USE_JQ" = true ]; then
        ACCESS_TOKEN=$(echo "$signup_body" | jq -r '.data.session.token' 2>/dev/null)
        REFRESH_TOKEN=$(echo "$signup_body" | jq -r '.data.session.refresh_token' 2>/dev/null)
        echo "$signup_body" | jq '.' 2>/dev/null
    else
        # Fallback parsing without jq
        ACCESS_TOKEN=$(echo "$signup_body" | grep -o '"token":"[^"]*' | head -1 | cut -d'"' -f4)
        REFRESH_TOKEN=$(echo "$signup_body" | grep -o '"refresh_token":"[^"]*' | head -1 | cut -d'"' -f4)
        echo "Response: $signup_body"
    fi
    if [ -n "$ACCESS_TOKEN" ] && [ "$ACCESS_TOKEN" != "null" ]; then
        echo "Access Token: ${ACCESS_TOKEN:0:30}..."
    fi
    if [ -n "$REFRESH_TOKEN" ] && [ "$REFRESH_TOKEN" != "null" ]; then
        echo "Refresh Token: ${REFRESH_TOKEN:0:30}..."
    fi
elif [ "$signup_http_code" = "409" ]; then
    echo -e "${YELLOW}⚠ User already exists (HTTP 409) - will use login instead${NC}"
    echo "Response: $signup_body"
else
    echo -e "${RED}✗ Signup failed (HTTP $signup_http_code)${NC}"
    echo "Response: $signup_body"
    echo -e "${YELLOW}Continuing with login test...${NC}"
fi
echo ""

# Test 3: Login
echo -e "${YELLOW}Test 3: Login${NC}"
login_response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "test@example.com",
        "password": "TestPass123!"
    }' || echo "000")

login_http_code=$(echo "$login_response" | tail -n1)
login_body=$(echo "$login_response" | sed '$d')

if [ "$login_http_code" = "200" ]; then
    echo -e "${GREEN}✓ Login successful${NC}"
    if [ "$USE_JQ" = true ]; then
        LOGIN_ACCESS_TOKEN=$(echo "$login_body" | jq -r '.data.session.token' 2>/dev/null)
        LOGIN_REFRESH_TOKEN=$(echo "$login_body" | jq -r '.data.session.refresh_token' 2>/dev/null)
        echo "$login_body" | jq '.' 2>/dev/null
    else
        LOGIN_ACCESS_TOKEN=$(echo "$login_body" | grep -o '"token":"[^"]*' | head -1 | cut -d'"' -f4)
        LOGIN_REFRESH_TOKEN=$(echo "$login_body" | grep -o '"refresh_token":"[^"]*' | head -1 | cut -d'"' -f4)
        echo "Response: $login_body"
    fi
    # Use login token if signup token wasn't extracted properly
    if [ -z "$ACCESS_TOKEN" ] || [ "$ACCESS_TOKEN" = "null" ]; then
        ACCESS_TOKEN="$LOGIN_ACCESS_TOKEN"
    fi
    if [ -z "$REFRESH_TOKEN" ] || [ "$REFRESH_TOKEN" = "null" ]; then
        REFRESH_TOKEN="$LOGIN_REFRESH_TOKEN"
    fi
else
    echo -e "${RED}✗ Login failed (HTTP $login_http_code)${NC}"
    echo "Response: $login_body"
fi
echo ""

# Test 4: Get Current User
if [ -n "$ACCESS_TOKEN" ] && [ "$ACCESS_TOKEN" != "null" ]; then
    echo -e "${YELLOW}Test 4: Get Current User${NC}"
    user_response=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/users/me" \
        -H "Authorization: Bearer $ACCESS_TOKEN" || echo "000")
    
    user_http_code=$(echo "$user_response" | tail -n1)
    user_body=$(echo "$user_response" | sed '$d')
    
    if [ "$user_http_code" = "200" ]; then
        echo -e "${GREEN}✓ Get current user successful${NC}"
        if [ "$USE_JQ" = true ]; then
            echo "$user_body" | jq '.' 2>/dev/null
        else
            echo "Response: $user_body"
        fi
    else
        echo -e "${RED}✗ Get current user failed (HTTP $user_http_code)${NC}"
        echo "Response: $user_body"
    fi
    echo ""
else
    echo -e "${YELLOW}Skipping authenticated tests - no access token available${NC}"
    echo ""
fi

# Test 5: Update Theme
if [ -n "$ACCESS_TOKEN" ] && [ "$ACCESS_TOKEN" != "null" ]; then
    echo -e "${YELLOW}Test 5: Update Theme${NC}"
    theme_response=$(curl -s -w "\n%{http_code}" -X PUT "$API_URL/users/me/settings/theme" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "theme": "dark",
            "contrast": "high"
        }' || echo "000")
    
    theme_http_code=$(echo "$theme_response" | tail -n1)
    theme_body=$(echo "$theme_response" | sed '$d')
    
    if [ "$theme_http_code" = "200" ]; then
        echo -e "${GREEN}✓ Update theme successful${NC}"
        if [ "$USE_JQ" = true ]; then
            echo "$theme_body" | jq '.' 2>/dev/null
        else
            echo "Response: $theme_body"
        fi
    else
        echo -e "${RED}✗ Update theme failed (HTTP $theme_http_code)${NC}"
        echo "Response: $theme_body"
    fi
    echo ""
fi

# Test 6: Get Theme
if [ -n "$ACCESS_TOKEN" ] && [ "$ACCESS_TOKEN" != "null" ]; then
    echo -e "${YELLOW}Test 6: Get Theme${NC}"
    get_theme_response=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/users/me/settings/theme" \
        -H "Authorization: Bearer $ACCESS_TOKEN" || echo "000")
    
    get_theme_http_code=$(echo "$get_theme_response" | tail -n1)
    get_theme_body=$(echo "$get_theme_response" | sed '$d')
    
    if [ "$get_theme_http_code" = "200" ]; then
        echo -e "${GREEN}✓ Get theme successful${NC}"
        if [ "$USE_JQ" = true ]; then
            echo "$get_theme_body" | jq '.' 2>/dev/null
        else
            echo "Response: $get_theme_body"
        fi
    else
        echo -e "${RED}✗ Get theme failed (HTTP $get_theme_http_code)${NC}"
        echo "Response: $get_theme_body"
    fi
    echo ""
fi

# Test 7: Refresh Token
if [ -n "$REFRESH_TOKEN" ] && [ "$REFRESH_TOKEN" != "null" ]; then
    echo -e "${YELLOW}Test 7: Refresh Token${NC}"
    refresh_response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/auth/refresh" \
        -H "Content-Type: application/json" \
        -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}" || echo "000")
    
    refresh_http_code=$(echo "$refresh_response" | tail -n1)
    refresh_body=$(echo "$refresh_response" | sed '$d')
    
    if [ "$refresh_http_code" = "200" ]; then
        echo -e "${GREEN}✓ Refresh token successful${NC}"
        if [ "$USE_JQ" = true ]; then
            NEW_ACCESS_TOKEN=$(echo "$refresh_body" | jq -r '.data.token' 2>/dev/null)
            echo "$refresh_body" | jq '.' 2>/dev/null
        else
            NEW_ACCESS_TOKEN=$(echo "$refresh_body" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
            echo "Response: $refresh_body"
        fi
        if [ -n "$NEW_ACCESS_TOKEN" ] && [ "$NEW_ACCESS_TOKEN" != "null" ]; then
            echo "New Access Token: ${NEW_ACCESS_TOKEN:0:30}..."
        fi
    else
        echo -e "${RED}✗ Refresh token failed (HTTP $refresh_http_code)${NC}"
        echo "Response: $refresh_body"
    fi
    echo ""
fi

# Test 8: Logout
if [ -n "$ACCESS_TOKEN" ] && [ "$ACCESS_TOKEN" != "null" ]; then
    echo -e "${YELLOW}Test 8: Logout${NC}"
    logout_response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/auth/logout" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"revoke_all_sessions": false}' || echo "000")
    
    logout_http_code=$(echo "$logout_response" | tail -n1)
    logout_body=$(echo "$logout_response" | sed '$d')
    
    if [ "$logout_http_code" = "200" ]; then
        echo -e "${GREEN}✓ Logout successful${NC}"
        if [ "$USE_JQ" = true ]; then
            echo "$logout_body" | jq '.' 2>/dev/null
        else
            echo "Response: $logout_body"
        fi
    else
        echo -e "${RED}✗ Logout failed (HTTP $logout_http_code)${NC}"
        echo "Response: $logout_body"
    fi
    echo ""
fi

echo "=================================="
echo -e "${GREEN}API Tests Completed!${NC}"

