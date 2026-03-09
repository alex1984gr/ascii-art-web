// Package main contains the HTTP server and request handlers for the ASCII art web application
package main

// Import required packages for HTTP server, templating, and error handling
import (
	"errors"        // Functions for error comparison and checking
	"html/template" // HTML template parsing and execution
	"net/http"      // HTTP server and request/response handling
	"os"            // Operating system functions for file error checking

	"ascii-art/pipeline" // Custom package containing ASCII art rendering logic
)

// templatePath defines the file path to the HTML template for the web interface
const templatePath = "templates/index.html"

// pageData is a struct that holds data to be passed to the HTML template
type pageData struct {
	Input  string // The text input submitted by the user
	Banner string // The selected banner style (standard, shadow, thinkertoy)
	Result string // The generated ASCII art result to display
}

// runServer initializes the HTTP routes and starts the web server on port 8080
func runServer() error {
	// Register homeHandler to handle requests to the root path "/"
	http.HandleFunc("/", homeHandler)
	// Register asciiArtHandler to handle POST requests to "/ascii-art"
	http.HandleFunc("/ascii-art", asciiArtHandler)
	// Start the HTTP server listening on port 8080 and return any error
	return http.ListenAndServe(":8080", nil)
}

// homeHandler handles GET requests to the root path and displays the main form
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the requested path is exactly "/" (reject other paths)
	if r.URL.Path != "/" {
		// Return 404 Not Found for any path other than root
		http.NotFound(w, r)
		return
	}
	// Check if the HTTP method is GET (reject POST, PUT, etc.)
	if r.Method != http.MethodGet {
		// Return 400 Bad Request for non-GET methods
		http.Error(w, "400 bad request", http.StatusBadRequest)
		return
	}

	// Render the page with default banner "standard" and empty input/result
	renderPage(w, pageData{Banner: "standard"}, http.StatusOK)
}

// asciiArtHandler handles POST requests to generate ASCII art from form data
func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the requested path is exactly "/ascii-art"
	
	// Parse the form data from the request body
	if err := r.ParseForm(); err != nil {
		// Return 400 Bad Request if form parsing fails
		http.Error(w, "400 bad request", http.StatusBadRequest)
		return
	}

	// Extract the "text" field from the form (user's input text)
	input := r.FormValue("text")
	// Extract the "banner" field from the form (selected banner style)
	banner := r.FormValue("banner")
	// If no banner was selected, default to "standard"
	if banner == "" {
		banner = "standard"
	}

	// Call the pipeline to render ASCII art with the provided input and banner
	result, err := pipeline.RenderASCII(input, banner)
	// Check if an error occurred during rendering
	if err != nil {
		// Use a switch to determine the appropriate HTTP status code based on error type
		switch {
		// If error is invalid input or invalid banner, return 400 Bad Request
		case errors.Is(err, pipeline.ErrInvalidInput), errors.Is(err, pipeline.ErrInvalidBanner):
			http.Error(w, "400 bad request", http.StatusBadRequest)
		// If error is file not found (banner file missing), return 404 Not Found
		case errors.Is(err, os.ErrNotExist):
			http.Error(w, "404 not found", http.StatusNotFound)
		// For any other unexpected error, return 500 Internal Server Error
		default:
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Render the page with the input, banner, and generated ASCII art result
	renderPage(w, pageData{Input: input, Banner: banner, Result: result}, http.StatusOK)
}

// renderPage loads the HTML template, populates it with data, and writes it to the response
func renderPage(w http.ResponseWriter, data pageData, status int) {
	// Parse the HTML template file from disk
	tmpl, err := template.ParseFiles(templatePath)
	// Check if template parsing failed
	if err != nil {
		// If the template file doesn't exist, return 404 Not Found
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "404 not found", http.StatusNotFound)
			return
		}
		// For any other template error, return 500 Internal Server Error
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	// Set the HTTP status code in the response header
	w.WriteHeader(status)
	// Execute the template with the provided data and write output to response writer
	if err := tmpl.Execute(w, data); err != nil {
		// If template execution fails, return 500 Internal Server Error
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
}
