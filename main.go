package main

import (
	"fmt"
	"os"
	"log"
	"net/http"

	"github.com/yodertv/test-atomic-writes/lib"
)

func main() {
	// Define the port to listen on
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// 1. Setup a multiplexer (router)
	mux := http.NewServeMux()

	// 2. Wire a handler function to it's api path.
	mux.HandleFunc("/api/start", lib.StartHandler)

	// 3. Example API endpoint (if your service builds/handles other requests)
	mux.HandleFunc("/api/restart", lib.RestartHandler)

	// 5. Any route not caught by specific handlers will hit the "/" handler
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
	    w.Header().Set("Content-Type", "application/json")
	    w.WriteHeader(http.StatusNotFound)
	    fmt.Fprintln(w, `{"error": "This API endpoint does not exist"}`)
	})

	// Start the server
	log.Printf("test-atomic-writes-server starting on port %s...", port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
