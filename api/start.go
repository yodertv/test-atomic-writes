package handler

import (
	"os"
	"io"
    "fmt"
    "time"
    "syscall"
    "runtime"
    "os/exec"
    "net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    currentTime := time.Now().Format(time.RFC850)
    fmt.Fprintf(w, "%v\n", currentTime)
    fmt.Fprintf(w, "%#v\n", runtime.GOOS)
    fmt.Fprintf(w, "%#v\n", os.Environ())
    fmt.Fprintf(w, "pagesize=%d\n", syscall.Getpagesize())
	var cmdString []string
    cmd := cmdMake(w, "./test-atomic-writes", cmdString)
	err := cmd.Start()
	if err != nil { fmt.Fprintf(w, "Start failed: %v\n", err) }
	err = cmd.Wait()
	if err != nil { fmt.Fprintf(w, "Wait failed: %v\n", err) }
}

// subPscat returns an exec.Cmd with the pscat subtest set up to be started or run with the subcommand arguments.
// It copies the io of the parent. It marks the environment with a "TESTING" flag.
// Should the pscat command care to check for the flag, it could do so like:
//    if os.Getenv("PSCAT_TESTING_IN_PROGRESS") == "1" {
//       fmt.Println("Testing in progress!")
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
