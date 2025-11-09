// test-atomic-writes.go implements a multi-writer test and validates the data written
// to prove that atomic writes work as expected in append mode on a posix compliant OS.
// I.e. each write is complete and not interleved with parralel writes.

package api

import (
	"os"
	"fmt"
	"syscall"
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

func Validate_bytes(path string, count int, size int, workers int) int {
	// Read all data as bytes and validate atomic write pattern expected.
	stat := syscall.Stat_t {}
	err := error(nil)
	fd := int(0)
	fd, err = open_file(path)
	check(err)
	err = syscall.Fstat(fd, &stat)
	check(err)
	filesize := int(stat.Size)
	efs := count * size * workers
	emc := count * workers
	d, err := map_log(fd, int(filesize))
	check(err)
	lastByte := byte(d[0])
	lastMsg := lastByte
	msgSrcChangeCnt := 1
	msgCnt := 0
	errCnt := 0
	byteCount := 1
	if filesize != efs {
		errCnt++
		fmt.Printf("validate: filesize missmatch: filesize=%d expected file size=%d\n", filesize, efs)
	}
	for i := 1 ; i < filesize ; i++ {
		if lastByte == d[i] {
			byteCount++
		} else if d[i] != '\n' {
			errCnt++
			fmt.Printf("validate: interleaved write: errCnt=%d byteCount=%d, size%d, msgSrcChangeCnt=%d\n", errCnt, byteCount, size, msgSrcChangeCnt)
		} else {
			i++
			msgCnt++
			byteCount++
			if byteCount != size {
				errCnt++
				fmt.Printf("validate: bad messages size: byteCount=%d, size=%d\n", byteCount, size)
			}
			if i == filesize { break }
			if lastMsg != d[i] {
				msgSrcChangeCnt++
				lastMsg = d[i]
			}
			byteCount = 1
		}
		lastByte = d[i]
		lastMsg = lastByte
	}
	if msgCnt != emc {
		fmt.Printf("validate: missed messages: msgCnt=%d expected=%d difference=%d\n", msgCnt, emc, emc - msgCnt)
		errCnt++
	}
	fmt.Printf("validate: source changes=%d shuffle=%.2f msgCnt=%d errCnt=%d\n", msgSrcChangeCnt, float32(msgSrcChangeCnt)/float32(workers), msgCnt, errCnt)
	return (errCnt)
}

func close_log (fd int, file string){
	syscall.Close(fd)
	// Leave the file for examination post test.
	// syscall.Unlink(file)
}

func Write_bytes (count int, size int, workers int, worker int, filename string) {
	// As a worker append cnt messages unique to the worker.
	if worker > -1 {
		fd, err := open_file(filename)
		check(err)
		defer close_log(fd, filename)
		msg := make([]byte, size)
		for i := 0; i < size; i++ {
			msg[i] = byte('!' + worker)
		}
		msg[size-1] = '\n'
		for i := 0 ; i < count ; i++ {
			nw, err := append_msg(fd, []byte(msg))
			check(err)
			if nw != size {
				err = fmt.Errorf("bytes writen don't match len: %d <> %d", size, nw,)
				check(err)
			}
		}
	} else {
		// As the parent start the workers with shared IO, wait for them, and test the results.
		err := error(nil)
		attr := syscall.ProcAttr{Dir: "", Env: nil, Files: []uintptr{ 0, 1, 2 }, Sys: nil} // stdin, out, err passed to child.
		wstatus := syscall.WaitStatus(0)
		rusage := syscall.Rusage{}
		pids := make([]int, workers)
		options := int(0)
		syscall.Unlink(filename)
		os.Args = append(os.Args, []string{ "-worker", "0" } ...)
		for i:=0 ; i < workers ; i++ {
			workerNumber := fmt.Sprintf("%d", i)
			os.Args[len(os.Args) - 1] = workerNumber
			pids[i], err = syscall.ForkExec(os.Args[0], os.Args[:], &attr)
			check(err)
		}
		fmt.Printf("Each line of file %s will be %d bytes, written by %d workers, writing %d lines each.\n", filename, size, workers, count)
		// wait for them all to finish
		for i:=0 ; i < workers ; i++ {
			_, err := syscall.Wait4(pids[i], &wstatus, options, &rusage)
			check(err)
		}
	}
}
