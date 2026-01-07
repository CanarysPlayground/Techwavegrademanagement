package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Grade Management API - Ready for AI delegation!")
	})

	port := ":8080"
	fmt.Printf("Starting Grade Management API on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}