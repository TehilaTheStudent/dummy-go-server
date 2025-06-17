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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on :%s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
