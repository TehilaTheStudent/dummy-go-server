package main

import (
	"log"
	"net/http"
	"os"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)
func main() {
	r := chi.NewRouter()
	// Use both chi's logger and our custom logger for demonstration
	r.Use(middleware.Logger)
	r.Use(requestLogger)
	r.Use(corsMiddleware)

	r.Get("/health", healthHandler)

	r.Route("/users", func(r chi.Router) {
		r.Get("/", listUsers)
		r.Post("/", createUser)
		r.Get("/{id}", getUser)
		r.Put("/{id}", updateUser)
		r.Delete("/{id}", deleteUser)
	})
	
	// OpenAPI schema routes
	r.Get("/schemas", listAvailableSchemas) // List all available schemas
	r.Get("/schemas/*", serveOpenAPISchema) // Serve specific schema files
	
	// Conversion routes
	r.Get("/convert/yaml-to-json/*", convertYAMLToJSON)
	r.Get("/convert/json-to-yaml/*", convertJSONToYAML)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on 0.0.0.0:%s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
