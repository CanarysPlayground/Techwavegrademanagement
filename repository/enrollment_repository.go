// Package repository provides data storage and retrieval operations for the Grade Management API.
// It implements an in-memory storage solution with thread-safe operations.
package repository

import (
	"errors"
	"sync"
	"techwave/models"
)

var (
	// ErrNotFound is returned when an enrollment is not found in the repository.
	// This typically results in a 404 HTTP response.
	ErrNotFound = errors.New("enrollment not found")

	// ErrAlreadyExists is returned when attempting to create an enrollment with a duplicate ID.
	// This typically results in a 409 HTTP response.
	ErrAlreadyExists = errors.New("enrollment already exists")
)

// EnrollmentRepository manages enrollment data storage with thread-safe operations.
// It uses an in-memory map for storage with read-write mutex for concurrent access.
//
// Note: This is an in-memory implementation. Data will be lost on server restart.
// For production use, replace with persistent storage (database, Redis, etc.).
type EnrollmentRepository struct {
	mu          sync.RWMutex                 // Protects concurrent access to enrollments map
	enrollments map[string]*models.Enrollment // In-memory storage indexed by enrollment ID
}

// NewEnrollmentRepository creates a new enrollment repository instance.
// Initializes the internal map for storing enrollments.
//
// Returns:
//   - A configured EnrollmentRepository ready for use
//
// Example:
//
//	repo := repository.NewEnrollmentRepository()
//	enrollment := &models.Enrollment{ID: "123", StudentID: "s1", CourseID: "c1", Status: "active"}
//	err := repo.Create(enrollment)
func NewEnrollmentRepository() *EnrollmentRepository {
	return &EnrollmentRepository{
		enrollments: make(map[string]*models.Enrollment),
	}
}

// Create adds a new enrollment to the repository.
// Uses write lock to ensure thread-safe insertion.
//
// Parameters:
//   - enrollment: The enrollment to create (must have unique ID)
//
// Returns:
//   - nil on success
//   - ErrAlreadyExists if an enrollment with the same ID exists
//
// Example:
//
//	enrollment := &models.Enrollment{
//	    ID: uuid.New().String(),
//	    StudentID: "student-123",
//	    CourseID: "course-456",
//	    Status: "active",
//	}
//	err := repo.Create(enrollment)
//	if err == repository.ErrAlreadyExists {
//	    // Handle duplicate
//	}
func (r *EnrollmentRepository) Create(enrollment *models.Enrollment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.enrollments[enrollment.ID]; exists {
		return ErrAlreadyExists
	}

	r.enrollments[enrollment.ID] = enrollment
	return nil
}

// GetByID retrieves an enrollment by its unique identifier.
// Uses read lock to allow concurrent read operations.
//
// Parameters:
//   - id: The enrollment UUID to retrieve
//
// Returns:
//   - The enrollment if found
//   - ErrNotFound if no enrollment with the given ID exists
//
// Example:
//
//	enrollment, err := repo.GetByID("123e4567-e89b-12d3-a456-426614174000")
//	if err == repository.ErrNotFound {
//	    // Handle not found
//	}
func (r *EnrollmentRepository) GetByID(id string) (*models.Enrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	enrollment, exists := r.enrollments[id]
	if !exists {
		return nil, ErrNotFound
	}

	return enrollment, nil
}

// GetAll retrieves all enrollments from the repository.
// Uses read lock to ensure thread-safe access. Returns a new slice to prevent
// external modifications to the internal storage.
//
// Returns:
//   - A slice of all enrollments (may be empty if no enrollments exist)
//
// Example:
//
//	enrollments := repo.GetAll()
//	for _, enrollment := range enrollments {
//	    fmt.Printf("Enrollment: %s - Student: %s\n", enrollment.ID, enrollment.StudentID)
//	}
func (r *EnrollmentRepository) GetAll() []*models.Enrollment {
	r.mu.RLock()
	defer r.mu.RUnlock()

	enrollments := make([]*models.Enrollment, 0, len(r.enrollments))
	for _, enrollment := range r.enrollments {
		enrollments = append(enrollments, enrollment)
	}

	return enrollments
}

// Update modifies an existing enrollment in the repository.
// Uses write lock to ensure thread-safe modification. Creates a copy of the
// enrollment to prevent external modifications.
//
// Parameters:
//   - id: The enrollment UUID to update
//   - enrollment: The updated enrollment data (ID will be overwritten with the id parameter)
//
// Returns:
//   - nil on success
//   - ErrNotFound if no enrollment with the given ID exists
//
// Example:
//
//	updated := &models.Enrollment{
//	    StudentID: "student-123",
//	    CourseID: "course-456",
//	    Status: "completed",
//	}
//	err := repo.Update("existing-id", updated)
//	if err == repository.ErrNotFound {
//	    // Handle not found
//	}
func (r *EnrollmentRepository) Update(id string, enrollment *models.Enrollment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.enrollments[id]; !exists {
		return ErrNotFound
	}

	// Create a copy to avoid modifying the input
	updated := *enrollment
	updated.ID = id
	r.enrollments[id] = &updated
	return nil
}

// Delete removes an enrollment from the repository.
// Uses write lock to ensure thread-safe deletion. This is a hard delete operation.
//
// Parameters:
//   - id: The enrollment UUID to delete
//
// Returns:
//   - nil on success
//   - ErrNotFound if no enrollment with the given ID exists
//
// Example:
//
//	err := repo.Delete("123e4567-e89b-12d3-a456-426614174000")
//	if err == repository.ErrNotFound {
//	    // Handle not found
//	}
func (r *EnrollmentRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.enrollments[id]; !exists {
		return ErrNotFound
	}

	delete(r.enrollments, id)
	return nil
}
