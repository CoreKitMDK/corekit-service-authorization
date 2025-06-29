package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CoreKitMDK/corekit-service-authorization/v2/pkg/authorization"
	"github.com/google/uuid"
)

// TestHandle ensures that Handle executes without error and returns the
// HTTP 200 status code indicating no errors.
func TestHandle(t *testing.T) {
	entityID := uuid.New()
	rightID1 := uuid.New()
	rightID2 := uuid.New()

	revokeRightsReq := authorization.RevokeRightsRequest{
		Entity: entityID,
		Rights: []uuid.UUID{rightID1, rightID2},
	}

	reqBody, err := json.Marshal(revokeRightsReq)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	var (
		w   = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "http://example.com/test", bytes.NewBuffer(reqBody))
		res *http.Response
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Caller", "test-caller")

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
