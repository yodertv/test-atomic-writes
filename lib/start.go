package lib

import (
	"os"
	"io"
    "fmt"
    "flag"
    "time"
    "errors"
    "syscall"
    "runtime"
    "os/exec"
    "net/http"
)

var (
    Args = Cmdline_args{}
    Empty_args = Cmdline_args{}
)

func StartHandler(w http.ResponseWriter, r *http.Request) {
    currentTime := time.Now().Format(time.RFC850)
    fmt.Fprintf(w, "Start time: %v\n", currentTime)
    fmt.Fprintf(w, "Comandline Args: %v\n", os.Args)
    fmt.Fprintf(w, "PID: %v\n", os.Getpid())
    fmt.Fprintf(w, "GOOS: %#v\n", runtime.GOOS)
    fmt.Fprintf(w, "Pagesize: %d\n", syscall.Getpagesize())
    name, err := os.Hostname()
    fmt.Fprintf(w, "Hostname: %s, err: %v\n", name, err)
    name, err = os.Getwd()
    fmt.Fprintf(w, "Getwd: %s, err: %v\n", name, err)
    if Args == Empty_args {
        Parse_args(&Args)
        if len(flag.Args()) != 0 { // Extraneous arguments
            flag.Usage();
            fmt.Fprintf(w, "Parse_args: See server log for Usage info.\n")
            return
        }
    }
    fmt.Fprintf(w, "Params: %v", Args)

    // The idea to fork an executable doesn't work in vercel lamda. I can't put executables into the api directory and the api
    // directory can't see any of the public files served by the vercel run-time. Which makes sense from a security perspective.

    var cmdName = "touch"
    var cmdArgs []string = []string{Args.Filename}
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

    Write_bytes(Args.Count, Args.Size, Args.Workers, -1, Args.Filename)
    res := Validate_bytes(Args.Filename, Args.Count, Args.Size, Args.Workers)
    fmt.Fprintf(w, "Validate_bytes returned: %d\n", res)

    cmdName = "ls"
    cmdArgs = []string{"-la", Args.Filename}
    cmd = cmdMake(w, cmdName, cmdArgs)
    err = cmd.Start()
    if err != nil { fmt.Fprintf(w, "Start failed for command %s %s: %v\n", cmdName, cmdArgs, err) }
    err = cmd.Wait()
    if err != nil { fmt.Fprintf(w, "Wait failed for command %s %s: %v\n",cmdName, cmdArgs, err) }
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
