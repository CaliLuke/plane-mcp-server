// Package plane provides a lightweight HTTP client for the Plane API.
package plane

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is a lightweight Plane API client.
type Client struct {
	baseURL    string // e.g. "https://api.plane.so/api/v1"
	apiKey     string
	token      string // OAuth bearer token
	httpClient *http.Client
}

// NewClient creates a new Plane API client. Provide either apiKey or accessToken, not both.
func NewClient(baseURL, apiKey, accessToken string) *Client {
	base := strings.TrimRight(baseURL, "/") + "/api/v1"
	return &Client{
		baseURL: base,
		apiKey:  apiKey,
		token:   accessToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// HttpError represents an HTTP error from the Plane API.
type HttpError struct {
	StatusCode int
	Message    string
	Payload    any
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

func (c *Client) buildURL(endpoint string) string {
	endpoint = strings.TrimLeft(endpoint, "/")
	if endpoint == "" {
		return c.baseURL + "/"
	}
	return c.baseURL + "/" + endpoint + "/"
}

func (c *Client) headers() http.Header {
	h := http.Header{}
	h.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		h.Set("X-Api-Key", c.apiKey)
	}
	if c.token != "" {
		h.Set("Authorization", "Bearer "+c.token)
	}
	return h
}

func (c *Client) do(ctx context.Context, method, endpoint string, body any, params map[string]string) (json.RawMessage, error) {
	u := c.buildURL(endpoint)

	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			q.Set(k, v)
		}
		u += "?" + q.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	for k, vals := range c.headers() {
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode == 204 || len(respBody) == 0 {
		return nil, nil
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return json.RawMessage(respBody), nil
	}

	var payload any
	if json.Unmarshal(respBody, &payload) != nil {
		payload = string(respBody)
	}
	return nil, &HttpError{
		StatusCode: resp.StatusCode,
		Message:    fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status),
		Payload:    payload,
	}
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, endpoint string, params map[string]string) (json.RawMessage, error) {
	return c.do(ctx, http.MethodGet, endpoint, nil, params)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, endpoint string, body any) (json.RawMessage, error) {
	return c.do(ctx, http.MethodPost, endpoint, body, nil)
}

// Patch performs a PATCH request.
func (c *Client) Patch(ctx context.Context, endpoint string, body any) (json.RawMessage, error) {
	return c.do(ctx, http.MethodPatch, endpoint, body, nil)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, endpoint string) error {
	_, err := c.do(ctx, http.MethodDelete, endpoint, nil, nil)
	return err
}
