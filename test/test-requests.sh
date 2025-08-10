#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

PROXY_URL="http://localhost:8080"

echo_header() {
    echo -e "\n${BLUE}=== $1 ===${NC}"
}

echo_test() {
    echo -e "${YELLOW}Testing: $1${NC}"
}

echo_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

echo_error() {
    echo -e "${RED}✗ $1${NC}"
}

test_request() {
    local description="$1"
    local url="$2"
    local host_header="$3"
    
    echo_test "$description"
    
    if [ -n "$host_header" ]; then
        response=$(curl -s -H "Host: $host_header" "$url" 2>/dev/null)
    else
        response=$(curl -s "$url" 2>/dev/null)
    fi
    
    if [ $? -eq 0 ]; then
        server_name=$(echo "$response" | grep -o '"server": "[^"]*"' | cut -d'"' -f4)
        path=$(echo "$response" | grep -o '"path": "[^"]*"' | cut -d'"' -f4)
        echo -e "  Response: Server=${server_name}, Path=${path}"
        echo_success "Request successful"
    else
        echo_error "Request failed"
    fi
    
    sleep 0.5
}

# Check if proxy is running
echo_header "Checking Proxy Server"
if curl -s "$PROXY_URL" >/dev/null 2>&1; then
    echo_success "Proxy server is running on port 8080"
else
    echo_error "Proxy server is not running on port 8080"
    echo "Please start the proxy server first: ./go-gate"
    exit 1
fi

# Test 1: Path-based routing - API endpoints
echo_header "Testing Path-based Routing (/api/*)"
for i in {1..5}; do
    test_request "API Request #$i" "$PROXY_URL/api/users"
done

# Test 2: Default routing (should load balance across all servers)
echo_header "Testing Default Routing (load balancing)"
for i in {1..6}; do
    test_request "Default Request #$i" "$PROXY_URL/default/path"
done

# Test 3: Host-based routing - admin.example.com
echo_header "Testing Host-based Routing (admin.example.com)"
for i in {1..3}; do
    test_request "Admin Request #$i" "$PROXY_URL/admin/dashboard" "admin.example.com"
done

# Test 4: Host-based routing - *.example.com (should go to web server)
echo_header "Testing Wildcard Host Routing (*.example.com)"
for i in {1..3}; do
    test_request "Web Request #$i" "$PROXY_URL/home" "www.example.com"
done

# Test 5: Health check endpoints
echo_header "Testing Health Check Endpoints"
test_request "Health Check via API route" "$PROXY_URL/api/health"
test_request "Health Check via default route" "$PROXY_URL/health"

# Test 6: POST requests
echo_header "Testing POST Requests"
echo_test "POST to API endpoint"
response=$(curl -s -X POST -H "Content-Type: application/json" -d '{"test": "data"}' "$PROXY_URL/api/submit" 2>/dev/null)
if [ $? -eq 0 ]; then
    server_name=$(echo "$response" | grep -o '"server": "[^"]*"' | cut -d'"' -f4)
    echo -e "  Response: Server=${server_name}"
    echo_success "POST request successful"
else
    echo_error "POST request failed"
fi

echo_header "Test Summary"
echo "All tests completed. Check the proxy logs to see request routing in action."
echo ""
echo "To see real-time logs, run the proxy with: ./go-gate"
echo "To add custom hosts for domain testing, add to /etc/hosts:"
echo "127.0.0.1 admin.example.com"
echo "127.0.0.1 www.example.com"