package handler

import (
	"os"
    "fmt"
    "time"
    "syscall"
    "runtime"
    "net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    currentTime := time.Now().Format(time.RFC850)
    fmt.Fprintf(w, "%v\n", currentTime)
    fmt.Fprintf(w, "%#v\n", runtime.GOOS)
    fmt.Fprintf(w, "%#v\n", os.Environ())
    fmt.Fprintf(w, "pagesize=%d\n", syscall.Getpagesize())
}
