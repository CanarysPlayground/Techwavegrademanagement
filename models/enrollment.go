package models

import (
	"errors"
	"time"
)

// Enrollment represents a student enrollment in a course
type Enrollment struct {
	ID             string    `json:"id"`
	StudentID      string    `json:"student_id"`
	CourseID       string    `json:"course_id"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ValidStatuses contains the allowed status values
var ValidStatuses = map[string]bool{
	"pending":   true,
	"active":    true,
	"completed": true,
}

// Validate checks if the enrollment data is valid
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
