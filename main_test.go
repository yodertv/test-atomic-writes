// main_test.go
//
// Test the pscat utility
//
package main_test

import (
	"os"
	"os/exec"
	"testing"
)

// tCheck is a helper function for fatal tests
func tCheck(e error, t *testing.T) {
	t.Helper()
    if e != nil {
        t.Fatal(e)
    }
}

// cmdMake returns an exec.Cmd with the command s set up to be started or run with the args.
// It copies the io of the parent. It marks the environment with a "TESTING" flag.
// Should the command care to check for the flag, it could do so like:
//    if os.Getenv("CMD_IN_PROGRESS") == "1" {
//       fmt.Println("CMD in progress.")
//    }
func cmdMake(s string, args []string) *exec.Cmd {
    env := []string{
		"CMD_IN_PROGRESS=1",
    }
    cmd := exec.Command(s, args...)
    cmd.Env = append(os.Environ(), env...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd
}

// TestAtomicWrites runs the cmdline executable with default arguments
func TestAtomicWrites(t *testing.T){
	cmdString := "./test-atomic-writes"
	cmdArgs := []string{}
    cmd := cmdMake(cmdString, cmdArgs)
	err := cmd.Start()
	if err != nil { t.Errorf("%s cmd start failed: %v\n", cmdString, err) }
	err = cmd.Wait()
	if err != nil { t.Errorf("Wait failed: %v\n", err) }
}
