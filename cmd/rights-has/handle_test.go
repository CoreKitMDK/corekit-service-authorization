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

func TestHandle(t *testing.T) {
	entityID := uuid.New()

	hasRightsReq := authorization.HasRightsRequest{
		Entity:    entityID,
		Resources: []string{"resource1", "resource2"},
	}

	reqBody, err := json.Marshal(hasRightsReq)
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
	fmt.Println(string(body))

	if res.StatusCode != 200 {
		t.Fatalf("unexpected response code: %v", res.StatusCode)
	}

	time.Sleep(5 * time.Second)
}
