package main

import (
	"fmt"
	"os"
	"log"
	"flag"
	"io/fs"
	"embed"
	"net/http"

	"github.com/yodertv/test-atomic-writes/lib"
)

//go:embed public/*
var publicFiles embed.FS
var cl = &lib.Args

// Act as an http server when no args are present. Otherwise, act as test-atomic-writes cmdline program.
func main() {
	var res = 0
	if len(flag.Args()) == 0 {
	    handleListenAndServe()
	} else {
		lib.Parse_args(cl)
		if len(flag.Args()) != 0 { // Extraneous arguments
			flag.Usage();
			os.Exit(1)
		}
		log.Printf("test-atomic-writes-server starting with args: %v...", cl)
		if (!cl.Readonly){
			// Note that this function forks a process for each worker.
			lib.Write_bytes(cl.Count, cl.Size, cl.Workers, cl.Worker, cl.Filename)
		}
	    if cl.Worker == -1 { // The orchestrating process has worker index = -1. Only need to validate the file once after all the workers finish.
			res = lib.Validate_bytes(cl.Filename, cl.Count, cl.Size, cl.Workers)
		}
		os.Exit(res)
	}
}

func handleListenAndServe() {
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

	// Registers the file server. By appending a trailing slash, it catches all routes
	// that didn't match the specific API endpoints above.
	mux.Handle("/", fileServer)

	// Start the server
	log.Printf("test-atomic-writes-server listening on port %s...", port)
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
