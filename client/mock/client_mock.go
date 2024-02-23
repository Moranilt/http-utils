package client_mock

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/Moranilt/http-utils/mock"
)

const (
	ERR_Unexpected_Url     = "call %q expected url %q, got %q"
	ERR_Unexpected_Data    = "call %q expected data %v, got %v"
	ERR_Unexpected_Headers = "call %q expected headers %v, got %v"
)

type MockedClient struct {
	history *mock.MockHistory[*mockClientData]
}

type mockClientData struct {
	url      string
	body     []byte
	response *http.Response
	headers  map[string]string
	err      error
}

// Create new Client instance
func New() *MockedClient {
	return &MockedClient{
		history: mock.NewMockHistory[*mockClientData](),
	}
}

// Store expected call of Post method with expected data and error
func (m *MockedClient) ExpectPost(url string, body []byte, err error, response *http.Response, headers map[string]string) {
	m.history.Push("Post", &mockClientData{
		url:      url,
		body:     body,
		response: response,
		headers:  headers,
		err:      err,
	}, err)
}

// Store expected call of Get method with expected data and error
func (m *MockedClient) ExpectGet(url string, err error, response *http.Response, headers map[string]string) {
	m.history.Push("Get", &mockClientData{
		url:      url,
		body:     nil,
		response: response,
		headers:  headers,
		err:      err,
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
func (m *MockedClient) Post(ctx context.Context, url string, body []byte, headers map[string]string) (*http.Response, error) {
	item, err := m.checkCall("Post", url, body, headers)
	if err != nil {
		return nil, err
	}

	return item.Data.response, item.Data.err
}

// Check if call of Get method was expected and returning expected response and error
func (m *MockedClient) Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	item, err := m.checkCall("Get", url, nil, headers)
	if err != nil {
		return nil, err
	}

	return item.Data.response, item.Data.err
}

// Check if call of method was expected and returning expected data of this call
func (m *MockedClient) checkCall(name string, url string, body []byte, headers map[string]string) (*mock.MockHistoryItem[*mockClientData], error) {
	item, err := m.history.Get(name)
	if err != nil {
		return nil, err
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
