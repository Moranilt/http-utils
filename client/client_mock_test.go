package client

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/Moranilt/http-utils/mock"
	"github.com/stretchr/testify/assert"
)

var postClientTests = []struct {
	name      string
	body      []byte
	url       string
	response  *http.Response
	runExpect bool
}{
	{
		name: "default call",
		body: []byte(`{"name":"John"}`),
		url:  "http://test.com",
		response: &http.Response{
			StatusCode: http.StatusOK,
		},
		runExpect: true,
	},
	{
		name:      "unexpected call",
		body:      []byte(`{"name":"Jane"}`),
		url:       "http://test.com",
		runExpect: false,
		response: &http.Response{
			StatusCode: http.StatusOK,
		},
	},
}

var getClientTests = []struct {
	name          string
	url           string
	response      *http.Response
	runExpect     bool
	expectedError error
}{
	{
		name: "default call of Get",
		url:  "http://test.com",
		response: &http.Response{
			StatusCode: http.StatusOK,
		},
		runExpect:     true,
		expectedError: nil,
	},
	{
		name: "unexpected call of Get",
		url:  "http://test.com/users",
		response: &http.Response{
			StatusCode: http.StatusOK,
		},
		runExpect:     false,
		expectedError: fmt.Errorf(mock.ERR_Events_Is_Empty, "Get"),
	},
}

func TestMockHttpUtil(t *testing.T) {
	mockClient := NewMock()

	var testMethods = []struct {
		name       string
		method     string
		expectFunc func(url string, body []byte, err error, response *http.Response, headers Headers)
		executor   func(ctx context.Context, url string, body []byte, headers Headers) (*http.Response, error)
	}{
		{
			name:   "Post",
			method: http.MethodPost,
			expectFunc: func(url string, body []byte, err error, response *http.Response, headers Headers) {
				mockClient.ExpectPost(url, body, err, response, headers)
			},
			executor: func(ctx context.Context, url string, body []byte, headers Headers) (*http.Response, error) {
				return mockClient.Post(ctx, url, body, headers)
			},
		},
		{
			name:   "Put",
			method: http.MethodPut,
			expectFunc: func(url string, body []byte, err error, response *http.Response, headers Headers) {
				mockClient.ExpectPut(url, body, err, response, headers)
			},
			executor: func(ctx context.Context, url string, body []byte, headers Headers) (*http.Response, error) {
				return mockClient.Put(ctx, url, body, headers)
			},
		},
		{
			name:   "Patch",
			method: http.MethodPatch,
			expectFunc: func(url string, body []byte, err error, response *http.Response, headers Headers) {
				mockClient.ExpectPatch(url, body, err, response, headers)
			},
			executor: func(ctx context.Context, url string, body []byte, headers Headers) (*http.Response, error) {
				return mockClient.Patch(ctx, url, body, headers)
			},
		},
		{
			name:   "Delete",
			method: http.MethodDelete,
			expectFunc: func(url string, body []byte, err error, response *http.Response, headers Headers) {
				mockClient.ExpectDelete(url, body, err, response, headers)
			},
			executor: func(ctx context.Context, url string, body []byte, headers Headers) (*http.Response, error) {
				return mockClient.Delete(ctx, url, body, headers)
			},
		},
	}

	for _, method := range testMethods {
		t.Run(method.method, func(t *testing.T) {
			for _, test := range postClientTests {
				t.Run(test.name, func(t *testing.T) {
					expectedError := fmt.Errorf(mock.ERR_Events_Is_Empty, method.name)
					if test.runExpect {
						method.expectFunc(test.url, test.body, nil, test.response, nil)
					}

					resp, err := method.executor(context.Background(), test.url, test.body, nil)
					if err != nil {
						assert.Equal(t, expectedError.Error(), err.Error())
					}
					if resp != nil && !reflect.DeepEqual(resp, test.response) {
						assert.Equal(t, *test.response, *resp)
					}
					assert.Nil(t, mockClient.AllExpectationsDone())
					mockClient.Reset()
				})
			}
		})
	}

	t.Run("Do", func(t *testing.T) {
		for _, test := range postClientTests {
			t.Run(test.name, func(t *testing.T) {
				mockClient := NewMock()

				if test.runExpect {
					mockClient.ExpectDo(test.url, MethodPost, test.body, nil, test.response, nil)
				}

				resp, err := mockClient.Do(context.Background(), MethodPost, test.url, test.body, nil)
				if err != nil {
					assert.Equal(t, fmt.Errorf(mock.ERR_Events_Is_Empty, "Do"), err)
				}
				if resp != nil && !reflect.DeepEqual(resp, test.response) {
					assert.Equal(t, *test.response, *resp)
				}
				assert.Nil(t, mockClient.AllExpectationsDone())
				mockClient.Reset()
			})
		}
	})

	for _, test := range getClientTests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := NewMock()

			if test.runExpect {
				mockClient.ExpectGet(test.url, nil, test.response, nil)
			}

			resp, err := mockClient.Get(context.Background(), test.url, nil)
			if err != nil && err.Error() != test.expectedError.Error() {
				t.Errorf("not expected error: %v", err)
			}

			if resp != nil && !reflect.DeepEqual(resp, test.response) {
				t.Errorf("not equal post responses, expect %#v, go %#v", *test.response, *resp)
			}

			if err := mockClient.AllExpectationsDone(); err != nil {
				t.Error(err)
			}
		})
	}
}

var checkCallTests = []struct {
	name           string
	expectedURL    string
	actualURL      string
	runExpects     bool
	expectedBody   []byte
	unexpectedBody []byte
	expectedError  error
	expectedMethod Method
}{
	{
		name:           "not expected call",
		actualURL:      "http://test.com",
		runExpects:     false,
		expectedURL:    "",
		expectedBody:   []byte{},
		unexpectedBody: []byte{},
		expectedError:  fmt.Errorf(mock.ERR_Events_Is_Empty, "Post"),
		expectedMethod: MethodPost,
	},
	{
		name:           "not expected url",
		expectedURL:    "http://test.com",
		actualURL:      "http://test.com/users",
		runExpects:     true,
		expectedBody:   []byte{},
		unexpectedBody: []byte{},
		expectedError:  fmt.Errorf(ERR_Unexpected_Url, "Post", "http://test.com", "http://test.com/users"),
		expectedMethod: MethodPost,
	},
	{
		name:           "not expected body",
		expectedURL:    "http://test.com/users",
		actualURL:      "http://test.com/users",
		runExpects:     true,
		expectedBody:   []byte("expected body"),
		unexpectedBody: []byte("unexpected body"),
		expectedError:  fmt.Errorf(ERR_Unexpected_Data, "Post", "expected body", "unexpected body"),
		expectedMethod: MethodPost,
	},
	{
		name:           "not expected method",
		expectedURL:    "http://test.com/users",
		actualURL:      "http://test.com/users",
		runExpects:     true,
		expectedBody:   []byte{},
		unexpectedBody: []byte{},
		expectedError:  fmt.Errorf(ERR_Unexpected_Method, "Post", MethodPost, MethodPut),
		expectedMethod: MethodPut,
	},
}

func TestHttpCheckCall(t *testing.T) {
	for _, test := range checkCallTests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := &MockedClient{
				history: mock.NewMockHistory[*mockClientData](),
			}

			if test.runExpects {
				mockClient.ExpectPost(test.expectedURL, test.expectedBody, nil, &http.Response{
					StatusCode: http.StatusOK,
				}, nil)
			}

			item, err := mockClient.checkCall("Post", test.expectedMethod, test.actualURL, test.unexpectedBody, nil)
			if err.Error() != test.expectedError.Error() {
				t.Errorf("got error %q, expected %q", err, test.expectedError)
			}

			if item != nil {
				t.Errorf("expected item to be nil, got %#v", *item)
			}
		})
	}
}
