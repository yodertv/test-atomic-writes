// test-atomic-writes-lib_test.go

package api

import (
	"testing"
)

var (
    res = 0
    cl = Cmdline_args{}
)

// This test requires a main to parse args and perform worker tasks.
func TestMain(m *testing.M) {
    // Initialize test
    Parse_args(&cl)
    // Run all tests and capture the result, res, only in the paraent executable.
    if cl.Worker == -1 { // The orchestrating process has worker index = -1.
        res = m.Run()
    }
    return
}

func TestAtomicWrites(t *testing.T) {
    if !cl.Readonly {
        // Note that this function forks a process for each worker.
        Write_bytes(cl.Count, cl.Size, cl.Workers, cl.Worker, cl.Filename)
    }
    if cl.Worker == -1 { // The orchestrating process has worker index = -1. Only need to validate the file once after all the workers finish.
        res = Validate_bytes(cl.Filename, cl.Count, cl.Size, cl.Workers)
    }
    if (res != 0) {
        t.Errorf("Validate_bytes returned non zero status: %d.", res)
    }
}

func TestReadOnly(t *testing.T) {
    res = Validate_bytes(cl.Filename, cl.Count, cl.Size, cl.Workers)
    if (res != 0) {
        t.Errorf("Validate_bytes returned non zero status")
    }
}
