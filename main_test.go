// main_test.go
//
// Test the test-atomic-writes utility.
//
package main_test

import (
	"os"
	"os/exec"
	"fmt"
	"testing"

	"github.com/yodertv/test-atomic-writes/lib"
)

var (
    res = 0
    cl = lib.Cmdline_args{}
)

// tCheck is a helper function for fatal tests
func tCheck(e error, t *testing.T) {
	t.Helper()
    if e != nil {
        t.Fatal(e)
    }
}

// This test requires a main to parse args and perform worker tasks.
func TestMain(m *testing.M) {
    fmt.Printf("Entering TestMain...\n")
    // Initialize test
    lib.Parse_args(&cl)
    // Run all tests and capture the result, res, only in the parent executable.
    if cl.Worker == -1 { // The orchestrating process has worker index = -1.
        res = m.Run()
    }
    fmt.Printf("Exiting TestMain...\n")
    return
}

func TestExecModWrite(t *testing.T) {
    if !cl.Readonly {
        // Note that this function forks a process for each worker.
        lib.Write_bytes(cl.Count, cl.Size, cl.Workers, cl.Worker, cl.Filename)
    }
    if cl.Worker == -1 { // The orchestrating process has worker index = -1. Only need to validate the file once after all the workers finish.
        res = lib.Validate_bytes(cl.Filename, cl.Count, cl.Size, cl.Workers)
    }
    if (res != 0) {
        t.Errorf("Validate_bytes returned non zero status: %d.", res)
    }
}

// Just validate an existing file. Fails if arguments are different from those used when the file was produced.
func TestExecModRead(t *testing.T) {
    res = lib.Validate_bytes(cl.Filename, cl.Count, cl.Size, cl.Workers)
    if (res != 0) {
        t.Errorf("Validate_bytes returned non zero status")
    }
}

// cmdMake returns an exec.Cmd with the command s set up to be started or run with the args.
// It copies the io of the parent. It marks the environment with a "TESTING" flag.
// Should the command care to check for the flag, it could do so like:
//    if os.Getenv("CMD_IN_PROGRESS") == "1" {
//       fmt.Println("CMD in progress.")
//    }
// Needed to run tests of the utility command produced, test-atomic-writes.
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

// TestAtomicWritesExec runs the cmdline executable with default arguments
func TestExecWrite(t *testing.T){
	cmdString := "./test-atomic-writes"
	cmdArgs := []string{}
    cmd := cmdMake(cmdString, cmdArgs)
	err := cmd.Start()
	if err != nil { t.Errorf("%s cmd start failed: %v\n", cmdString, err) }
	err = cmd.Wait()
	if err != nil { t.Errorf("Wait failed: %v\n", err) }
}

// TestReadOnlyExec runs the cmdline executable with default arguments readonly.
// Todo: This test panics when run before TestAtomicWrites ever has. Should simply fail instead.
func TestExecRead(t *testing.T){
	cmdString := "./test-atomic-writes"
	cmdArgs := []string{"-readonly"}
    cmd := cmdMake(cmdString, cmdArgs)
	err := cmd.Start()
	if err != nil { t.Errorf("%s cmd start failed: %v\n", cmdString, err) }
	err = cmd.Wait()
	if err != nil { t.Errorf("Wait failed: %v\n", err) }
}
