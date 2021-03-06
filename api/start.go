package handler

import (
	"os"
	"io"
    "fmt"
    "time"
    "strings"
    "syscall"
    "runtime"
    "os/exec"
    "net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    currentTime := time.Now().Format(time.RFC850)
    fmt.Fprintf(w, "%v\n", currentTime)
    fmt.Fprintf(w, "%#v\n", runtime.GOOS)
    fmt.Fprintf(w, "%v\n", strings.Join(os.Environ(), "\n"))
    fmt.Fprintf(w, "pagesize=%d\n", syscall.Getpagesize())
    name, err := os.Hostname()
    fmt.Fprintf(w, "%s, %v\n", name, err)
	var cmdString []string = []string{}
    cmd := cmdMake(w, "./handler", cmdString)
	err = cmd.Start()
	if err != nil { fmt.Fprintf(w, "Start failed: %v\n", err) }
	err = cmd.Wait()
	if err != nil { fmt.Fprintf(w, "Wait failed: %v\n", err) }
    cmdString = []string{}
    cmd = cmdMake(w, "pwd", cmdString)
    err = cmd.Start()
    if err != nil { fmt.Fprintf(w, "Start failed: %v\n", err) }
    err = cmd.Wait()
    if err != nil { fmt.Fprintf(w, "Wait failed: %v\n", err) }
    cmdString = []string{"-la", ".", ".."}
    cmd = cmdMake(w, "ls", cmdString)
    err = cmd.Start()
    if err != nil { fmt.Fprintf(w, "Start failed: %v\n", err) }
    err = cmd.Wait()
    if err != nil { fmt.Fprintf(w, "Wait failed: %v\n", err) }
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
