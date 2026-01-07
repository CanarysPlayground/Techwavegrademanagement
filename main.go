package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Enrollment represents a student's enrollment in a course
// This model tracks the relationship between students and courses with status tracking
type Enrollment struct {
	ID             int       `json:"id"`              // Unique identifier for the enrollment
	StudentID      int       `json:"student_id"`      // ID of the enrolled student
	CourseID       int       `json:"course_id"`       // ID of the course
	EnrollmentDate time.Time `json:"enrollment_date"` // Date when the enrollment was created
	Status         string    `json:"status"`          // Current status: pending, active, or completed
}

// EnrollmentStore manages enrollment data in memory with thread-safe operations
type EnrollmentStore struct {
	mu          sync.RWMutex            // Mutex for thread-safe access
	enrollments map[int]*Enrollment     // In-memory storage of enrollments
	nextID      int                     // Counter for generating unique IDs
}

// Global store instance
var store = &EnrollmentStore{
	enrollments: make(map[int]*Enrollment),
	nextID:      1,
}

// CreateEnrollmentRequest represents the request body for creating an enrollment
type CreateEnrollmentRequest struct {
	StudentID int    `json:"student_id"`
	CourseID  int    `json:"course_id"`
	Status    string `json:"status"`
}

// UpdateEnrollmentRequest represents the request body for updating an enrollment
type UpdateEnrollmentRequest struct {
	StudentID *int    `json:"student_id,omitempty"`
	CourseID  *int    `json:"course_id,omitempty"`
	Status    *string `json:"status,omitempty"`
}

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// validateStatus checks if the provided status is valid
// Valid statuses are: pending, active, completed
func validateStatus(status string) bool {
	validStatuses := map[string]bool{
		"pending":   true,
		"active":    true,
		"completed": true,
	}
	return validStatuses[status]
}

// sendJSONResponse sends a JSON response with the given status code
func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// sendErrorResponse sends a JSON error response
func sendErrorResponse(w http.ResponseWriter, statusCode int, errorMsg string, message string) {
	sendJSONResponse(w, statusCode, ErrorResponse{
		Error:   errorMsg,
		Message: message,
	})
}

// CreateEnrollmentHandler handles POST /api/enrollments
// Creates a new enrollment with validation
func CreateEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateEnrollmentRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validate required fields
	if req.StudentID <= 0 {
		sendErrorResponse(w, http.StatusBadRequest, "Validation error", "student_id must be a positive integer")
		return
	}
	if req.CourseID <= 0 {
		sendErrorResponse(w, http.StatusBadRequest, "Validation error", "course_id must be a positive integer")
		return
	}
	if req.Status == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Validation error", "status is required")
		return
	}

	// Validate status
	if !validateStatus(req.Status) {
		sendErrorResponse(w, http.StatusBadRequest, "Validation error", "status must be one of: pending, active, completed")
		return
	}

	// Create enrollment
	store.mu.Lock()
	defer store.mu.Unlock()

	enrollment := &Enrollment{
		ID:             store.nextID,
		StudentID:      req.StudentID,
		CourseID:       req.CourseID,
		EnrollmentDate: time.Now(),
		Status:         req.Status,
	}

	store.enrollments[store.nextID] = enrollment
	store.nextID++

	sendJSONResponse(w, http.StatusCreated, enrollment)
}

// GetAllEnrollmentsHandler handles GET /api/enrollments
// Returns all enrollments in the system
func GetAllEnrollmentsHandler(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	// Convert map to slice for JSON response
	enrollments := make([]*Enrollment, 0, len(store.enrollments))
	for _, enrollment := range store.enrollments {
		enrollments = append(enrollments, enrollment)
	}

	sendJSONResponse(w, http.StatusOK, enrollments)
}

// GetEnrollmentByIDHandler handles GET /api/enrollments/{id}
// Returns a specific enrollment by ID
func GetEnrollmentByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid ID", "ID must be a valid integer")
		return
	}

	// Find enrollment
	store.mu.RLock()
	enrollment, exists := store.enrollments[id]
	store.mu.RUnlock()

	if !exists {
		sendErrorResponse(w, http.StatusNotFound, "Not found", fmt.Sprintf("Enrollment with ID %d not found", id))
		return
	}

	sendJSONResponse(w, http.StatusOK, enrollment)
}

// UpdateEnrollmentHandler handles PUT /api/enrollments/{id}
// Updates an existing enrollment
func UpdateEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid ID", "ID must be a valid integer")
		return
	}

	// Parse request body
	var req UpdateEnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validate status if provided
	if req.Status != nil && !validateStatus(*req.Status) {
		sendErrorResponse(w, http.StatusBadRequest, "Validation error", "status must be one of: pending, active, completed")
		return
	}

	// Validate IDs if provided
	if req.StudentID != nil && *req.StudentID <= 0 {
		sendErrorResponse(w, http.StatusBadRequest, "Validation error", "student_id must be a positive integer")
		return
	}
	if req.CourseID != nil && *req.CourseID <= 0 {
		sendErrorResponse(w, http.StatusBadRequest, "Validation error", "course_id must be a positive integer")
		return
	}

	// Update enrollment
	store.mu.Lock()
	defer store.mu.Unlock()

	enrollment, exists := store.enrollments[id]
	if !exists {
		sendErrorResponse(w, http.StatusNotFound, "Not found", fmt.Sprintf("Enrollment with ID %d not found", id))
		return
	}

	// Apply updates
	if req.StudentID != nil {
		enrollment.StudentID = *req.StudentID
	}
	if req.CourseID != nil {
		enrollment.CourseID = *req.CourseID
	}
	if req.Status != nil {
		enrollment.Status = *req.Status
	}

	sendJSONResponse(w, http.StatusOK, enrollment)
}

// DeleteEnrollmentHandler handles DELETE /api/enrollments/{id}
// Deletes an enrollment by ID
func DeleteEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid ID", "ID must be a valid integer")
		return
	}

	// Delete enrollment
	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.enrollments[id]; !exists {
		sendErrorResponse(w, http.StatusNotFound, "Not found", fmt.Sprintf("Enrollment with ID %d not found", id))
		return
	}

	delete(store.enrollments, id)

	sendJSONResponse(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Enrollment with ID %d deleted successfully", id),
	})
}

// HealthCheckHandler handles GET /
// Returns a simple health check response
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	sendJSONResponse(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"message": "Grade Management API - Ready for AI delegation!",
	})
}

func main() {
	// Create a new router
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/", HealthCheckHandler).Methods("GET")

	// API routes with /api prefix
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Enrollment CRUD endpoints
	apiRouter.HandleFunc("/enrollments", CreateEnrollmentHandler).Methods("POST")
	apiRouter.HandleFunc("/enrollments", GetAllEnrollmentsHandler).Methods("GET")
	apiRouter.HandleFunc("/enrollments/{id}", GetEnrollmentByIDHandler).Methods("GET")
	apiRouter.HandleFunc("/enrollments/{id}", UpdateEnrollmentHandler).Methods("PUT")
	apiRouter.HandleFunc("/enrollments/{id}", DeleteEnrollmentHandler).Methods("DELETE")

	// Start server
	port := ":8080"
	fmt.Printf("Starting Grade Management API on port %s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET    /                      - Health check")
	fmt.Println("  POST   /api/enrollments       - Create enrollment")
	fmt.Println("  GET    /api/enrollments       - List all enrollments")
	fmt.Println("  GET    /api/enrollments/{id}  - Get enrollment by ID")
	fmt.Println("  PUT    /api/enrollments/{id}  - Update enrollment")
	fmt.Println("  DELETE /api/enrollments/{id}  - Delete enrollment")
	log.Fatal(http.ListenAndServe(port, router))
}