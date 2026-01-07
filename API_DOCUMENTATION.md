# API Documentation - Grade Management System

Complete technical documentation for developers integrating with or maintaining the Grade Management API.

## Quick Reference

| Endpoint | Method | Purpose | Status Codes |
|----------|--------|---------|--------------|
| `/api/enrollments` | POST | Create enrollment | 201, 400, 409, 500 |
| `/api/enrollments` | GET | List all enrollments | 200 |
| `/api/enrollments/{id}` | GET | Get specific enrollment | 200, 404, 500 |
| `/api/enrollments/{id}` | PUT | Update enrollment | 200, 400, 404, 500 |
| `/api/enrollments/{id}` | DELETE | Delete enrollment | 200, 404, 500 |

## Architecture Overview

### Package Structure

```
techwave/
├── main.go              # HTTP server and routing setup
├── models/              # Domain models and business logic
│   └── enrollment.go
├── handlers/            # HTTP request handlers
│   └── enrollment_handler.go
└── repository/          # Data storage layer
    └── enrollment_repository.go
```

### Design Patterns

- **Repository Pattern**: Separates data access logic
- **Handler Pattern**: Decouples HTTP concerns from business logic
- **Dependency Injection**: Handlers receive repository instances
- **Thread-Safe Operations**: RWMutex for concurrent access

### Data Flow

```
HTTP Request
    ↓
Handler (validation, serialization)
    ↓
Model (business logic validation)
    ↓
Repository (data storage)
    ↓
Response (JSON)
```

## Domain Model

### Enrollment

Represents a student's enrollment in a course.

**Structure:**
```go
type Enrollment struct {
    ID             string    `json:"id"`              // UUID
    StudentID      string    `json:"student_id"`      // Required
    CourseID       string    `json:"course_id"`       // Required
    EnrollmentDate time.Time `json:"enrollment_date"` // Auto-set
    Status         string    `json:"status"`          // Required: pending|active|completed
    CreatedAt      time.Time `json:"created_at"`      // Auto-set
    UpdatedAt      time.Time `json:"updated_at"`      // Auto-updated
}
```

**Business Rules:**
1. `student_id` and `course_id` are required
2. `status` must be one of: "pending", "active", "completed"
3. `id` is auto-generated (UUID v4)
4. `enrollment_date` defaults to current time if not provided
5. `created_at` set once on creation
6. `updated_at` updated on every modification

**Validation:**
```go
func (e *Enrollment) Validate() error
```

Returns error if:
- `student_id` is empty
- `course_id` is empty
- `status` is empty or not in valid set

### Status Lifecycle

```
   pending ──────┐
      │          │
      ▼          │
   active ───────┤
      │          │
      ▼          ▼
  completed    deleted
```

## API Endpoints

### 1. Create Enrollment

**Endpoint:** `POST /api/enrollments`

**Purpose:** Create a new student enrollment

**Request Headers:**
- `Content-Type: application/json`

**Request Body:**
```json
{
  "student_id": "string (required)",
  "course_id": "string (required)",
  "status": "string (required: pending|active|completed)",
  "enrollment_date": "string (optional: ISO 8601 timestamp)"
}
```

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

**Error Responses:**
- `400 Bad Request`: Invalid JSON or validation failure
  ```json
  {"error": "student_id is required"}
  {"error": "status must be one of: pending, active, completed"}
  ```
- `409 Conflict`: Enrollment with same ID already exists
  ```json
  {"error": "Enrollment already exists"}
  ```
- `500 Internal Server Error`: Server-side error
  ```json
  {"error": "Failed to create enrollment"}
  ```

**Code Example (Go):**
```go
enrollment := &models.Enrollment{
    StudentID: "student-123",
    CourseID: "course-456",
    Status: "active",
}
if err := enrollment.Validate(); err != nil {
    // Handle validation error
}
```

**curl Example:**
```bash
curl -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "student-123",
    "course_id": "course-456",
    "status": "active"
  }'
```

**PowerShell Example:**
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

---

### 2. Get All Enrollments

**Endpoint:** `GET /api/enrollments`

**Purpose:** Retrieve all enrollments

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

**Notes:**
- Returns empty array `[]` if no enrollments exist
- No pagination (returns all records)
- Thread-safe read operation

**curl Example:**
```bash
curl http://localhost:8080/api/enrollments
```

**PowerShell Example:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments" -Method Get
```

---

### 3. Get Enrollment by ID

**Endpoint:** `GET /api/enrollments/{id}`

**Purpose:** Retrieve a specific enrollment

**URL Parameters:**
- `id` (required): Enrollment UUID

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

**Error Responses:**
- `404 Not Found`: Enrollment doesn't exist
  ```json
  {"error": "Enrollment not found"}
  ```
- `500 Internal Server Error`: Server-side error
  ```json
  {"error": "Failed to retrieve enrollment"}
  ```

**curl Example:**
```bash
curl http://localhost:8080/api/enrollments/123e4567-e89b-12d3-a456-426614174000
```

**PowerShell Example:**
```powershell
$id = "123e4567-e89b-12d3-a456-426614174000"
Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments/$id" -Method Get
```

---

### 4. Update Enrollment

**Endpoint:** `PUT /api/enrollments/{id}`

**Purpose:** Update an existing enrollment (full replacement)

**URL Parameters:**
- `id` (required): Enrollment UUID

**Request Headers:**
- `Content-Type: application/json`

**Request Body:**
```json
{
  "student_id": "string (required)",
  "course_id": "string (required)",
  "status": "string (required: pending|active|completed)",
  "enrollment_date": "string (optional: ISO 8601 timestamp)"
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

**Error Responses:**
- `400 Bad Request`: Invalid JSON or validation failure
- `404 Not Found`: Enrollment doesn't exist
- `500 Internal Server Error`: Server-side error

**Notes:**
- All fields must be provided (PUT replaces entire resource)
- `id` in URL takes precedence over any ID in body
- `updated_at` is automatically set to current time
- `created_at` remains unchanged

**curl Example:**
```bash
curl -X PUT http://localhost:8080/api/enrollments/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "student-123",
    "course_id": "course-456",
    "status": "completed"
  }'
```

**PowerShell Example:**
```powershell
$id = "123e4567-e89b-12d3-a456-426614174000"
$body = @{
    student_id = "student-123"
    course_id = "course-456"
    status = "completed"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments/$id" `
  -Method Put `
  -Body $body `
  -ContentType "application/json"
```

---

### 5. Delete Enrollment

**Endpoint:** `DELETE /api/enrollments/{id}`

**Purpose:** Permanently delete an enrollment

**URL Parameters:**
- `id` (required): Enrollment UUID

**Response (200 OK):**
```json
{
  "message": "Enrollment deleted successfully"
}
```

**Error Responses:**
- `404 Not Found`: Enrollment doesn't exist
- `500 Internal Server Error`: Server-side error

**Notes:**
- This is a hard delete (permanent)
- Cannot be undone
- Idempotent (deleting non-existent resource returns 404)

**curl Example:**
```bash
curl -X DELETE http://localhost:8080/api/enrollments/123e4567-e89b-12d3-a456-426614174000
```

**PowerShell Example:**
```powershell
$id = "123e4567-e89b-12d3-a456-426614174000"
Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments/$id" -Method Delete
```

## Repository Layer

### EnrollmentRepository

Thread-safe in-memory storage for enrollments.

**Methods:**

#### Create
```go
func (r *EnrollmentRepository) Create(enrollment *models.Enrollment) error
```
- Adds new enrollment
- Returns `ErrAlreadyExists` if ID exists
- Thread-safe with write lock

#### GetByID
```go
func (r *EnrollmentRepository) GetByID(id string) (*models.Enrollment, error)
```
- Retrieves enrollment by ID
- Returns `ErrNotFound` if not exists
- Thread-safe with read lock

#### GetAll
```go
func (r *EnrollmentRepository) GetAll() []*models.Enrollment
```
- Returns all enrollments
- Returns empty slice if none exist
- Thread-safe with read lock

#### Update
```go
func (r *EnrollmentRepository) Update(id string, enrollment *models.Enrollment) error
```
- Updates existing enrollment
- Returns `ErrNotFound` if not exists
- Creates copy to prevent external modification
- Thread-safe with write lock

#### Delete
```go
func (r *EnrollmentRepository) Delete(id string) error
```
- Removes enrollment
- Returns `ErrNotFound` if not exists
- Thread-safe with write lock

**Usage Example:**
```go
repo := repository.NewEnrollmentRepository()

enrollment := &models.Enrollment{
    ID: uuid.New().String(),
    StudentID: "student-123",
    CourseID: "course-456",
    Status: "active",
}

if err := repo.Create(enrollment); err != nil {
    if err == repository.ErrAlreadyExists {
        // Handle duplicate
    }
}

found, err := repo.GetByID(enrollment.ID)
if err == repository.ErrNotFound {
    // Handle not found
}
```

## Error Handling

### Standard Error Format

All errors return JSON:
```json
{
  "error": "Human-readable error message"
}
```

### Repository Errors

```go
var (
    ErrNotFound = errors.New("enrollment not found")
    ErrAlreadyExists = errors.New("enrollment already exists")
)
```

### HTTP Status Code Mapping

| Error Condition | HTTP Status | Example |
|----------------|-------------|---------|
| Invalid JSON | 400 | Malformed request body |
| Validation failure | 400 | Missing required field |
| Resource not found | 404 | GET non-existent ID |
| Duplicate resource | 409 | Create with existing ID |
| Database error | 500 | Storage failure |
| JSON marshal error | 500 | Response serialization failure |

## Performance Considerations

### Thread Safety

- All repository operations use RWMutex
- Read operations (`GetByID`, `GetAll`) use read locks (concurrent)
- Write operations (`Create`, `Update`, `Delete`) use write locks (exclusive)

### Memory Usage

- In-memory storage: O(n) space where n = number of enrollments
- Each enrollment ~300 bytes
- No automatic cleanup (grows until server restart)

### Scalability Limitations

1. **No Persistence**: Data lost on restart
2. **No Pagination**: All records returned in GetAll
3. **No Filtering**: Cannot query by student_id or course_id
4. **No Caching**: Direct repository access every time

### Production Recommendations

1. Replace in-memory repo with database (PostgreSQL, MySQL)
2. Add pagination (`?page=1&limit=20`)
3. Add filtering (`?student_id=xyz&status=active`)
4. Add caching layer (Redis)
5. Add connection pooling
6. Add request rate limiting
7. Add request logging and metrics

## Security Considerations

### Current State (Development)

⚠️ **No authentication or authorization**
- All endpoints publicly accessible
- No user identity tracking
- No access control

### Production Requirements

1. **Authentication**: JWT tokens or OAuth2
2. **Authorization**: Role-based access control
   - Students can only view their own enrollments
   - Instructors can view enrollments in their courses
   - Admins have full access
3. **Input Validation**: Already implemented
4. **Rate Limiting**: Prevent abuse
5. **HTTPS**: Encrypt data in transit
6. **CORS**: Configure allowed origins
7. **SQL Injection**: Not applicable (no SQL), but prepare for DB migration

## Testing

### Unit Testing Example

```go
func TestEnrollmentValidation(t *testing.T) {
    tests := []struct {
        name        string
        enrollment  models.Enrollment
        expectError bool
    }{
        {
            name: "valid enrollment",
            enrollment: models.Enrollment{
                StudentID: "s1",
                CourseID: "c1",
                Status: "active",
            },
            expectError: false,
        },
        {
            name: "missing student_id",
            enrollment: models.Enrollment{
                CourseID: "c1",
                Status: "active",
            },
            expectError: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.enrollment.Validate()
            if (err != nil) != tt.expectError {
                t.Errorf("Validate() error = %v, expectError %v", err, tt.expectError)
            }
        })
    }
}
```

### Integration Testing

See `TESTING.md` for comprehensive test scenarios and examples.

## Godoc Comments

All public functions and types include godoc comments. Generate documentation:

```bash
# Start godoc server
godoc -http=:6060

# Visit in browser
http://localhost:6060/pkg/techwave/
```

## Migration Path

### From In-Memory to PostgreSQL

1. **Install driver:**
   ```bash
   go get github.com/lib/pq
   ```

2. **Update repository:**
   ```go
   type EnrollmentRepository struct {
       db *sql.DB
   }
   
   func (r *EnrollmentRepository) Create(e *models.Enrollment) error {
       query := `INSERT INTO enrollments (...) VALUES (...)`
       _, err := r.db.Exec(query, e.ID, e.StudentID, ...)
       return err
   }
   ```

3. **Update initialization:**
   ```go
   db, err := sql.Open("postgres", connStr)
   repo := repository.NewEnrollmentRepository(db)
   ```

No changes needed to handlers or models!

## Common Integration Patterns

### Pattern 1: Create and Track

```bash
# Create enrollment
ENROLLMENT_ID=$(curl -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{"student_id":"s1","course_id":"c1","status":"pending"}' \
  | jq -r '.id')

# Poll until active
while true; do
  STATUS=$(curl -s http://localhost:8080/api/enrollments/$ENROLLMENT_ID | jq -r '.status')
  [ "$STATUS" = "active" ] && break
  sleep 5
done
```

### Pattern 2: Batch Operations

```bash
# Create multiple enrollments
for student in s1 s2 s3; do
  curl -X POST http://localhost:8080/api/enrollments \
    -H "Content-Type: application/json" \
    -d "{\"student_id\":\"$student\",\"course_id\":\"c1\",\"status\":\"pending\"}"
done

# Get all and filter by course
curl -s http://localhost:8080/api/enrollments \
  | jq '.[] | select(.course_id == "c1")'
```

### Pattern 3: Status Progression

```go
// Helper function to progress enrollment status
func progressEnrollment(client *http.Client, id string, status string) error {
    // Get current enrollment
    resp, err := client.Get(fmt.Sprintf("http://localhost:8080/api/enrollments/%s", id))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    var enrollment models.Enrollment
    if err := json.NewDecoder(resp.Body).Decode(&enrollment); err != nil {
        return err
    }
    
    // Update status
    enrollment.Status = status
    body, _ := json.Marshal(enrollment)
    
    req, _ := http.NewRequest("PUT", 
        fmt.Sprintf("http://localhost:8080/api/enrollments/%s", id),
        bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    
    _, err = client.Do(req)
    return err
}
```

## Changelog

### Version 1.0.0
- Initial release
- CRUD operations for enrollments
- In-memory storage
- Input validation
- Thread-safe operations
- Comprehensive documentation

---

**For more information, see:**
- `README.md` - Setup and usage guide
- `TESTING.md` - Comprehensive testing guide
- Inline godoc comments in source files
