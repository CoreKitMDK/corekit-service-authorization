package function

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestHandle ensures that Handle executes without error and returns the
// HTTP 200 status code indicating no errors.
func TestHandle(t *testing.T) {
	var (
		w   = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "http://example.com/test?key=foo", nil)
		res *http.Response
	)

	Handle(w, req)
	res = w.Result()
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err == nil {
		fmt.Println(string(body))
	}

	if res.StatusCode != 200 {
		t.Fatalf("unexpected response code: %v", res.StatusCode)
	}

	time.Sleep(5 * time.Second)
}
