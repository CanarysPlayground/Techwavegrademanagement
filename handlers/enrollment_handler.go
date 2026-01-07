package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"techwave/cache"
	"techwave/models"
	"techwave/repository"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// EnrollmentHandler handles HTTP requests for enrollments
type EnrollmentHandler struct {
	repo  *repository.EnrollmentRepository
	cache *cache.EnrollmentCache
}

// NewEnrollmentHandler creates a new enrollment handler
func NewEnrollmentHandler(repo *repository.EnrollmentRepository, cache *cache.EnrollmentCache) *EnrollmentHandler {
	return &EnrollmentHandler{
		repo:  repo,
		cache: cache,
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
// Implements cache-aside pattern with Redis caching
func (h *EnrollmentHandler) GetEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Try to get from cache first
	if h.cache != nil {
		cachedEnrollment, err := h.cache.Get(id)
		if err == nil && cachedEnrollment != nil {
			// Cache HIT
			w.Header().Set("X-Cache-Status", "HIT")
			respondWithJSON(w, http.StatusOK, cachedEnrollment)
			return
		}
		// Cache MISS - continue to database
		log.Printf("Cache MISS for enrollment ID: %s", id)
	}

	// Get from database
	enrollment, err := h.repo.GetByID(id)
	if err != nil {
		if err == repository.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Enrollment not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve enrollment")
		return
	}

	// Store in cache for next time (cache-aside pattern)
	if h.cache != nil {
		if err := h.cache.Set(enrollment); err != nil {
			log.Printf("Failed to cache enrollment: %v", err)
			// Don't fail the request if caching fails
		}
	}

	// Set cache status to MISS
	w.Header().Set("X-Cache-Status", "MISS")
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

	// Invalidate cache after update
	if h.cache != nil {
		if err := h.cache.Delete(id); err != nil {
			log.Printf("Failed to invalidate cache for enrollment %s: %v", id, err)
		}
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

	// Invalidate cache after delete
	if h.cache != nil {
		if err := h.cache.Delete(id); err != nil {
			log.Printf("Failed to invalidate cache for enrollment %s: %v", id, err)
		}
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
