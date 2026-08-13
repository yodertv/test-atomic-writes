// test-atomic-writes-lib_test.go

package lib_test

import (
    "fmt"
	"testing"

    "github.com/yodertv/test-atomic-writes/lib"
)

var (
    res = 0
    cl = lib.Cmdline_args{}
)

// This test requires a main to parse args and perform worker tasks.
func TestMain(m *testing.M) {
    fmt.Printf("Entering TestMain...\n")
    // Initialize test
    lib.Parse_args(&cl)
    // Run all tests and capture the result, res, only in the parent executable.
    if cl.Worker == -1 { // The orchestrating process has worker index = -1.
        cl.Filename = "../testdata"
        res = m.Run()
    }
    fmt.Printf("Exiting TestMain with %d result...\n", res)
    return
}

func TestAtomicWrites(t *testing.T) {
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

func TestReadOnly(t *testing.T) {
    res = lib.Validate_bytes(cl.Filename, cl.Count, cl.Size, cl.Workers)
    if (res != 0) {
        t.Errorf("Validate_bytes returned non zero status")
    }
}
