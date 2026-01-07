package repository

import (
	"errors"
	"sync"
	"techwave/models"
)

var (
	// ErrNotFound is returned when an enrollment is not found
	ErrNotFound = errors.New("enrollment not found")
	// ErrAlreadyExists is returned when an enrollment already exists
	ErrAlreadyExists = errors.New("enrollment already exists")
)

// EnrollmentRepository manages enrollment data storage
type EnrollmentRepository struct {
	mu          sync.RWMutex
	enrollments map[string]*models.Enrollment
}

// NewEnrollmentRepository creates a new enrollment repository
func NewEnrollmentRepository() *EnrollmentRepository {
	return &EnrollmentRepository{
		enrollments: make(map[string]*models.Enrollment),
	}
}

// Create adds a new enrollment to the repository
func (r *EnrollmentRepository) Create(enrollment *models.Enrollment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.enrollments[enrollment.ID]; exists {
		return ErrAlreadyExists
	}

	r.enrollments[enrollment.ID] = enrollment
	return nil
}

// GetByID retrieves an enrollment by ID
func (r *EnrollmentRepository) GetByID(id string) (*models.Enrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	enrollment, exists := r.enrollments[id]
	if !exists {
		return nil, ErrNotFound
	}

	return enrollment, nil
}

// GetAll retrieves all enrollments
func (r *EnrollmentRepository) GetAll() []*models.Enrollment {
	r.mu.RLock()
	defer r.mu.RUnlock()

	enrollments := make([]*models.Enrollment, 0, len(r.enrollments))
	for _, enrollment := range r.enrollments {
		enrollments = append(enrollments, enrollment)
	}

	return enrollments
}

// Update modifies an existing enrollment
func (r *EnrollmentRepository) Update(id string, enrollment *models.Enrollment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.enrollments[id]; !exists {
		return ErrNotFound
	}

	enrollment.ID = id
	r.enrollments[id] = enrollment
	return nil
}

// Delete removes an enrollment from the repository
func (r *EnrollmentRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.enrollments[id]; !exists {
		return ErrNotFound
	}

	delete(r.enrollments, id)
	return nil
}
