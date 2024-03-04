package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

	// Create request
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	reqHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Call Post method
	resp, err := c.Post(ctx, server.URL, []byte(`{"key":"value"}`), reqHeaders)
	if err != nil {
		t.Fatal(err)
	}

	// Verify response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %d", resp.StatusCode)
	}

	if resp.Request.Method != http.MethodPost {
		t.Errorf("Expected method POST, got %s", resp.Request.Method)
	}

	if resp.Request.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", resp.Request.Header.Get("Content-Type"))
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
