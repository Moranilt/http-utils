package client

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

var timeout atomic.Value

func init() {
	timeout.Store(10 * time.Second)
}

// Set request timeout
func SetTimeout(val time.Duration) {
	timeout.Store(val)
}

// Get request timeout
func Timeout() time.Duration {
	return timeout.Load().(time.Duration)
}

type client struct {
	client *http.Client
}

type Client interface {
	Post(ctx context.Context, url string, body []byte, headers map[string]string) (*http.Response, error)
	Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error)
}

// Create new Client instance
func New() Client {
	return &client{
		client: &http.Client{},
	}
}

// Send request with method POST
func (c *client) Post(ctx context.Context, url string, body []byte, headers map[string]string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	request, cancel := c.setRequestTimeout(request)
	res, err := c.client.Do(request)
	if err != nil {
		cancel()
		return nil, err
	}

	return res, nil
}

// Send request with method GET
func (c *client) Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	request, cancel := c.setRequestTimeout(request)
	res, err := c.client.Do(request)
	if err != nil {
		cancel()
		return nil, err
	}

	return res, nil
}

// Set request timeout
func (c *client) setRequestTimeout(req *http.Request) (*http.Request, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(req.Context(), Timeout())
	return req.WithContext(ctx), cancel
}

// Get IP address from request using X-Real-IP or X-Forwarded-For headers
func GetIP(req *http.Request) string {
	if req == nil {
		return ""
	}
	ip := req.Header.Get("X-Real-IP")
	if ip == "" {
		ip = req.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			ip = req.RemoteAddr
		}
		return ip
	}
	return ip
}
