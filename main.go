// Package main is the entry point for the ASCII art web server application
package main

// Import the log package for error logging and fatal error handling
import "log"

// main is the entry point function that starts when the program runs
func main() {
	// Call runServer() to start the HTTP server and capture any error it returns
	if err := runServer(); err != nil {
		// If runServer returns an error, log it and terminate the program with exit code 1
		log.Fatal(err)
	}
}
