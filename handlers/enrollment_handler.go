// Package handlers provides HTTP request handlers for the Grade Management API.
// It implements RESTful endpoints for enrollment management with proper error handling.
package handlers

import (
	"encoding/json"
	"net/http"
	"techwave/models"
	"techwave/repository"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// EnrollmentHandler handles HTTP requests for enrollment operations.
// It provides CRUD functionality for student enrollments with validation and error handling.
type EnrollmentHandler struct {
	repo *repository.EnrollmentRepository
}

// NewEnrollmentHandler creates a new enrollment handler with the provided repository.
//
// Parameters:
//   - repo: The enrollment repository for data persistence
//
// Returns:
//   - A configured EnrollmentHandler instance
//
// Example:
//
//	repo := repository.NewEnrollmentRepository()
//	handler := handlers.NewEnrollmentHandler(repo)
func NewEnrollmentHandler(repo *repository.EnrollmentRepository) *EnrollmentHandler {
	return &EnrollmentHandler{
		repo: repo,
	}
}

// CreateEnrollment handles POST /api/enrollments to create a new enrollment.
// It validates the request, generates IDs and timestamps, and persists the enrollment.
//
// Request Body (JSON):
//
//	{
//	  "student_id": "student-123",     // Required
//	  "course_id": "course-456",       // Required
//	  "status": "active",              // Required: pending/active/completed
//	  "enrollment_date": "2024-01-15"  // Optional: defaults to now
//	}
//
// Response Codes:
//   - 201 Created: Enrollment created successfully
//   - 400 Bad Request: Invalid payload or validation failure
//   - 409 Conflict: Enrollment with same ID already exists
//   - 500 Internal Server Error: Database or system error
//
// Example using curl:
//
//	curl -X POST http://localhost:8080/api/enrollments \
//	  -H "Content-Type: application/json" \
//	  -d '{"student_id":"s1","course_id":"c1","status":"active"}'
func (h *EnrollmentHandler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	var enrollment models.Enrollment

	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the enrollment
	if err := enrollment.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Set timestamps and generate ID
	enrollment.ID = uuid.New().String()
	enrollment.CreatedAt = time.Now()
	enrollment.UpdatedAt = time.Now()

	// Set enrollment date if not provided
	if enrollment.EnrollmentDate.IsZero() {
		enrollment.EnrollmentDate = time.Now()
	}

	// Create the enrollment
	if err := h.repo.Create(&enrollment); err != nil {
		if err == repository.ErrAlreadyExists {
			respondWithError(w, http.StatusConflict, "Enrollment already exists")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create enrollment")
		return
	}

	respondWithJSON(w, http.StatusCreated, enrollment)
}

// GetEnrollment handles GET /api/enrollments/{id} to retrieve a specific enrollment.
//
// URL Parameters:
//   - id: The enrollment UUID
//
// Response Codes:
//   - 200 OK: Enrollment found and returned
//   - 404 Not Found: Enrollment with specified ID doesn't exist
//   - 500 Internal Server Error: Database or system error
//
// Example using curl:
//
//	curl http://localhost:8080/api/enrollments/123e4567-e89b-12d3-a456-426614174000
//
// Example using PowerShell:
//
//	Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments/{id}" -Method Get
func (h *EnrollmentHandler) GetEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	enrollment, err := h.repo.GetByID(id)
	if err != nil {
		if err == repository.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Enrollment not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve enrollment")
		return
	}

	respondWithJSON(w, http.StatusOK, enrollment)
}

// GetAllEnrollments handles GET /api/enrollments to retrieve all enrollments.
// Returns an array of all enrollment records in the system.
//
// Response Codes:
//   - 200 OK: Returns array of enrollments (may be empty)
//
// Example using curl:
//
//	curl http://localhost:8080/api/enrollments
//
// Example using PowerShell:
//
//	Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments" -Method Get
func (h *EnrollmentHandler) GetAllEnrollments(w http.ResponseWriter, r *http.Request) {
	enrollments := h.repo.GetAll()
	respondWithJSON(w, http.StatusOK, enrollments)
}

// UpdateEnrollment handles PUT /api/enrollments/{id} to update an existing enrollment.
// The entire enrollment resource is replaced with the new data.
//
// URL Parameters:
//   - id: The enrollment UUID to update
//
// Request Body (JSON):
//
//	{
//	  "student_id": "student-123",
//	  "course_id": "course-456",
//	  "status": "completed",
//	  "enrollment_date": "2024-01-15"
//	}
//
// Response Codes:
//   - 200 OK: Enrollment updated successfully
//   - 400 Bad Request: Invalid payload or validation failure
//   - 404 Not Found: Enrollment with specified ID doesn't exist
//   - 500 Internal Server Error: Database or system error
//
// Example using curl:
//
//	curl -X PUT http://localhost:8080/api/enrollments/{id} \
//	  -H "Content-Type: application/json" \
//	  -d '{"student_id":"s1","course_id":"c1","status":"completed"}'
//
// Example using PowerShell:
//
//	$body = @{student_id="s1"; course_id="c1"; status="completed"} | ConvertTo-Json
//	Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments/{id}" -Method Put -Body $body -ContentType "application/json"
func (h *EnrollmentHandler) UpdateEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var enrollment models.Enrollment
	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the enrollment
	if err := enrollment.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update timestamp and set ID
	enrollment.ID = id
	enrollment.UpdatedAt = time.Now()

	// Update the enrollment
	if err := h.repo.Update(id, &enrollment); err != nil {
		if err == repository.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Enrollment not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to update enrollment")
		return
	}

	respondWithJSON(w, http.StatusOK, enrollment)
}

// DeleteEnrollment handles DELETE /api/enrollments/{id} to remove an enrollment.
// This is a hard delete operation and cannot be undone.
//
// URL Parameters:
//   - id: The enrollment UUID to delete
//
// Response Codes:
//   - 200 OK: Enrollment deleted successfully
//   - 404 Not Found: Enrollment with specified ID doesn't exist
//   - 500 Internal Server Error: Database or system error
//
// Example using curl:
//
//	curl -X DELETE http://localhost:8080/api/enrollments/{id}
//
// Example using PowerShell:
//
//	Invoke-RestMethod -Uri "http://localhost:8080/api/enrollments/{id}" -Method Delete
func (h *EnrollmentHandler) DeleteEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.repo.Delete(id); err != nil {
		if err == repository.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Enrollment not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to delete enrollment")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Enrollment deleted successfully"})
}

// respondWithError sends a standardized JSON error response.
//
// Parameters:
//   - w: HTTP response writer
//   - code: HTTP status code
//   - message: Error message to return to client
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON marshals the payload to JSON and sends it as HTTP response.
// Sets appropriate Content-Type header and handles marshaling errors.
//
// Parameters:
//   - w: HTTP response writer
//   - code: HTTP status code
//   - payload: Data to serialize as JSON
//
// Error handling:
//   - If JSON marshaling fails, returns 500 with generic error message
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
