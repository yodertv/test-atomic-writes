package start

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
    fmt.Fprintf(w, "pagesize=%d\n", syscall.Getpagesize())
    name, err := os.Hostname()
    fmt.Fprintf(w, "%s, %v\n", name, err)
    name, err = os.Getwd() 
    fmt.Fprintf(w, "%s, %v\n", name, err)

/*
    // This idea doesn't work in vercel. I can't put executables into the api directory and the api 
    // directory can't see any of the public files served by the vercel run-time.
	var cmdName string = "./test-atomic-writes"
    var cmdArgs []string = []string{}
    var cmd *exec.Cmd = cmdMake(w, cmdName, cmdArgs)
	err = cmd.Start()
    if err != nil { fmt.Fprintf(w, "Start failed for command %s %s: %v\n", cmdName, cmdArgs, err) }
	err = cmd.Wait()
    if err != nil { fmt.Fprintf(w, "Wait failed for command %s %s: %v\n",cmdName, cmdArgs, err) }
 */
    cmdName := "pwd"
    cmdArgs := []string{}
    cmd := cmdMake(w, cmdName, cmdArgs)
    err = cmd.Start()
    if err != nil { fmt.Fprintf(w, "Start failed for command %s %s: %v\n", cmdName, cmdArgs, err) }
    err = cmd.Wait()
    if err != nil { fmt.Fprintf(w, "Wait failed for command %s %s: %v\n",cmdName, cmdArgs, err) }
    cmdName = "ls"
    cmdArgs = []string{"-la", "/", ".", ".."}
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
