# Mock utility
This package can help you to create your own test mocks.

A package is an implementation of a queue.

## Example
If your repository has an http client to send requests to an external service and in tests you should mock this request without sending the actual request to the server and check that you've called the request with the expected data.


**repository.go**:
```go
type Client interface {
  Get(ctx context.Context, url string) (*http.Response, error)
}

type UserRequest struct {
  ID int
}

type UserResponse struct {
  Name string `json:"name"`
}

type Repository struct {
  client Client
}

func (r *Repository) GetUserName(ctx context.Context, req *UserRequest) (*string, error) {
  // ...
  url := fmt.Sprintf("https://test.com/user/{%v}", req.ID)
  response, err := r.client.Get(ctx, url)

  bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var userResp UserResponse

	err := json.NewDecoder(bytes).Decode(&userResp)
	if err != nil {
		return nil, err
	}

	return userResp.Name, nil
}
```

**client_mock.go**:
```go
const (
	ERR_Unexpected_Url     = "call %q expected url %q, got %q"
	ERR_Unexpected_Data    = "call %q expected data %v, got %v"
)

type MockedClient struct {
	history *mock.MockHistory[*mockClientData]
}

type mockClientData struct {
	url      string
	response *http.Response
	err      error
}

// Create new Client instance
func New() *MockedClient {
	return &MockedClient{
		history: mock.NewMockHistory[*mockClientData](),
	}
}

// Store expected call of Get method with expected data and error
func (m *MockedClient) ExpectGet(url string, err error, response *http.Response) {
	m.history.Push("Get", &mockClientData{
		url:      url,
		response: response,
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

// Check if call of Get method was expected and returning expected response and error
func (m *MockedClient) Get(ctx context.Context, url string) (*http.Response, error) {
	item, err := m.checkCall("Get", url, nil)
	if err != nil {
		return nil, err
	}

	return item.Data.response, item.Data.err
}

// Check if call of method was expected and returning expected data of this call
func (m *MockedClient) checkCall(name string, url string, body []byte) (*mock.MockHistoryItem[*mockClientData], error) {
	item, err := m.history.Get(name)
	if err != nil {
		return nil, err
	}

	if item.Data.url != url {
		return nil, fmt.Errorf(ERR_Unexpected_Url, name, item.Data.url, url)
	}

	return item, nil
}
```

**repository_test.go**:
```go
func TestRepository_GetUserName(t *testing.T) {
  var (
    ctx = context.Background()
    req = &UserRequest{
      ID: "1",
    }
    expectedName = "test"

    url = fmt.Sprintf("https://test.com/user/{%v}", req.ID)

    client = client_mock.New()
    repo = &Repository{
      client: client,
    }
  )
  
  client.ExpectGet(url, nil, &http.Response{
    Body: io.NopCloser(strings.NewReader(`{"name": "test"}`)),
    StatusCode: http.StatusOK,
  })

  name, err := repo.GetUserName(ctx, req)
  assert.NoError(t, err)
  assert.Equal(t, expectedName, *name)

  if err := client.AllExpectationsDone(); err != nil {
    t.Error(err)
  }
}
```