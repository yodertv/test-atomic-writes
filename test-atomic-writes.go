// test-atomic-writes implements a multi-writer test and validates the data written
// to prove that atomic writes work as expected in append mode on a posix compliant OS.
// I.e. each write is complete and not interleved with parralel writes.

package main

import (
	"os"
	"fmt"
	"flag"
	"syscall"
	"test_atomic_writes"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func trunc_log(fd int, ln int) (err error){
	return syscall.Ftruncate(fd, int64(ln))
}

func map_log(fd int, length int) (d []byte, err error){
	var (
		off int64 = 0
		prot int = syscall.PROT_READ | syscall.PROT_WRITE
		flags int = syscall.MAP_SHARED
	)
	pgsz := int(syscall.Getpagesize())
	length = pgsz * ((length / pgsz) + 1)
	d, err = syscall.Mmap(fd, off, length, prot, flags)
	check(err)
	return
}

// Open the file in append mode for sequential, atomic writing.
func open_file(path string) (fd int, err error) {
	fd, err = syscall.Open(path, syscall.O_RDWR | syscall.O_CREAT | syscall.O_APPEND, syscall.S_IREAD | syscall.S_IWRITE )
	check(err)
	return
}

// Map the file's contents into memory begining at the 0 and extending for length.
func open_log(path string, length int) (fd int, d []byte, err error) {
	fd, _ = open_file(path)
	d, err = map_log(fd, length)
	if err != nil {
		fmt.Println("open_log: map failure ", err)
	}
	return
}

func append_msg(fd int, msg []byte) (nw int, err error){
	nw, err = syscall.Write(fd, msg)
	if err != nil {
		fmt.Printf("Write threw error: %s\n", err)
	}
	return // naked return
}

func close_log (fd int, file string){
	syscall.Close(fd)
	// Leave the file for examination post test.
	// syscall.Unlink(file)
}

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
		test_atomic_writes.Write_bytes(count, size, workers, worker, false, filename)
    }
    if worker == -1 { // The orchestrating process has worker index = -1. Only need to validate the file once after all the workers finish.
	    res = test_atomic_writes.Validate_bytes(filename, count, size, workers)
    }
    os.Exit(res)
}
