// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"techwave/cache"
	"techwave/handlers"
	"techwave/models"
	"techwave/repository"

	"github.com/alicebob/miniredis/v2"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestServer creates a test server with mock Redis
func setupTestServer(t *testing.T) (*httptest.Server, *miniredis.Miniredis, *cache.EnrollmentCache) {
	// Start mini Redis
	mr, err := miniredis.Run()
	require.NoError(t, err)

	// Create Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Initialize components
	enrollmentRepo := repository.NewEnrollmentRepository()
	enrollmentCache := cache.NewEnrollmentCache(redisClient)
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentRepo, enrollmentCache)

	// Setup router
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Grade Management API - Cache: enabled")
	}).Methods("GET")

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/enrollments", enrollmentHandler.CreateEnrollment).Methods("POST")
	apiRouter.HandleFunc("/enrollments", enrollmentHandler.GetAllEnrollments).Methods("GET")
	apiRouter.HandleFunc("/enrollments/{id}", enrollmentHandler.GetEnrollment).Methods("GET")
	apiRouter.HandleFunc("/enrollments/{id}", enrollmentHandler.UpdateEnrollment).Methods("PUT")
	apiRouter.HandleFunc("/enrollments/{id}", enrollmentHandler.DeleteEnrollment).Methods("DELETE")

	server := httptest.NewServer(router)
	return server, mr, enrollmentCache
}

// TestCompleteCRUDWorkflow tests the complete CRUD workflow
func TestCompleteCRUDWorkflow(t *testing.T) {
	server, mr, _ := setupTestServer(t)
	defer server.Close()
	defer mr.Close()

	// 1. CREATE enrollment
	createPayload := map[string]interface{}{
		"student_id": "student-123",
		"course_id":  "course-456",
		"status":     "pending",
	}
	createBody, _ := json.Marshal(createPayload)

	resp, err := http.Post(server.URL+"/api/enrollments", "application/json", bytes.NewBuffer(createBody))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created models.Enrollment
	err = json.NewDecoder(resp.Body).Decode(&created)
	require.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "student-123", created.StudentID)
	assert.Equal(t, "course-456", created.CourseID)
	assert.Equal(t, "pending", created.Status)
	resp.Body.Close()

	enrollmentID := created.ID

	// 2. GET single enrollment (cache MISS)
	resp, err = http.Get(server.URL + "/api/enrollments/" + enrollmentID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "MISS", resp.Header.Get("X-Cache-Status"))

	var fetched models.Enrollment
	err = json.NewDecoder(resp.Body).Decode(&fetched)
	require.NoError(t, err)
	assert.Equal(t, enrollmentID, fetched.ID)
	resp.Body.Close()

	// 3. GET same enrollment (cache HIT)
	resp, err = http.Get(server.URL + "/api/enrollments/" + enrollmentID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "HIT", resp.Header.Get("X-Cache-Status"))
	resp.Body.Close()

	// 4. UPDATE enrollment (should invalidate cache)
	updatePayload := map[string]interface{}{
		"student_id": "student-123",
		"course_id":  "course-456",
		"status":     "active",
	}
	updateBody, _ := json.Marshal(updatePayload)

	req, _ := http.NewRequest(http.MethodPut, server.URL+"/api/enrollments/"+enrollmentID, bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updated models.Enrollment
	err = json.NewDecoder(resp.Body).Decode(&updated)
	require.NoError(t, err)
	assert.Equal(t, "active", updated.Status)
	resp.Body.Close()

	// 5. GET after update (cache MISS due to invalidation)
	resp, err = http.Get(server.URL + "/api/enrollments/" + enrollmentID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "MISS", resp.Header.Get("X-Cache-Status"))
	resp.Body.Close()

	// 6. GET all enrollments
	resp, err = http.Get(server.URL + "/api/enrollments")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var allEnrollments []models.Enrollment
	err = json.NewDecoder(resp.Body).Decode(&allEnrollments)
	require.NoError(t, err)
	assert.NotEmpty(t, allEnrollments)
	resp.Body.Close()

	// 7. DELETE enrollment (should invalidate cache)
	req, _ = http.NewRequest(http.MethodDelete, server.URL+"/api/enrollments/"+enrollmentID, nil)
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// 8. GET after delete (should return 404)
	resp, err = http.Get(server.URL + "/api/enrollments/" + enrollmentID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()
}

// TestCachePerformance validates cached response performance
func TestCachePerformance(t *testing.T) {
	server, mr, _ := setupTestServer(t)
	defer server.Close()
	defer mr.Close()

	// Create enrollment
	createPayload := map[string]interface{}{
		"student_id": "perf-student",
		"course_id":  "perf-course",
		"status":     "active",
	}
	createBody, _ := json.Marshal(createPayload)

	resp, err := http.Post(server.URL+"/api/enrollments", "application/json", bytes.NewBuffer(createBody))
	require.NoError(t, err)

	var created models.Enrollment
	json.NewDecoder(resp.Body).Decode(&created)
	resp.Body.Close()

	enrollmentID := created.ID

	// First GET (cache MISS) - prime the cache
	resp, _ = http.Get(server.URL + "/api/enrollments/" + enrollmentID)
	resp.Body.Close()

	// Second GET (cache HIT) - measure performance
	start := time.Now()
	resp, err = http.Get(server.URL + "/api/enrollments/" + enrollmentID)
	duration := time.Since(start)
	require.NoError(t, err)

	assert.Equal(t, "HIT", resp.Header.Get("X-Cache-Status"))
	assert.Less(t, duration.Milliseconds(), int64(100), "Cached response should be under 100ms")
	resp.Body.Close()

	t.Logf("✓ Cached response time: %v ms", duration.Milliseconds())
}

// TestCacheInvalidation validates cache invalidation on UPDATE and DELETE
func TestCacheInvalidation(t *testing.T) {
	server, mr, _ := setupTestServer(t)
	defer server.Close()
	defer mr.Close()

	// Create enrollment
	createPayload := map[string]interface{}{
		"student_id": "cache-test",
		"course_id":  "cache-course",
		"status":     "pending",
	}
	createBody, _ := json.Marshal(createPayload)

	resp, _ := http.Post(server.URL+"/api/enrollments", "application/json", bytes.NewBuffer(createBody))
	var created models.Enrollment
	json.NewDecoder(resp.Body).Decode(&created)
	resp.Body.Close()

	enrollmentID := created.ID

	// GET to populate cache
	resp, _ = http.Get(server.URL + "/api/enrollments/" + enrollmentID)
	assert.Equal(t, "MISS", resp.Header.Get("X-Cache-Status"))
	resp.Body.Close()

	// GET again to confirm cache HIT
	resp, _ = http.Get(server.URL + "/api/enrollments/" + enrollmentID)
	assert.Equal(t, "HIT", resp.Header.Get("X-Cache-Status"))
	resp.Body.Close()

	// UPDATE to invalidate cache
	updatePayload := map[string]interface{}{
		"student_id": "cache-test",
		"course_id":  "cache-course",
		"status":     "active",
	}
	updateBody, _ := json.Marshal(updatePayload)
	req, _ := http.NewRequest(http.MethodPut, server.URL+"/api/enrollments/"+enrollmentID, bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = http.DefaultClient.Do(req)
	resp.Body.Close()

	// GET after update should be cache MISS
	resp, _ = http.Get(server.URL + "/api/enrollments/" + enrollmentID)
	assert.Equal(t, "MISS", resp.Header.Get("X-Cache-Status"))
	resp.Body.Close()

	// GET again to populate cache
	resp, _ = http.Get(server.URL + "/api/enrollments/" + enrollmentID)
	assert.Equal(t, "HIT", resp.Header.Get("X-Cache-Status"))
	resp.Body.Close()

	// DELETE to invalidate cache
	req, _ = http.NewRequest(http.MethodDelete, server.URL+"/api/enrollments/"+enrollmentID, nil)
	resp, _ = http.DefaultClient.Do(req)
	resp.Body.Close()

	// GET after delete should return 404
	resp, _ = http.Get(server.URL + "/api/enrollments/" + enrollmentID)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()
}

// TestValidationErrors tests error handling and validation
func TestValidationErrors(t *testing.T) {
	server, mr, _ := setupTestServer(t)
	defer server.Close()
	defer mr.Close()

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "missing student_id",
			payload:        map[string]interface{}{"course_id": "course-1", "status": "pending"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "student_id is required",
		},
		{
			name:           "missing course_id",
			payload:        map[string]interface{}{"student_id": "student-1", "status": "pending"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "course_id is required",
		},
		{
			name:           "missing status",
			payload:        map[string]interface{}{"student_id": "student-1", "course_id": "course-1"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "status is required",
		},
		{
			name:           "invalid status",
			payload:        map[string]interface{}{"student_id": "student-1", "course_id": "course-1", "status": "invalid"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "status must be one of: pending, active, completed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			resp, err := http.Post(server.URL+"/api/enrollments", "application/json", bytes.NewBuffer(body))
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var errorResp map[string]string
			json.NewDecoder(resp.Body).Decode(&errorResp)
			assert.Contains(t, errorResp["error"], tt.expectedError)
			resp.Body.Close()
		})
	}
}

// TestNotFoundErrors tests 404 error handling
func TestNotFoundErrors(t *testing.T) {
	server, mr, _ := setupTestServer(t)
	defer server.Close()
	defer mr.Close()

	nonExistentID := "00000000-0000-0000-0000-000000000000"

	// GET non-existent enrollment
	resp, err := http.Get(server.URL + "/api/enrollments/" + nonExistentID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()

	// UPDATE non-existent enrollment
	updatePayload := map[string]interface{}{
		"student_id": "student",
		"course_id":  "course",
		"status":     "active",
	}
	updateBody, _ := json.Marshal(updatePayload)
	req, _ := http.NewRequest(http.MethodPut, server.URL+"/api/enrollments/"+nonExistentID, bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()

	// DELETE non-existent enrollment
	req, _ = http.NewRequest(http.MethodDelete, server.URL+"/api/enrollments/"+nonExistentID, nil)
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()
}

// TestResponseSchemaValidation ensures all responses match expected schemas
func TestResponseSchemaValidation(t *testing.T) {
	server, mr, _ := setupTestServer(t)
	defer server.Close()
	defer mr.Close()

	// Create enrollment and validate response structure
	createPayload := map[string]interface{}{
		"student_id": "schema-student",
		"course_id":  "schema-course",
		"status":     "pending",
	}
	createBody, _ := json.Marshal(createPayload)

	resp, _ := http.Post(server.URL+"/api/enrollments", "application/json", bytes.NewBuffer(createBody))
	
	var enrollment models.Enrollment
	err := json.NewDecoder(resp.Body).Decode(&enrollment)
	require.NoError(t, err)
	
	// Validate all required fields are present
	assert.NotEmpty(t, enrollment.ID)
	assert.NotEmpty(t, enrollment.StudentID)
	assert.NotEmpty(t, enrollment.CourseID)
	assert.NotEmpty(t, enrollment.Status)
	assert.NotZero(t, enrollment.CreatedAt)
	assert.NotZero(t, enrollment.UpdatedAt)
	assert.NotZero(t, enrollment.EnrollmentDate)
	
	resp.Body.Close()
}
