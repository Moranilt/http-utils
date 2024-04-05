package client

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/Moranilt/http-utils/mock"
)

const (
	ERR_Unexpected_Method  = "call %q expected method %q, got %q"
	ERR_Unexpected_Url     = "call %q expected url %q, got %q"
	ERR_Unexpected_Data    = "call %q expected data %v, got %v"
	ERR_Unexpected_Headers = "call %q expected headers %v, got %v"
)

type MockedClient struct {
	history *mock.MockHistory[*mockClientData]
}

type Expecter interface {
	ExpectPost(url string, body []byte, err error, response *http.Response, headers Headers)
	ExpectPut(url string, body []byte, err error, response *http.Response, headers Headers)
	ExpectPatch(url string, body []byte, err error, response *http.Response, headers Headers)
	ExpectDelete(url string, body []byte, err error, response *http.Response, headers Headers)
	ExpectGet(url string, err error, response *http.Response, headers Headers)
	ExpectDo(url string, method Method, body []byte, err error, response *http.Response, headers Headers)

	AllExpectationsDone() error
	Reset()
}

type Mocker interface {
	Expecter
	Client
}

type mockClientData struct {
	url      string
	body     []byte
	response *http.Response
	headers  Headers
	err      error
	method   Method
}

// Create new Client instance
func NewMock() Mocker {
	return &MockedClient{
		history: mock.NewMockHistory[*mockClientData](),
	}
}

// Store expected call of Post method with expected data and error
func (m *MockedClient) ExpectDo(url string, method Method, body []byte, err error, response *http.Response, headers Headers) {
	m.history.Push("Do", &mockClientData{
		url:      url,
		body:     body,
		response: response,
		headers:  headers,
		err:      err,
		method:   method,
	}, err)
}

// Store expected call of Post method with expected data and error
func (m *MockedClient) ExpectPost(url string, body []byte, err error, response *http.Response, headers Headers) {
	m.history.Push("Post", &mockClientData{
		url:      url,
		body:     body,
		response: response,
		headers:  headers,
		err:      err,
		method:   MethodPost,
	}, err)
}

// Store expected call of Post method with expected data and error
func (m *MockedClient) ExpectPut(url string, body []byte, err error, response *http.Response, headers Headers) {
	m.history.Push("Put", &mockClientData{
		url:      url,
		body:     body,
		response: response,
		headers:  headers,
		err:      err,
		method:   MethodPut,
	}, err)
}

// Store expected call of Post method with expected data and error
func (m *MockedClient) ExpectPatch(url string, body []byte, err error, response *http.Response, headers Headers) {
	m.history.Push("Patch", &mockClientData{
		url:      url,
		body:     body,
		response: response,
		headers:  headers,
		err:      err,
		method:   MethodPatch,
	}, err)
}

// Store expected call of Post method with expected data and error
func (m *MockedClient) ExpectDelete(url string, body []byte, err error, response *http.Response, headers Headers) {
	m.history.Push("Delete", &mockClientData{
		url:      url,
		body:     body,
		response: response,
		headers:  headers,
		err:      err,
		method:   MethodDelete,
	}, err)
}

// Store expected call of Get method with expected data and error
func (m *MockedClient) ExpectGet(url string, err error, response *http.Response, headers Headers) {
	m.history.Push("Get", &mockClientData{
		url:      url,
		body:     nil,
		response: response,
		headers:  headers,
		err:      err,
		method:   MethodGet,
	}, err)
}

// Check if all expected calls were done
func (m *MockedClient) AllExpectationsDone() error {
	return m.history.AllExpectationsDone()
}

// Reset all expected calls
func (m *MockedClient) Reset() {
	m.history.Clear()
}

// Check if call of Post method was expected and returning expected response and error
func (m *MockedClient) Post(ctx context.Context, url string, body []byte, headers Headers) (*http.Response, error) {
	item, err := m.checkCall("Post", MethodPost, url, body, headers)
	if err != nil {
		return nil, err
	}

	return item.Data.response, item.Data.err
}

// Check if call of Post method was expected and returning expected response and error
func (m *MockedClient) Put(ctx context.Context, url string, body []byte, headers Headers) (*http.Response, error) {
	item, err := m.checkCall("Put", MethodPut, url, body, headers)
	if err != nil {
		return nil, err
	}

	return item.Data.response, item.Data.err
}

// Check if call of Post method was expected and returning expected response and error
func (m *MockedClient) Patch(ctx context.Context, url string, body []byte, headers Headers) (*http.Response, error) {
	item, err := m.checkCall("Patch", MethodPatch, url, body, headers)
	if err != nil {
		return nil, err
	}

	return item.Data.response, item.Data.err
}

// Check if call of Post method was expected and returning expected response and error
func (m *MockedClient) Delete(ctx context.Context, url string, body []byte, headers Headers) (*http.Response, error) {
	item, err := m.checkCall("Delete", MethodDelete, url, body, headers)
	if err != nil {
		return nil, err
	}

	return item.Data.response, item.Data.err
}

// Check if call of Get method was expected and returning expected response and error
func (m *MockedClient) Get(ctx context.Context, url string, headers Headers) (*http.Response, error) {
	item, err := m.checkCall("Get", MethodGet, url, nil, headers)
	if err != nil {
		return nil, err
	}

	return item.Data.response, item.Data.err
}

func (m *MockedClient) Do(ctx context.Context, method Method, url string, body []byte, headers Headers) (*http.Response, error) {
	item, err := m.checkCall("Do", method, url, body, headers)
	if err != nil {
		return nil, err
	}

	return item.Data.response, item.Data.err
}

// Check if call of method was expected and returning expected data of this call
func (m *MockedClient) checkCall(name string, method Method, url string, body []byte, headers Headers) (*mock.MockHistoryItem[*mockClientData], error) {
	item, err := m.history.Get(name)
	if err != nil {
		return nil, err
	}

	if item.Data.method != method {
		return nil, fmt.Errorf(ERR_Unexpected_Method, name, item.Data.method, method)
	}

	if item.Data.url != url {
		return nil, fmt.Errorf(ERR_Unexpected_Url, name, item.Data.url, url)
	}

	if !reflect.DeepEqual(item.Data.body, body) {
		return nil, fmt.Errorf(ERR_Unexpected_Data, name, string(item.Data.body), string(body))
	}

	if !reflect.DeepEqual(item.Data.headers, headers) {
		return nil, fmt.Errorf(ERR_Unexpected_Headers, name, item.Data.headers, headers)
	}

	return item, nil
}
