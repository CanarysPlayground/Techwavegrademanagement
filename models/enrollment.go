// Package models defines the domain entities and business logic for the Grade Management API.
package models

import (
	"errors"
	"time"
)

// Enrollment represents a student enrollment in a course with comprehensive tracking.
// It manages the lifecycle of a student's participation in a course from enrollment to completion.
//
// Business Rules:
//   - Each enrollment must have a unique ID (UUID format)
//   - StudentID and CourseID are required and reference external entities
//   - Status must be one of: "pending", "active", or "completed"
//   - EnrollmentDate defaults to current time if not provided
//   - Timestamps (CreatedAt, UpdatedAt) are automatically managed
//
// Example:
//
//	enrollment := &models.Enrollment{
//	    StudentID: "student-123",
//	    CourseID: "course-456",
//	    Status: "active",
//	}
type Enrollment struct {
	// ID is the unique identifier for the enrollment (auto-generated UUID)
	ID string `json:"id"`

	// StudentID references the student taking the course (required)
	StudentID string `json:"student_id"`

	// CourseID references the course being taken (required)
	CourseID string `json:"course_id"`

	// EnrollmentDate is when the student enrolled (defaults to current time)
	EnrollmentDate time.Time `json:"enrollment_date"`

	// Status indicates the current state: "pending", "active", or "completed" (required)
	Status string `json:"status"`

	// CreatedAt is when the enrollment record was created (auto-set)
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is when the enrollment was last modified (auto-updated)
	UpdatedAt time.Time `json:"updated_at"`
}

// ValidStatuses contains the allowed status values for enrollment lifecycle.
// Business rules:
//   - "pending": Student enrolled but course hasn't started
//   - "active": Student is currently taking the course
//   - "completed": Student has finished the course
var ValidStatuses = map[string]bool{
	"pending":   true,
	"active":    true,
	"completed": true,
}

// Validate checks if the enrollment data meets all business rules and constraints.
// It ensures required fields are present and status values are valid.
//
// Returns:
//   - nil if validation passes
//   - error describing the validation failure
//
// Example:
//
//	enrollment := &Enrollment{StudentID: "", CourseID: "course-1", Status: "active"}
//	err := enrollment.Validate()
//	// Returns: "student_id is required"
func (e *Enrollment) Validate() error {
	if e.StudentID == "" {
		return errors.New("student_id is required")
	}
	if e.CourseID == "" {
		return errors.New("course_id is required")
	}
	if e.Status == "" {
		return errors.New("status is required")
	}
	if !ValidStatuses[e.Status] {
		return errors.New("status must be one of: pending, active, completed")
	}
	return nil
}
