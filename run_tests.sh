#!/bin/bash

BASE_URL="http://localhost:8000"
FAILED=0
PASSED=0

# Helper function to run a test
# Usage: run_test <name> <endpoint> <data_file> <expected_string>
run_test() {
    local name=$1
    local endpoint=$2
    local data_file=$3
    local expected=$4

    echo -n "Test: $name ... "
    
    # Run curl and capture response
    response=$(curl -s -X POST "$BASE_URL$endpoint" -H "Content-Type: application/json" -d @"$data_file")
    
    # Check if expected string exists in response
    if echo "$response" | grep -q "$expected"; then
        echo -e "\e[32mPASS\e[0m"
        PASSED=$((PASSED + 1))
    else
        echo -e "\e[31mFAIL\e[0m"
        echo "  Expected to find: $expected"
        echo "  Actual response:  $response"
        FAILED=$((FAILED + 1))
    fi
}

echo "Starting Automated Admission Webhook Tests..."
echo "============================================"

# Validation Tests
run_test "Validation Allowed Case" "/api/validate" "tests/validate-allowed.json" '"allowed":true'
run_test "Validation Disallowed Case" "/api/validate" "tests/validate-disallowed.json" '"allowed":false'
run_test "Validation Developer Case" "/api/validate" "tests/validate-another-allowed.json" '"allowed":true'

# Mutation Tests
run_test "Mutation Add SC" "/api/mutate" "tests/mutate-missing-sc.json" '"patch":"W3sib3AiOiJhZGQiLCJwYXRoIjoiL3NwZWMvY29udGFpbmVycy8wL3NlY3VyaXR5Q29udGV4dCIsInZhbHVlIjp7InJ1bkFzVXNlciI6MTAwMH19XQ=="'
run_test "Mutation Fix User" "/api/mutate" "tests/mutate-wrong-user.json" '"patch":"W3sib3AiOiJhZGQiLCJwYXRoIjoiL3NwZWMvY29udGFpbmVycy8wL3NlY3VyaXR5Q29udGV4dC9ydW5Bc1VzZXIiLCJ2YWx1ZSI6MTAwMH1d"'
run_test "Mutation No Op" "/api/mutate" "tests/mutate-correct-user.json" '"patchType":""'

echo "============================================"
echo "Summary: $PASSED passed, $FAILED failed"

if [ $FAILED -ne 0 ]; then
    exit 1
fi
