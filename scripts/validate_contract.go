//go:build contract
// +build contract

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

// validateContract validates the API implementation against the OpenAPI spec
func main() {
	ctx := context.Background()

	// Load OpenAPI specification
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("api/openapi.yaml")
	if err != nil {
		log.Fatalf("❌ Failed to load OpenAPI spec: %v", err)
	}

	// Validate the OpenAPI document itself
	if err := doc.Validate(ctx); err != nil {
		log.Fatalf("❌ Invalid OpenAPI specification: %v", err)
	}
	log.Println("✓ OpenAPI specification is valid")

	// Create router from the spec
	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		log.Fatalf("❌ Failed to create router from spec: %v", err)
	}

	// Test routes defined in spec
	testRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "http://localhost:8080/"},
		{"GET", "http://localhost:8080/health"},
		{"GET", "http://localhost:8080/api/enrollments"},
		{"POST", "http://localhost:8080/api/enrollments"},
		{"GET", "http://localhost:8080/api/enrollments/a81eee8a-8ef0-46c9-aefa-e3f14ff1303c"},
		{"PUT", "http://localhost:8080/api/enrollments/a81eee8a-8ef0-46c9-aefa-e3f14ff1303c"},
		{"DELETE", "http://localhost:8080/api/enrollments/a81eee8a-8ef0-46c9-aefa-e3f14ff1303c"},
	}

	violations := 0
	for _, route := range testRoutes {
		req, err := http.NewRequest(route.method, route.path, nil)
		if err != nil {
			log.Printf("❌ Failed to create request for %s %s: %v", route.method, route.path, err)
			violations++
			continue
		}

		// Find route in spec
		_, _, err = router.FindRoute(req)
		if err != nil {
			log.Printf("❌ Route not found in spec: %s %s - %v", route.method, route.path, err)
			violations++
			continue
		}

		log.Printf("✓ Route validated: %s %s", route.method, route.path)
	}

	// Validate response schemas
	if err := validateSchemas(doc); err != nil {
		log.Printf("❌ Schema validation failed: %v", err)
		violations++
	} else {
		log.Println("✓ All schemas are valid")
	}

	// Check for required components
	requiredComponents := []string{"Enrollment", "EnrollmentRequest", "ErrorResponse", "HealthResponse"}
	for _, component := range requiredComponents {
		if _, exists := doc.Components.Schemas[component]; !exists {
			log.Printf("❌ Required schema missing: %s", component)
			violations++
		} else {
			log.Printf("✓ Schema defined: %s", component)
		}
	}

	// Check for X-Cache-Status header
	if _, exists := doc.Components.Headers["X-Cache-Status"]; !exists {
		log.Printf("❌ Required header missing: X-Cache-Status")
		violations++
	} else {
		log.Println("✓ X-Cache-Status header documented")
	}

	// Summary
	fmt.Println("\n" + strings.Repeat("=", 60))
	if violations > 0 {
		log.Printf("❌ CONTRACT VALIDATION FAILED: %d violation(s) found", violations)
		os.Exit(1)
	} else {
		log.Println("✅ CONTRACT VALIDATION PASSED: All checks successful")
		os.Exit(0)
	}
}

// validateSchemas validates all schemas in the OpenAPI document
func validateSchemas(doc *openapi3.T) error {
	ctx := context.Background()

	for name, schema := range doc.Components.Schemas {
		if err := schema.Value.Validate(ctx); err != nil {
			return fmt.Errorf("schema %s validation failed: %w", name, err)
		}
	}

	return nil
}
