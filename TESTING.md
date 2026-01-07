# Testing Guide - Grade Management API

Complete guide for testing the Grade Management API with examples and common scenarios.

## Table of Contents

1. [Setup](#setup)
2. [Basic Testing](#basic-testing)
3. [Complete Workflow Tests](#complete-workflow-tests)
4. [Validation Tests](#validation-tests)
5. [Error Scenario Tests](#error-scenario-tests)
6. [PowerShell Testing](#powershell-testing)
7. [Automated Testing](#automated-testing)

## Setup

### Start the Server

```bash
# Build and run
go build && ./techwave

# Or run directly
go run main.go
```

Expected output:
```
Starting Grade Management API on port :8080
```

### Verify Server is Running

```bash
curl http://localhost:8080
```

Expected: `Grade Management API - Ready for AI delegation!`

## Basic Testing

### Test 1: Create a Simple Enrollment

**Request:**
```bash
curl -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "student-001",
    "course_id": "course-101",
    "status": "pending"
  }'
```

**Expected Response (201):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "student_id": "student-001",
  "course_id": "course-101",
  "enrollment_date": "2024-01-15T10:00:00Z",
  "status": "pending",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

**What to Check:**
- ✅ Status code is 201
- ✅ Response contains auto-generated UUID in `id`
- ✅ `enrollment_date`, `created_at`, `updated_at` are set
- ✅ All input fields are preserved

### Test 2: Retrieve All Enrollments

**Request:**
```bash
curl http://localhost:8080/api/enrollments
```

**Expected Response (200):**
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "student_id": "student-001",
    "course_id": "course-101",
    "enrollment_date": "2024-01-15T10:00:00Z",
    "status": "pending",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
]
```

**What to Check:**
- ✅ Status code is 200
- ✅ Response is an array
- ✅ Contains previously created enrollment

### Test 3: Retrieve Specific Enrollment

**Request:**
```bash
# Replace {id} with actual enrollment ID
curl http://localhost:8080/api/enrollments/123e4567-e89b-12d3-a456-426614174000
```

**Expected Response (200):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "student_id": "student-001",
  "course_id": "course-101",
  "enrollment_date": "2024-01-15T10:00:00Z",
  "status": "pending",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

**What to Check:**
- ✅ Status code is 200
- ✅ Returned enrollment matches requested ID

### Test 4: Update Enrollment

**Request:**
```bash
# Replace {id} with actual enrollment ID
curl -X PUT http://localhost:8080/api/enrollments/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "student-001",
    "course_id": "course-101",
    "status": "active"
  }'
```

**Expected Response (200):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "student_id": "student-001",
  "course_id": "course-101",
  "enrollment_date": "2024-01-15T10:00:00Z",
  "status": "active",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T11:30:00Z"
}
```

**What to Check:**
- ✅ Status code is 200
- ✅ Status changed from "pending" to "active"
- ✅ `updated_at` timestamp is newer than `created_at`
- ✅ `created_at` remained unchanged

### Test 5: Delete Enrollment

**Request:**
```bash
# Replace {id} with actual enrollment ID
curl -X DELETE http://localhost:8080/api/enrollments/123e4567-e89b-12d3-a456-426614174000
```

**Expected Response (200):**
```json
{
  "message": "Enrollment deleted successfully"
}
```

**What to Check:**
- ✅ Status code is 200
- ✅ Success message returned
- ✅ Subsequent GET returns 404

## Complete Workflow Tests

### Scenario: Student Enrollment Lifecycle

This tests the complete journey from enrollment to completion.

**Step 1: Create Pending Enrollment**
```bash
ENROLLMENT=$(curl -s -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "student-123",
    "course_id": "golang-101",
    "status": "pending"
  }')

ENROLLMENT_ID=$(echo $ENROLLMENT | jq -r '.id')
echo "Created enrollment: $ENROLLMENT_ID"
```

**Step 2: Verify Enrollment Created**
```bash
curl http://localhost:8080/api/enrollments/$ENROLLMENT_ID | jq .
```

Expected: Status is "pending"

**Step 3: Activate Enrollment (Course Started)**
```bash
curl -s -X PUT http://localhost:8080/api/enrollments/$ENROLLMENT_ID \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "student-123",
    "course_id": "golang-101",
    "status": "active"
  }' | jq .
```

Expected: Status is "active"

**Step 4: Complete Course**
```bash
curl -s -X PUT http://localhost:8080/api/enrollments/$ENROLLMENT_ID \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "student-123",
    "course_id": "golang-101",
    "status": "completed"
  }' | jq .
```

Expected: Status is "completed"

**Step 5: View Final State**
```bash
curl http://localhost:8080/api/enrollments/$ENROLLMENT_ID | jq .
```

**Step 6: Cleanup**
```bash
curl -X DELETE http://localhost:8080/api/enrollments/$ENROLLMENT_ID
```

### Scenario: Multiple Enrollments

Test handling multiple enrollments simultaneously.

```bash
# Create 3 enrollments
for i in {1..3}; do
  curl -s -X POST http://localhost:8080/api/enrollments \
    -H "Content-Type: application/json" \
    -d "{
      \"student_id\": \"student-00$i\",
      \"course_id\": \"course-101\",
      \"status\": \"active\"
    }" | jq -r '.id'
done

# List all enrollments
curl -s http://localhost:8080/api/enrollments | jq '. | length'
```

Expected: Returns 3

## Validation Tests

### Test Missing student_id

**Request:**
```bash
curl -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{
    "course_id": "course-101",
    "status": "active"
  }'
```

**Expected Response (400):**
```json
{
  "error": "student_id is required"
}
```

### Test Missing course_id

**Request:**
```bash
curl -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "student-001",
    "status": "active"
  }'
```

**Expected Response (400):**
```json
{
  "error": "course_id is required"
}
```

### Test Missing status

**Request:**
```bash
curl -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "student-001",
    "course_id": "course-101"
  }'
```

**Expected Response (400):**
```json
{
  "error": "status is required"
}
```

### Test Invalid Status

**Request:**
```bash
curl -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "student-001",
    "course_id": "course-101",
    "status": "invalid_status"
  }'
```

**Expected Response (400):**
```json
{
  "error": "status must be one of: pending, active, completed"
}
```

### Test All Valid Statuses

**Test 1: pending**
```bash
curl -s -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{"student_id":"s1","course_id":"c1","status":"pending"}' \
  | jq -r '.status'
```
Expected: "pending"

**Test 2: active**
```bash
curl -s -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{"student_id":"s2","course_id":"c1","status":"active"}' \
  | jq -r '.status'
```
Expected: "active"

**Test 3: completed**
```bash
curl -s -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{"student_id":"s3","course_id":"c1","status":"completed"}' \
  | jq -r '.status'
```
Expected: "completed"

## Error Scenario Tests

### Test Get Non-existent Enrollment

**Request:**
```bash
curl http://localhost:8080/api/enrollments/00000000-0000-0000-0000-000000000000
```

**Expected Response (404):**
```json
{
  "error": "Enrollment not found"
}
```

### Test Update Non-existent Enrollment

**Request:**
```bash
curl -X PUT http://localhost:8080/api/enrollments/00000000-0000-0000-0000-000000000000 \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "student-001",
    "course_id": "course-101",
    "status": "active"
  }'
```

**Expected Response (404):**
```json
{
  "error": "Enrollment not found"
}
```

### Test Delete Non-existent Enrollment

**Request:**
```bash
curl -X DELETE http://localhost:8080/api/enrollments/00000000-0000-0000-0000-000000000000
```

**Expected Response (404):**
```json
{
  "error": "Enrollment not found"
}
```

### Test Malformed JSON

**Request:**
```bash
curl -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{invalid json}'
```

**Expected Response (400):**
```json
{
  "error": "Invalid request payload"
}
```

## PowerShell Testing

### Complete Test Suite

```powershell
# Function to test API endpoint
function Test-Enrollment {
    param(
        [string]$TestName,
        [string]$Method,
        [string]$Endpoint,
        [object]$Body = $null
    )
    
    Write-Host "`n=== $TestName ===" -ForegroundColor Cyan
    
    try {
        if ($Body) {
            $jsonBody = $Body | ConvertTo-Json
            $result = Invoke-RestMethod -Uri "http://localhost:8080$Endpoint" `
                -Method $Method `
                -Body $jsonBody `
                -ContentType "application/json"
        } else {
            $result = Invoke-RestMethod -Uri "http://localhost:8080$Endpoint" `
                -Method $Method
        }
        
        Write-Host "✅ Success" -ForegroundColor Green
        $result | ConvertTo-Json -Depth 10
        return $result
    } catch {
        Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
        return $null
    }
}

# Test 1: Create Enrollment
$createBody = @{
    student_id = "ps-student-001"
    course_id = "ps-course-101"
    status = "pending"
}
$enrollment = Test-Enrollment -TestName "Create Enrollment" -Method Post -Endpoint "/api/enrollments" -Body $createBody
$enrollmentId = $enrollment.id

# Test 2: Get All Enrollments
Test-Enrollment -TestName "Get All Enrollments" -Method Get -Endpoint "/api/enrollments"

# Test 3: Get Specific Enrollment
Test-Enrollment -TestName "Get Enrollment by ID" -Method Get -Endpoint "/api/enrollments/$enrollmentId"

# Test 4: Update Enrollment
$updateBody = @{
    student_id = "ps-student-001"
    course_id = "ps-course-101"
    status = "active"
}
Test-Enrollment -TestName "Update Enrollment" -Method Put -Endpoint "/api/enrollments/$enrollmentId" -Body $updateBody

# Test 5: Update to Completed
$completeBody = @{
    student_id = "ps-student-001"
    course_id = "ps-course-101"
    status = "completed"
}
Test-Enrollment -TestName "Complete Enrollment" -Method Put -Endpoint "/api/enrollments/$enrollmentId" -Body $completeBody

# Test 6: Delete Enrollment
Test-Enrollment -TestName "Delete Enrollment" -Method Delete -Endpoint "/api/enrollments/$enrollmentId"

# Test 7: Verify Deletion (should fail)
Write-Host "`n=== Verify Deletion (Expect 404) ===" -ForegroundColor Cyan
try {
    Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments/$enrollmentId" -Method Get
    Write-Host "❌ Should have returned 404" -ForegroundColor Red
} catch {
    Write-Host "✅ Correctly returned 404" -ForegroundColor Green
}

Write-Host "`n=== All Tests Complete ===" -ForegroundColor Green
```

### Save and Run

Save the above script as `test-api.ps1` and run:

```powershell
.\test-api.ps1
```

## Automated Testing

### Create Test Script

Create `test-all.sh`:

```bash
#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
TESTS_PASSED=0
TESTS_FAILED=0

# Test function
test_api() {
    local test_name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local expected_status=$5
    
    echo -e "\n${YELLOW}Testing: $test_name${NC}"
    
    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X $method "http://localhost:8080$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X $method "http://localhost:8080$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi
    
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)
    
    if [ "$status_code" = "$expected_status" ]; then
        echo -e "${GREEN}✅ PASS${NC} (Status: $status_code)"
        ((TESTS_PASSED++))
        echo "$body" | jq . 2>/dev/null || echo "$body"
    else
        echo -e "${RED}❌ FAIL${NC} (Expected: $expected_status, Got: $status_code)"
        ((TESTS_FAILED++))
        echo "$body"
    fi
}

# Start tests
echo "Starting API Tests..."

# Create enrollment
test_api "Create Enrollment" "POST" "/api/enrollments" \
    '{"student_id":"test-001","course_id":"course-101","status":"pending"}' \
    "201"

# Get all enrollments
test_api "Get All Enrollments" "GET" "/api/enrollments" "" "200"

# Test validation errors
test_api "Missing student_id" "POST" "/api/enrollments" \
    '{"course_id":"course-101","status":"pending"}' \
    "400"

test_api "Invalid status" "POST" "/api/enrollments" \
    '{"student_id":"test-001","course_id":"course-101","status":"invalid"}' \
    "400"

test_api "Get non-existent" "GET" "/api/enrollments/00000000-0000-0000-0000-000000000000" "" "404"

# Summary
echo -e "\n=========================================="
echo -e "Test Summary:"
echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Failed: $TESTS_FAILED${NC}"
echo -e "=========================================="

if [ $TESTS_FAILED -eq 0 ]; then
    exit 0
else
    exit 1
fi
```

Make executable and run:

```bash
chmod +x test-all.sh
./test-all.sh
```

## Performance Testing

### Simple Load Test

```bash
# Create 100 enrollments as fast as possible
for i in {1..100}; do
  curl -s -X POST http://localhost:8080/api/enrollments \
    -H "Content-Type: application/json" \
    -d "{\"student_id\":\"student-$i\",\"course_id\":\"course-101\",\"status\":\"active\"}" &
done
wait

# Check count
COUNT=$(curl -s http://localhost:8080/api/enrollments | jq '. | length')
echo "Created $COUNT enrollments"
```

Expected: Should create 100 enrollments successfully (check for thread-safety).

## Test Checklist

Before deployment, verify:

- [ ] All CRUD operations work correctly
- [ ] Validation rules are enforced
- [ ] Error responses have correct status codes
- [ ] Timestamps (created_at, updated_at) work correctly
- [ ] IDs are auto-generated as UUIDs
- [ ] Thread-safe operations (concurrent requests)
- [ ] Non-existent resources return 404
- [ ] Malformed JSON returns 400
- [ ] All status values (pending, active, completed) work
- [ ] Update changes updated_at timestamp
- [ ] Delete makes resource unavailable

## Continuous Testing

Run tests regularly during development:

```bash
# Watch mode - re-run tests on code changes
while true; do
  clear
  ./test-all.sh
  sleep 5
done
```

---

**Happy Testing! 🧪**
