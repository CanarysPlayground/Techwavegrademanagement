// Package main provides the entry point for the Grade Management API server.
// This RESTful API manages student enrollments in courses with full CRUD operations.
package main

import (
	"fmt"
	"log"
	"net/http"
	"techwave/handlers"
	"techwave/repository"

	"github.com/gorilla/mux"
)

// main initializes and starts the Grade Management API server.
// It sets up the enrollment repository, handlers, and HTTP routes.
// The server listens on port 8080 by default.
//
// Example usage:
//
//	go run main.go
//	# Server starts on http://localhost:8080
func main() {
	// Initialize repository
	enrollmentRepo := repository.NewEnrollmentRepository()

	// Initialize handlers
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentRepo)

	// Setup router
	router := mux.NewRouter()

	// Root endpoint
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Grade Management API - Ready for AI delegation!")
	}).Methods("GET")

	// API routes with /api prefix
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Enrollment routes
	apiRouter.HandleFunc("/enrollments", enrollmentHandler.CreateEnrollment).Methods("POST")
	apiRouter.HandleFunc("/enrollments", enrollmentHandler.GetAllEnrollments).Methods("GET")
	apiRouter.HandleFunc("/enrollments/{id}", enrollmentHandler.GetEnrollment).Methods("GET")
	apiRouter.HandleFunc("/enrollments/{id}", enrollmentHandler.UpdateEnrollment).Methods("PUT")
	apiRouter.HandleFunc("/enrollments/{id}", enrollmentHandler.DeleteEnrollment).Methods("DELETE")

	port := ":8080"
	fmt.Printf("Starting Grade Management API on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
