package api

import (
	"os"
	"io"
    "fmt"
    "time"
    "errors"
    "syscall"
    "runtime"
    "os/exec"
    "net/http"
    "path/filepath"

    "github.com/yodertv/test-atomic-writes/lib"
)

var args = lib.Cmdline_args{ 50, 4096, 3, -1, false, "../testdata" }

func StartHandler(w http.ResponseWriter, r *http.Request) {
    currentTime := time.Now().Format(time.RFC850)
    fmt.Fprintf(w, "%v\n", currentTime)
    fmt.Fprintf(w, "%#v\n", runtime.GOOS)
    fmt.Fprintf(w, "pagesize=%d\n", syscall.Getpagesize())
    name, err := os.Hostname()
    fmt.Fprintf(w, "Hostname=%s, err=%v\n", name, err)
    name, err = os.Getwd() 
    fmt.Fprintf(w, "Working directory=%s, err=%v\n", name, err)

    // The idea to fork an executable doesn't work in vercel lamda. I can't put executables into the api directory and the api
    // directory can't see any of the public files served by the vercel run-time. Which makes sense from a security perspective.

    var cmdName = "touch"
    var filePath = filepath.Join("/tmp", "hello")
    var cmdArgs []string = []string{filePath}
//    var cmd *exec.Cmd = cmdMake(w, cmdName, cmdArgs)
    var cmd *exec.Cmd = exec.Command(cmdName, cmdArgs...)
    // CombinedOutput runs the command and returns both stdout and stderr together
    output, cerr := cmd.CombinedOutput()
    if cerr != nil {
        fmt.Fprintf(w, "Start failed for command %s %s: %v\n", cmdName, cmdArgs, cerr)
        // 1. Extract the exit code
        var exitErr *exec.ExitError
        if errors.As(cerr, &exitErr) {
            fmt.Fprintf(w, "Command failed with Exit Code: %d\n", exitErr.ExitCode())
        } else {
            fmt.Fprintf(w, "Failed to run command %s %s (system error): %v\n", cmdName, cmdArgs, cerr)
        }
        // 2. Extract the stderr output
        fmt.Fprintf(w, "Error output (stderr): %s %v\n", string(output), cerr)
    } else {
        fmt.Fprintf(w, "%s\nSuccess error: %v\n", string(output), cerr)
    }

    cmdName = "pwd"
    cmdArgs = []string{}
    cmd = cmdMake(w, cmdName, cmdArgs)
	err = cmd.Start()
    if err != nil { fmt.Fprintf(w, "Start failed for command %s %s: %v\n", cmdName, cmdArgs, err) }
	err = cmd.Wait()
    if err != nil { fmt.Fprintf(w, "Wait failed for command %s %s: %v\n",cmdName, cmdArgs, err) }

    cmdName = "ls"
    cmdArgs = []string{"-laR", "/tmp/"}
    cmd = cmdMake(w, cmdName, cmdArgs)
    err = cmd.Start()
    if err != nil { fmt.Fprintf(w, "Start failed for command %s %s: %v\n", cmdName, cmdArgs, err) }
    err = cmd.Wait()
    if err != nil { fmt.Fprintf(w, "Wait failed for command %s %s: %v\n",cmdName, cmdArgs, err) }

    cmdName = "cat"
    cmdArgs = []string{filePath}
    cmd = cmdMake(w, cmdName, cmdArgs)
    err = cmd.Start()
    if err != nil { fmt.Fprintf(w, "Start failed for command %s %s: %v\n", cmdName, cmdArgs, err) }
    err = cmd.Wait()
    if err != nil { fmt.Fprintf(w, "Wait failed for command %s %s: %v\n",cmdName, cmdArgs, err) }

    res := lib.Validate_bytes(args.Filename, args.Count, args.Size, args.Workers)
    fmt.Fprintf(w, "Validate_bytes returned: %d\n", res)
}

// cmdMake returns an exec.Cmd with the command s set up to be started or run with args.
// It copies the io of the parent except for stdout which is replaced by the io.Writer w.
// It marks the environment with a "TESTING" flag.
// Should the pscat command care to check for the flag, it could do so like:
//    if os.Getenv("CMD_IN_PROGRESS") == "1" {
//       fmt.Println("Command in progress.")
//    }
func cmdMake(w io.Writer, s string, args []string) *exec.Cmd {
    env := []string{
		"CMD_IN_PROGRESS=1",
    }
    cmd := exec.Command(s, args...)
    cmd.Env = append(os.Environ(), env...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = w
    cmd.Stderr = os.Stderr
    return cmd
}
