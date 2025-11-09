// test-atomic-writes implements a multi-writer test and validates the data written
// to prove that atomic writes work as expected in append mode on a posix compliant OS.
// I.e. each write is complete and not interleved with parralel writes.

package main

import (
	"os"
	"fmt"
	"flag"

	"github.com/yodertv/test-atomic-writes/api"
)

func main() {
	var (
	    res int = 0
	    count int
	    size int
	    workers int
	    worker int
	    readonly bool
	    filename string
	)
    flag.StringVar(&filename, "f", "out.tmp", "f(ilename) to use for test")
    flag.IntVar(&count, "c", 10, "c(ount): of writes per worker, use at least 1")
    flag.IntVar(&size, "s", 30, "s(size): in bytes to write, use at least 2")
    flag.IntVar(&workers, "w", 4, "w(orkers): number of concurent writers, use at least 1 and no more than 222")
    flag.BoolVar(&readonly, "readonly", false, "readonly: just run the validate function, use identical flags")
    flag.IntVar(&worker, "worker", -1, "worker: perform the writes, internal")
    flag.Parse()
	if len(flag.Args()) != 0 { // Extraneous arguments
		flag.Usage();
		os.Exit(1)
	}
	if workers > 222 || workers < 1 || count < 1 || size < 2 {
		fmt.Println("Between 1 and 222 workers, writing at least 1 count each, of at least size 2 required to test.")
		flag.Usage();
		os.Exit(1)
	}
	if (!readonly){
		// Note that this function fork a process for each worker.
		api.Write_bytes(count, size, workers, worker, filename)
    }
    if worker == -1 { // The orchestrating process has worker index = -1. Only need to validate the file once after all the workers finish.
	    res = api.Validate_bytes(filename, count, size, workers)
    }
    os.Exit(res)
}
