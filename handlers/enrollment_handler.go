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

// EnrollmentHandler handles HTTP requests for enrollments
type EnrollmentHandler struct {
	repo *repository.EnrollmentRepository
}

// NewEnrollmentHandler creates a new enrollment handler
func NewEnrollmentHandler(repo *repository.EnrollmentRepository) *EnrollmentHandler {
	return &EnrollmentHandler{
		repo: repo,
	}
}

// CreateEnrollment handles POST /api/enrollments
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

// GetEnrollment handles GET /api/enrollments/{id}
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

// GetAllEnrollments handles GET /api/enrollments
func (h *EnrollmentHandler) GetAllEnrollments(w http.ResponseWriter, r *http.Request) {
	enrollments := h.repo.GetAll()
	respondWithJSON(w, http.StatusOK, enrollments)
}

// UpdateEnrollment handles PUT /api/enrollments/{id}
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

// DeleteEnrollment handles DELETE /api/enrollments/{id}
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

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON sends a JSON response
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
