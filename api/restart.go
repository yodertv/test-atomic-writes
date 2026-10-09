package api

import (
	"fmt"
	"net/http"
	"os"
)

// Handler complies with the required Vercel Go runtime signature
func RestartHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Restrict the request method if desired
	if r.Method != http.MethodGet {
		msg := fmt.Sprintf("Method Not Allowed: %v\n", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	// 2. Retrieve the hook URL from environment variables
	hookURL := os.Getenv("VERCEL_DEPLOY_HOOK_URL")
	if hookURL == "" {
		http.Error(w, "Deploy hook URL environment variable is not set", http.StatusInternalServerError)
		return
	}

	// 3. Trigger the deployment via an HTTP POST request
	resp, err := http.Post(hookURL, "application/json", nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to trigger deployment: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 4. Return the status matching Vercel's response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		http.Error(w, fmt.Sprintf("Vercel returned an error status: %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success", "message": "Deployment triggered successfully"}`))
}
