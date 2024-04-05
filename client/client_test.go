package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client
	c := New()
	// Set request timeout
	SetTimeout(1 * time.Second)

	var tests = []struct {
		name     string
		method   string
		executor func(ctx context.Context, url string, body []byte, headers map[string]string) (*http.Response, error)
	}{
		{
			name:     "Post",
			method:   http.MethodPost,
			executor: c.Post,
		},
		{
			name:     "Put",
			method:   http.MethodPut,
			executor: c.Put,
		},
		{
			name:     "Patch",
			method:   http.MethodPatch,
			executor: c.Patch,
		},
		{
			name:     "Delete",
			method:   http.MethodDelete,
			executor: c.Delete,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create request
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			reqHeaders := map[string]string{
				"Content-Type": "application/json",
			}

			resp, err := test.executor(ctx, server.URL, []byte(`{"key":"value"}`), reqHeaders)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, test.method, resp.Request.Method)
			assert.Equal(t, "application/json", resp.Request.Header.Get("Content-Type"))
		})
	}
}

func TestGet(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client
	c := New()

	// Create request
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	reqHeaders := map[string]string{
		"Accept": "application/json",
	}

	// Call Get method
	resp, err := c.Get(ctx, server.URL, reqHeaders)
	if err != nil {
		t.Fatal(err)
	}

	// Verify response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %d", resp.StatusCode)
	}

	if resp.Request.Method != http.MethodGet {
		t.Errorf("Expected method GET, got %s", resp.Request.Method)
	}

	if resp.Request.Header.Get("Accept") != "application/json" {
		t.Errorf("Expected Accept application/json, got %s", resp.Request.Header.Get("Accept"))
	}
}

func TestGetWithTimeout(t *testing.T) {
	// Create a mock server with a delay longer than the timeout
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client
	c := New()

	// Create request with a timeout shorter than the server delay
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	reqHeaders := map[string]string{
		"Accept": "application/json",
	}

	// Call Get method
	_, err := c.Get(ctx, server.URL, reqHeaders)

	// Verify that the request times out
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestGetIP(t *testing.T) {
	testCases := []struct {
		name     string
		request  *http.Request
		expected string
	}{
		{
			name:     "X-Real-IP header",
			request:  httptest.NewRequest("GET", "/", nil).WithContext(context.TODO()),
			expected: "192.168.0.1",
		},
		{
			name:     "X-Forwarded-For header",
			request:  httptest.NewRequest("GET", "/", nil).WithContext(context.TODO()),
			expected: "203.0.113.195",
		},
		{
			name:     "RemoteAddr with port",
			request:  httptest.NewRequest("GET", "/", nil).WithContext(context.TODO()),
			expected: "127.0.0.1",
		},
		{
			name:     "RemoteAddr without port",
			request:  httptest.NewRequest("GET", "/", nil).WithContext(context.TODO()),
			expected: "::1",
		},
		{
			name:     "nil request",
			request:  nil,
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.request != nil {
				switch tc.name {
				case "X-Real-IP header":
					tc.request.Header.Set("X-Real-IP", "192.168.0.1")
				case "X-Forwarded-For header":
					tc.request.Header.Set("X-Forwarded-For", "203.0.113.195")
				case "RemoteAddr with port":
					tc.request.RemoteAddr = "127.0.0.1:1234"
				case "RemoteAddr without port":
					tc.request.RemoteAddr = "::1"
				}
			}

			result := GetIP(tc.request)
			assert.Equal(t, tc.expected, result)
		})
	}
}
