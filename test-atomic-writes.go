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
    test_atomic_writes.Write_bytes(count, size, workers, worker, false, filename)
    res := test_atomic_writes.Validate_bytes(filename, count, size, workers)
    os.Exit(res)
}
