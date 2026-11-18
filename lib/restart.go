package lib

import (
	"os"
	"fmt"
	"flag"
	"time"
	"net/http"
)

// Handler complies with the required Vercel Go runtime signature
func RestartHandler(w http.ResponseWriter, r *http.Request) {
	// Restrict the request method if desired
	if r.Method != http.MethodGet {
		msg := fmt.Sprintf("Method Not Allowed: %v\n", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	// Parse args
	if Args == Empty_args {
        Parse_args(&Args)
        if len(flag.Args()) != 0 { // Extraneous arguments
            flag.Usage();
            fmt.Fprintf(w, "Parse_args: See server log for Usage info.\n")
            return
        }
    }

    // Common response data
    fmt.Fprintf(w, "Params: %v\n", Args)
    currentTime := time.Now().Format(time.RFC850)
    fmt.Fprintf(w, "Restart time: %v\n", currentTime)
    fmt.Fprintf(w, "PID: %v\n", os.Getpid())

    // Remove the test file named in Args
	cmdName := "rm"
    cmdArgs := []string{"-f", Args.Filename}
    cmd := cmdMake(w, cmdName, cmdArgs)
    err1 := cmd.Start()
    if err1 != nil { fmt.Fprintf(w, "Start failed for command %s %s: %v\n", cmdName, cmdArgs, err1) }
    err1 = cmd.Wait()
    if err1 != nil { fmt.Fprintf(w, "Wait failed for command %s %s: %v\n",cmdName, cmdArgs, err1) }

	// Retrieve the hook URL from environment variables
	hookURL := os.Getenv("VERCEL_DEPLOY_HOOK_URL")
	if hookURL == "" {
		w.Write([]byte(`{"status": "success", "message": "Deploy hook URL environment variable is not set. Skipping deployment trigger."}`))
		return
	}

	// Trigger the deployment via an HTTP POST request
	resp, err := http.Post(hookURL, "application/json", nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to trigger deployment: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Return the status matching Vercel's response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		http.Error(w, fmt.Sprintf("Vercel returned an error status: %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success", "message": "Deployment triggered successfully"}`))
}
