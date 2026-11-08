package main

import (
	"fmt"
	"os"
	"log"
	"io/fs"
	"embed"
	"net/http"

	"github.com/yodertv/test-atomic-writes/lib"
)

//go:embed public/*
var publicFiles embed.FS

func main() {
	// Define the port to listen on
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Setup a multiplexer (router)
	mux := http.NewServeMux()

	// Wire the StartHandler function to it's api path.
	mux.HandleFunc("/api/start", lib.StartHandler)

	// Restart API endpoint (if your service builds/handles other requests)
	mux.HandleFunc("/api/restart", lib.RestartHandler)

	// Any api requests not caught by specific handlers will hit the api home, "/api/", handler
	mux.HandleFunc("/api/", handleHome)

	// Create a sub-file system rooted at "public"
	subFS, err := fs.Sub(publicFiles, "public")
	if err != nil {
		log.Fatalf("Failed to create static fileServer: %v", err)
		return
	}

	// Serve static files using the embedded filesystem
	fileServer := http.FileServer(http.FS(subFS))
	// Static File Server handler that serves content out of the localized "./static" folder
	// fileServer := http.FileServer(http.Dir("./public"))

	// Registers the file server. By appending a trailing slash, it catches all routes
	// that didn't match the specific API endpoints above.
	mux.Handle("/", fileServer)

	// Start the server
	log.Printf("test-atomic-writes-server starting on port %s...", port)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/" {
	    w.Header().Set("Content-Type", "application/json")
	    w.WriteHeader(http.StatusNotFound)
	    fmt.Fprintln(w, `{"error": "This API endpoint does not exist"}`)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "Welcome to test-atomic-writes online!")
}
