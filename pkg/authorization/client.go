package authorization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"os"
	"time"
)

type Client struct {
	httpClient http.Client
	namespace  string
	hostname   string
	Tags       map[string]string
}

func NewClient(additionalTags *map[string]string) *Client {
	namespace := ""
	if data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		namespace = string(bytes.TrimSpace(data))
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	tags := map[string]string{}
	if additionalTags != nil {
		for k, v := range *additionalTags {
			tags[k] = v
		}
	}

	id := uuid.Must(uuid.NewV7())
	tags["x-trace-uid"] = id.String()
	tags["x-trace-namespace"] = namespace
	tags["x-trace-hostname"] = hostname

	return &Client{
		httpClient: http.Client{
			Timeout: 10 * time.Second,
		},
		namespace: namespace,
		hostname:  hostname,
		Tags:      tags,
	}
}

func (c *Client) getServiceURL(functionName string) string {
	if c.namespace == "" {
		return fmt.Sprintf("http://%s", functionName)
	}
	return fmt.Sprintf("http://%s.%s", functionName, c.namespace)
}

func (c *Client) addHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Caller", c.hostname)
	for k, v := range c.Tags {
		req.Header.Set(k, v)
	}
}

// GiveRights sends a GiveRightsRequest to the rights-give Knative function.
func (c *Client) GiveRights(req *GiveRightsRequest) (*GiveRightsResponse, error) {
	url := c.getServiceURL("rights-give")

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GiveRightsRequest: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	c.addHeaders(httpReq)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code: %d", httpResp.StatusCode)
	}

	var resp GiveRightsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode GiveRightsResponse: %w", err)
	}

	return &resp, nil
}

// GetRights sends a GetRightsRequest to the rights-get Knative function.
func (c *Client) GetRights(req *GetRightsRequest) (*GetRightsResponse, error) {
	url := c.getServiceURL("rights-get")

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GetRightsRequest: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Caller", c.hostname)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code: %d", httpResp.StatusCode)
	}

	var resp GetRightsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode GetRightsResponse: %w", err)
	}

	return &resp, nil
}

// HasRights sends a HasRightsRequest to the rights-has Knative function.
func (c *Client) HasRights(req *HasRightsRequest) (*HasRightsResponse, error) {
	url := c.getServiceURL("rights-has")

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal HasRightsRequest: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Caller", c.hostname)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code: %d", httpResp.StatusCode)
	}

	var resp HasRightsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode HasRightsResponse: %w", err)
	}

	return &resp, nil
}

// RevokeRights sends a RevokeRightsRequest to the rights-revoke Knative function.
func (c *Client) RevokeRights(req *RevokeRightsRequest) (*RevokeRightsResponse, error) {
	url := c.getServiceURL("rights-revoke")

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal RevokeRightsRequest: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Caller", c.hostname)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code: %d", httpResp.StatusCode)
	}

	var resp RevokeRightsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode RevokeRightsResponse: %w", err)
	}

	return &resp, nil
}
