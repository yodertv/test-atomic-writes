// test-atomic-writes implements a multi-writer test and validates the data written
// to prove that atomic writes work as expected in append mode on a posix compliant OS.
// I.e. each write is complete and not interleved with parralel writes.

package main

import (
	"os"
	"flag"

	"github.com/yodertv/test-atomic-writes/api"
)

func main() {
	var (
	    res int = 0
		cl = api.Cmdline_args{}
	)
	api.Parse_args(&cl)
	if len(flag.Args()) != 0 { // Extraneous arguments
		flag.Usage();
		os.Exit(1)
	}
	if (!cl.Readonly){
		// Note that this function forks a process for each worker.
		api.Write_bytes(cl.Count, cl.Size, cl.Workers, cl.Worker, cl.Filename)
    }
    if cl.Worker == -1 { // The orchestrating process has worker index = -1. Only need to validate the file once after all the workers finish.
	    res = api.Validate_bytes(cl.Filename, cl.Count, cl.Size, cl.Workers)
    }
    os.Exit(res)
}
