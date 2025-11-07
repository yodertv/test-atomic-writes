// start_test.go

package test_atomic_writes

import (
    "test_atomic_writes"
	"testing"
	"net/http"
    "net/http/httptest"
)

func TestAtomicWriteHandler(t *testing.T) {
    // Create a request to pass to our handler. We don't have any query parameters for now, so we'll
    // pass 'nil' as the third parameter.
    req, err := http.NewRequest("GET", "/api/start", nil)
    if err != nil {
        t.Fatal(err)
    }
    // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
    rr := httptest.NewRecorder()
    hndlr := http.HandlerFunc(test_atomic_writes.Handler)

    // Our handlers satisfy http.Handler, so we can call their ServeHTTP method
    // directly and pass in our Request and ResponseRecorder.
    hndlr.ServeHTTP(rr, req)
    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
    // Check the response body is what we expect.
    expected := 100
    if len(rr.Body.String()) < expected {
        t.Errorf("handler returned unexpected body size: got %v want %v",
            len(rr.Body.String()), expected)
    }
    t.Log(rr.Body.String())
}
