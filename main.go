package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Wrap the real writer with our custom one, pre-set 200
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		// Log method, path, status, duration
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, rw.statusCode, time.Since(start))
	})
}

// Application Struct
type application struct {
	logger *slog.Logger
}

// Handlers
func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "status: available\n")
	app.logger.Info("healthcheck handler called")
}

func (app *application) listBooks(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "list of books (coming soon)\n")
	app.logger.Info("listBooks handler called")
}

func (app *application) getBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "get book with id: %s\n", id)
	app.logger.Info("getBook handler called", "id", id)
}

func (app *application) createBook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "book created (coming soon)\n")
	app.logger.Info("createBook handler called")
}

func (app *application) deleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.WriteHeader(http.StatusNoContent) // 204: No Content, body not allowed
	app.logger.Info("deleteBook handler called", "id", id)
}

// Main Function
func main() {
	// Create structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{logger: logger}

	// Register routes with method-qualified patterns
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/healthcheck", app.healthcheck)
	mux.HandleFunc("GET /v1/books", app.listBooks)
	mux.HandleFunc("GET /v1/books/{id}", app.getBook)
	mux.HandleFunc("POST /v1/books", app.createBook)
	mux.HandleFunc("DELETE /v1/books/{id}", app.deleteBook)

	// Startup log
	logger.Info("starting server", "addr", ":4000")

	// Wrap entire router with logging middleware
	err := http.ListenAndServe(":4000", loggingMiddleware(mux))
	log.Fatal(err)
}
