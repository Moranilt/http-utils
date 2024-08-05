package client

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

type Method string

const (
	MethodGet    Method = http.MethodGet
	MethodPost   Method = http.MethodPost
	MethodPut    Method = http.MethodPut
	MethodPatch  Method = http.MethodPatch
	MethodDelete Method = http.MethodDelete
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
	Post(ctx context.Context, url string, body []byte, headers Headers) (*http.Response, error)
	Put(ctx context.Context, url string, body []byte, headers Headers) (*http.Response, error)
	Patch(ctx context.Context, url string, body []byte, headers Headers) (*http.Response, error)
	Delete(ctx context.Context, url string, body []byte, headers Headers) (*http.Response, error)
	Get(ctx context.Context, url string, headers Headers) (*http.Response, error)

	Do(ctx context.Context, method Method, url string, body []byte, headers Headers) (*http.Response, error)
}

// Create new Client instance
func New() Client {
	return &client{
		client: &http.Client{},
	}
}

// Send request with method POST
func (c *client) Post(ctx context.Context, url string, body []byte, headers Headers) (*http.Response, error) {
	return c.Do(ctx, MethodPost, url, body, headers)
}

// Send request with method PUT
func (c *client) Put(ctx context.Context, url string, body []byte, headers Headers) (*http.Response, error) {
	return c.Do(ctx, MethodPut, url, body, headers)
}

// Send request with method PATCH
func (c *client) Patch(ctx context.Context, url string, body []byte, headers Headers) (*http.Response, error) {
	return c.Do(ctx, MethodPatch, url, body, headers)
}

// Send request with method DELETE
func (c *client) Delete(ctx context.Context, url string, body []byte, headers Headers) (*http.Response, error) {
	return c.Do(ctx, MethodDelete, url, body, headers)
}

// Send request with method GET
func (c *client) Get(ctx context.Context, url string, headers Headers) (*http.Response, error) {
	return c.Do(ctx, MethodGet, url, nil, headers)
}

func (c *client) Do(ctx context.Context, method Method, url string, body []byte, headers Headers) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, string(method), url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if headers != nil {
		for _, key := range headers.Keys() {
			request.Header.Set(key, headers.Get(key))
		}
	}

	request, _ = c.setRequestTimeout(request)
	res, err := c.client.Do(request)
	if err != nil {
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
