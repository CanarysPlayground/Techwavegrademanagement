package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"techwave/cache"
	"techwave/handlers"
	"techwave/repository"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Initialize Redis client with connection pooling
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // Default to local Redis
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password:     os.Getenv("REDIS_PASSWORD"), // No password by default
		DB:           0,                           // Use default DB
		PoolSize:     10,                          // Connection pool size
		MinIdleConns: 5,                           // Minimum idle connections
	})

	// Test Redis connection (graceful fallback if unavailable)
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("WARNING: Redis unavailable, running without cache: %v", err)
		redisClient = nil // Disable caching
	} else {
		log.Println("✓ Redis connection established")
	}

	// Initialize repository
	enrollmentRepo := repository.NewEnrollmentRepository()

	// Initialize cache (nil-safe, graceful degradation)
	var enrollmentCache *cache.EnrollmentCache
	if redisClient != nil {
		enrollmentCache = cache.NewEnrollmentCache(redisClient)
		log.Println("✓ Cache layer enabled (5-minute TTL)")
	}

	// Initialize handlers with cache
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentRepo, enrollmentCache)

	// Setup router
	router := mux.NewRouter()

	// Root endpoint
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cacheStatus := "disabled"
		if redisClient != nil {
			cacheStatus = "enabled"
		}
		fmt.Fprintf(w, "Grade Management API - Cache: %s", cacheStatus)
	}).Methods("GET")

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		health := map[string]interface{}{
			"status": "healthy",
			"cache":  redisClient != nil,
		}
		if redisClient != nil {
			health["redis"] = redisClient.Ping(ctx).Err() == nil
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"healthy","cache":%v}`, health["cache"])
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
	fmt.Printf("🚀 Starting Grade Management API on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
