# Grade Management API

A RESTful API for managing student enrollments in courses, built with Go. This API provides complete CRUD operations for enrollment management with validation, error handling, and comprehensive documentation.

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Development Guide](#development-guide)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)

## 🎯 Overview

The Grade Management API helps educational institutions manage student enrollments efficiently. It tracks enrollment status throughout the student lifecycle from enrollment to course completion.

### Technology Stack

- **Language**: Go 1.x
- **HTTP Router**: Gorilla Mux
- **Storage**: In-memory (production-ready for database integration)
- **Architecture**: Clean architecture with separation of concerns

### Project Structure

```
techwave/
├── main.go              # Application entry point and HTTP server setup
├── models/              # Domain entities and business logic
│   └── enrollment.go    # Enrollment model with validation
├── handlers/            # HTTP request handlers
│   └── enrollment_handler.go
├── repository/          # Data storage layer
│   └── enrollment_repository.go
├── go.mod              # Go module dependencies
└── README.md           # This file
```

## ✨ Features

- **Full CRUD Operations**: Create, Read, Update, Delete enrollments
- **Data Validation**: Business rule enforcement for enrollment data
- **Status Management**: Track enrollment lifecycle (pending → active → completed)
- **Thread-Safe**: Concurrent request handling with proper synchronization
- **RESTful Design**: Standard HTTP methods and status codes
- **Error Handling**: Comprehensive error responses with proper HTTP codes

## 🚀 Getting Started

### Prerequisites

- Go 1.19 or higher
- Git (for cloning the repository)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd techwave
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Build the application**
   ```bash
   go build -o techwave
   ```

4. **Run the server**
   ```bash
   go run main.go
   ```

   Or use the built executable:
   ```bash
   ./techwave
   ```

The server will start on `http://localhost:8080`

### Quick Test

Test that the server is running:

**Using curl:**
```bash
curl http://localhost:8080
```

**Using PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080" -Method Get
```

Expected response: `Grade Management API - Ready for AI delegation!`

## 📚 API Documentation

### Base URL

```
http://localhost:8080/api
```

### Endpoints

#### 1. Create Enrollment

**POST** `/api/enrollments`

Creates a new student enrollment.

**Request Body:**
```json
{
  "student_id": "student-123",
  "course_id": "course-456",
  "status": "active",
  "enrollment_date": "2024-01-15T10:00:00Z"
}
```

**Required Fields:**
- `student_id`: Student identifier (string)
- `course_id`: Course identifier (string)
- `status`: One of "pending", "active", or "completed"

**Optional Fields:**
- `enrollment_date`: Date of enrollment (defaults to current time)

**Response (201 Created):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "student_id": "student-123",
  "course_id": "course-456",
  "enrollment_date": "2024-01-15T10:00:00Z",
  "status": "active",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

**Example (curl):**
```bash
curl -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "student-123",
    "course_id": "course-456",
    "status": "active"
  }'
```

**Example (PowerShell):**
```powershell
$body = @{
    student_id = "student-123"
    course_id = "course-456"
    status = "active"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments" `
  -Method Post `
  -Body $body `
  -ContentType "application/json"
```

#### 2. Get All Enrollments

**GET** `/api/enrollments`

Retrieves all enrollments in the system.

**Response (200 OK):**
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "student_id": "student-123",
    "course_id": "course-456",
    "enrollment_date": "2024-01-15T10:00:00Z",
    "status": "active",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
]
```

**Example (curl):**
```bash
curl http://localhost:8080/api/enrollments
```

**Example (PowerShell):**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments" -Method Get
```

#### 3. Get Enrollment by ID

**GET** `/api/enrollments/{id}`

Retrieves a specific enrollment by its ID.

**URL Parameters:**
- `id`: Enrollment UUID

**Response (200 OK):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "student_id": "student-123",
  "course_id": "course-456",
  "enrollment_date": "2024-01-15T10:00:00Z",
  "status": "active",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

**Example (curl):**
```bash
curl http://localhost:8080/api/enrollments/123e4567-e89b-12d3-a456-426614174000
```

**Example (PowerShell):**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments/123e4567-e89b-12d3-a456-426614174000" -Method Get
```

#### 4. Update Enrollment

**PUT** `/api/enrollments/{id}`

Updates an existing enrollment. All fields must be provided (full replacement).

**URL Parameters:**
- `id`: Enrollment UUID

**Request Body:**
```json
{
  "student_id": "student-123",
  "course_id": "course-456",
  "status": "completed",
  "enrollment_date": "2024-01-15T10:00:00Z"
}
```

**Response (200 OK):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "student_id": "student-123",
  "course_id": "course-456",
  "enrollment_date": "2024-01-15T10:00:00Z",
  "status": "completed",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T11:30:00Z"
}
```

**Example (curl):**
```bash
curl -X PUT http://localhost:8080/api/enrollments/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "student-123",
    "course_id": "course-456",
    "status": "completed"
  }'
```

**Example (PowerShell):**
```powershell
$body = @{
    student_id = "student-123"
    course_id = "course-456"
    status = "completed"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments/123e4567-e89b-12d3-a456-426614174000" `
  -Method Put `
  -Body $body `
  -ContentType "application/json"
```

#### 5. Delete Enrollment

**DELETE** `/api/enrollments/{id}`

Deletes an enrollment. This is a permanent operation.

**URL Parameters:**
- `id`: Enrollment UUID

**Response (200 OK):**
```json
{
  "message": "Enrollment deleted successfully"
}
```

**Example (curl):**
```bash
curl -X DELETE http://localhost:8080/api/enrollments/123e4567-e89b-12d3-a456-426614174000
```

**Example (PowerShell):**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments/123e4567-e89b-12d3-a456-426614174000" -Method Delete
```

### HTTP Status Codes

- **200 OK**: Request succeeded
- **201 Created**: Resource created successfully
- **400 Bad Request**: Invalid request data or validation failure
- **404 Not Found**: Resource not found
- **409 Conflict**: Resource already exists
- **500 Internal Server Error**: Server-side error

### Error Response Format

```json
{
  "error": "Descriptive error message"
}
```

## 👨‍💻 Development Guide

### Project Architecture

The application follows clean architecture principles with three main layers:

1. **Models Layer** (`models/`)
   - Domain entities and business logic
   - Data validation rules
   - No external dependencies

2. **Repository Layer** (`repository/`)
   - Data storage and retrieval
   - Currently in-memory, can be replaced with database
   - Thread-safe operations

3. **Handlers Layer** (`handlers/`)
   - HTTP request handling
   - JSON serialization/deserialization
   - Error response formatting

### Adding New Features

**To add a new field to Enrollment:**

1. Update the `Enrollment` struct in `models/enrollment.go`
2. Add validation logic in the `Validate()` method if needed
3. Update godoc comments to document the new field

**To add a new endpoint:**

1. Create handler method in `handlers/enrollment_handler.go`
2. Add route in `main.go`
3. Document with godoc comments including examples

### Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use godoc-style comments for all public functions
- Include code examples in complex function documentation
- Keep functions focused and small

### Generating Documentation

Generate HTML documentation from godoc comments:

```bash
godoc -http=:6060
```

Then visit `http://localhost:6060/pkg/techwave/`

## 🧪 Testing

### Manual Testing

**Complete workflow example:**

```bash
# 1. Create an enrollment
ENROLLMENT_ID=$(curl -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{"student_id":"s1","course_id":"c1","status":"pending"}' \
  | jq -r '.id')

# 2. Get the enrollment
curl http://localhost:8080/api/enrollments/$ENROLLMENT_ID

# 3. Update status to active
curl -X PUT http://localhost:8080/api/enrollments/$ENROLLMENT_ID \
  -H "Content-Type: application/json" \
  -d '{"student_id":"s1","course_id":"c1","status":"active"}'

# 4. Update status to completed
curl -X PUT http://localhost:8080/api/enrollments/$ENROLLMENT_ID \
  -H "Content-Type: application/json" \
  -d '{"student_id":"s1","course_id":"c1","status":"completed"}'

# 5. List all enrollments
curl http://localhost:8080/api/enrollments

# 6. Delete the enrollment
curl -X DELETE http://localhost:8080/api/enrollments/$ENROLLMENT_ID
```

**PowerShell workflow:**

```powershell
# 1. Create enrollment
$createBody = @{student_id="s1"; course_id="c1"; status="pending"} | ConvertTo-Json
$enrollment = Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments" -Method Post -Body $createBody -ContentType "application/json"
$enrollmentId = $enrollment.id

# 2. Get enrollment
$enrollment = Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments/$enrollmentId" -Method Get

# 3. Update to active
$updateBody = @{student_id="s1"; course_id="c1"; status="active"} | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments/$enrollmentId" -Method Put -Body $updateBody -ContentType "application/json"

# 4. List all
$all = Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments" -Method Get

# 5. Delete
Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments/$enrollmentId" -Method Delete
```

### Validation Testing

**Test required fields:**
```bash
# Missing student_id
curl -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{"course_id":"c1","status":"active"}'
# Expected: 400 Bad Request - "student_id is required"

# Missing course_id
curl -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{"student_id":"s1","status":"active"}'
# Expected: 400 Bad Request - "course_id is required"

# Invalid status
curl -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{"student_id":"s1","course_id":"c1","status":"invalid"}'
# Expected: 400 Bad Request - "status must be one of: pending, active, completed"
```

### Running Unit Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test ./models
go test ./handlers
go test ./repository
```

## 🔧 Troubleshooting

### Common Issues

#### 1. Port Already in Use

**Symptom:** `listen tcp :8080: bind: address already in use`

**Solution:**
```bash
# Find process using port 8080
lsof -i :8080  # Mac/Linux
netstat -ano | findstr :8080  # Windows

# Kill the process
kill -9 <PID>  # Mac/Linux
taskkill /PID <PID> /F  # Windows
```

#### 2. Import Errors

**Symptom:** `package techwave/models is not in GOROOT`

**Solution:**
```bash
# Ensure you're in the project directory
cd <project-root>

# Download dependencies
go mod download

# Tidy up go.mod
go mod tidy
```

#### 3. JSON Parsing Errors

**Symptom:** `400 Bad Request - Invalid request payload`

**Causes:**
- Malformed JSON syntax
- Missing quotes around strings
- Trailing commas

**Solution:** Validate JSON before sending:
```bash
# Use jq to validate
echo '{"student_id":"s1","course_id":"c1","status":"active"}' | jq .
```

#### 4. Enrollment Not Found

**Symptom:** `404 Not Found - Enrollment not found`

**Causes:**
- Using incorrect enrollment ID
- Enrollment was deleted
- Server restarted (in-memory storage cleared)

**Solution:** 
- Verify enrollment exists: `GET /api/enrollments`
- Check the ID is correct (UUID format)

#### 5. Cannot Update Enrollment

**Symptom:** `400 Bad Request` when updating

**Solution:** Ensure all required fields are included in PUT request:
```json
{
  "student_id": "required",
  "course_id": "required",
  "status": "required"
}
```

### Debug Mode

Enable detailed logging by modifying `main.go`:

```go
import "log"

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    // ... rest of main
}
```

### Getting Help

- Check godoc: `godoc -http=:6060`
- Review code comments in source files
- Verify JSON payload format matches examples
- Check server logs for detailed error messages

## 📝 Business Rules

### Enrollment Status Lifecycle

1. **pending**: Student has enrolled but course hasn't started
   - Initial state for new enrollments
   - Can transition to: active, deleted

2. **active**: Student is currently taking the course
   - Indicates ongoing participation
   - Can transition to: completed, deleted

3. **completed**: Student has finished the course
   - Terminal state (no further transitions expected)
   - Can transition to: deleted

### Data Constraints

- `student_id`: Required, any non-empty string
- `course_id`: Required, any non-empty string
- `status`: Required, must be one of: "pending", "active", "completed"
- `enrollment_date`: Optional, defaults to current timestamp
- `id`: Auto-generated UUID, cannot be set by client
- `created_at`: Auto-set on creation
- `updated_at`: Auto-updated on modification

## 🚀 Production Considerations

### Current Limitations

- **In-Memory Storage**: Data is lost on server restart
- **No Authentication**: All endpoints are publicly accessible
- **No Rate Limiting**: API can be overwhelmed by requests
- **No Pagination**: Large datasets return all records

### Recommended Enhancements

1. **Add Database**: Replace in-memory storage with PostgreSQL/MySQL
2. **Add Authentication**: Implement JWT or OAuth2
3. **Add Rate Limiting**: Protect against abuse
4. **Add Pagination**: Support `?page=1&limit=20` query parameters
5. **Add Logging**: Structured logging with log levels
6. **Add Metrics**: Prometheus metrics for monitoring
7. **Add Docker**: Container-based deployment
8. **Add CI/CD**: Automated testing and deployment

## 📄 License

[Your License Here]

## 🤝 Contributing

[Your Contributing Guidelines Here]

---

**Built with ❤️ using Go and GitHub Copilot**