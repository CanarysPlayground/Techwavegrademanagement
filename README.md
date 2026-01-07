# Grade Management API - Enrollment Service

## 🚀 Overview

A high-performance student enrollment management API built with Go, featuring Redis caching for sub-100ms response times and comprehensive API contract validation.

### ✨ Features

- ✅ **Complete CRUD Operations** - Create, Read, Update, Delete enrollments
- ⚡ **Redis Caching** - 5-minute TTL with cache-aside pattern
- 🔍 **X-Cache-Status Headers** - Debug cache hits/misses in real-time
- 🛡️ **API Contract Validation** - OpenAPI 3.0 spec with automated validation
- 🧪 **Integration Test Suite** - 100% pass rate with performance assertions
- 🔄 **Graceful Degradation** - Works with or without Redis
- 🐳 **Docker Ready** - Containerized Redis setup

## 📋 Prerequisites

- Go 1.22.5 or higher
- Docker & Docker Compose (for Redis)
- Git

## 🚀 Quick Start

### 1. Clone and Setup

```bash
git clone https://github.com/CanarysPlayground/Techwavegrademanagement.git
cd Techwavegrademanagement
go mod download
```

### 2. Start Redis

```bash
docker-compose up -d
```

### 3. Run the Server

```bash
go run main.go
```

Server starts on `http://localhost:8080`

### 4. Test the API

```powershell
# Create an enrollment
Invoke-RestMethod -Uri http://localhost:8080/api/enrollments -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"student_id":"42","course_id":"101","status":"pending"}'

# Get all enrollments
Invoke-RestMethod -Uri http://localhost:8080/api/enrollments

# Get specific enrollment (check X-Cache-Status header)
$response = Invoke-WebRequest -Uri "http://localhost:8080/api/enrollments/<enrollment-id>"
$response.Headers["X-Cache-Status"]  # Returns HIT or MISS
```

## 📡 API Endpoints

| Method | Endpoint | Description | Cache Behavior |
|--------|----------|-------------|----------------|
| GET | `/` | Root endpoint | N/A |
| GET | `/health` | Health check | N/A |
| POST | `/api/enrollments` | Create enrollment | No cache |
| GET | `/api/enrollments` | List all enrollments | No cache |
| GET | `/api/enrollments/{id}` | Get enrollment | Cached (5 min TTL) |
| PUT | `/api/enrollments/{id}` | Update enrollment | Invalidates cache |
| DELETE | `/api/enrollments/{id}` | Delete enrollment | Invalidates cache |

### Request/Response Examples

See the complete OpenAPI specification in [api/openapi.yaml](api/openapi.yaml) for detailed schemas and examples.

## 🧪 Testing

### Contract Validation

Validates API implementation against OpenAPI spec:

```bash
go run -tags contract scripts/validate_contract.go
```

**Expected Output:**
```
✓ OpenAPI specification is valid
✓ Route validated: GET http://localhost:8080/
✓ Route validated: POST http://localhost:8080/api/enrollments
... (all routes validated)
✅ CONTRACT VALIDATION PASSED: All checks successful
```

### Integration Tests

Comprehensive test suite with cache behavior validation:

```bash
go test -tags integration -v ./tests/integration_test.go
```

**Test Coverage:**
- ✅ Complete CRUD workflow
- ✅ Cache hit/miss/invalidation behavior
- ✅ Performance assertions (<100ms cached responses)
- ✅ Validation error handling
- ✅ 404 error scenarios
- ✅ Response schema validation

**Expected Results:**
```
=== RUN   TestCompleteCRUDWorkflow
--- PASS: TestCompleteCRUDWorkflow (0.02s)
=== RUN   TestCachePerformance
    ✓ Cached response time: 1 ms
--- PASS: TestCachePerformance (0.01s)
... (all tests passing)
PASS
```

## 🏗️ Architecture

### Components

```
├── main.go                    # Application entry point with Redis setup
├── api/
│   └── openapi.yaml           # OpenAPI 3.0 specification
├── cache/
│   └── enrollment_cache.go    # Redis caching layer (5-min TTL)
├── handlers/
│   └── enrollment_handler.go  # HTTP request handlers with cache integration
├── models/
│   └── enrollment.go          # Enrollment data model and validation
├── repository/
│   └── enrollment_repository.go # In-memory data storage
├── middleware/
│   └── cache_middleware.go    # X-Cache-Status header middleware
├── scripts/
│   └── validate_contract.go   # Contract validation script
└── tests/
    └── integration_test.go    # Integration test suite
```

### Caching Strategy

**Cache-Aside Pattern:**
1. Check cache on GET requests
2. On cache MISS: fetch from DB, store in cache
3. On cache HIT: return cached data (<100ms)
4. Invalidate cache on UPDATE/DELETE operations

**Cache Headers:**
- `X-Cache-Status: HIT` - Served from Redis cache
- `X-Cache-Status: MISS` - Fetched from database and cached
- `X-Cache-Status: SKIP` - Caching disabled/not applicable

## 🔧 Configuration

Environment variables:

```bash
REDIS_ADDR=localhost:6379      # Redis server address (default: localhost:6379)
REDIS_PASSWORD=                # Redis password (optional)
```

## 🚀 CI/CD Integration

### GitHub Actions Example

```yaml
name: API Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      redis:
        image: redis:7-alpine
        ports:
          - 6379:6379
    
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22.5'
      
      - name: Install dependencies
        run: go mod download
      
      - name: Contract validation
        run: go run -tags contract scripts/validate_contract.go
      
      - name: Integration tests
        run: go test -tags integration -v ./tests/integration_test.go
```

## 📊 Performance Benchmarks

| Operation | Response Time | Cache Status |
|-----------|--------------|--------------|
| Create Enrollment | ~5-10ms | N/A |
| Get Enrollment (cached) | **<1ms** | HIT |
| Get Enrollment (uncached) | ~5-10ms | MISS |
| Update Enrollment | ~5-10ms | Invalidates |
| Delete Enrollment | ~5-10ms | Invalidates |

## 🤝 Development Workflow

1. **Feature Development**
   - Update [api/openapi.yaml](api/openapi.yaml) with new endpoints
   - Implement handlers with cache support
   - Run contract validation
   - Write integration tests

2. **Testing**
   ```bash
   # Validate contract
   go run -tags contract scripts/validate_contract.go
   
   # Run tests
   go test -tags integration -v ./tests/integration_test.go
   ```

3. **Commit**
   ```bash
   git add .
   git commit -m "feat: Add new endpoint with caching"
   git push
   ```

## 📝 API Documentation

Full API documentation is available in [api/openapi.yaml](api/openapi.yaml). You can view it using:

- [Swagger Editor](https://editor.swagger.io/) - Paste the YAML content
- [Redoc](https://redocly.github.io/redoc/) - For beautiful docs rendering

## 🐛 Troubleshooting

### Redis Connection Issues

```bash
# Check Redis is running
docker ps | grep redis

# View Redis logs
docker-compose logs redis

# Restart Redis
docker-compose restart redis
```

### Cache Not Working

- Check `X-Cache-Status` header in responses
- Verify Redis connection in server logs: `✓ Redis connection established`
- API works in degraded mode without Redis (cache disabled)

## 📄 License

MIT License - See LICENSE file for details

## 🎯 Project Status

**TEC-17**: ✅ Redis caching implementation - **COMPLETE**
**TEC-18**: ✅ API contract validation & test suite - **COMPLETE**

### Completed Features
- ✅ Complete CRUD API for enrollments
- ✅ Redis caching with 5-minute TTL
- ✅ Cache invalidation on UPDATE/DELETE
- ✅ X-Cache-Status headers for debugging
- ✅ OpenAPI 3.0 specification
- ✅ Automated contract validation
- ✅ Comprehensive integration tests
- ✅ Performance benchmarks (<100ms cached)
- ✅ 100% test pass rate

---

**Built with GitHub Copilot** 🤖 | **Powered by Redis** ⚡ | **Validated by OpenAPI** 🛡️