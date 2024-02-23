# Client
It's an implementation of [http.Client](https://pkg.go.dev/net/http#Client)

You can use it to make requests to the server and easy to mock in tests using `client_mock` package.

## Example of usage

**repository.go**:
```go
type Repository struct {
  client Client
}

func New(client Client) *Repository {
  return &Repository{
    client: client,
  }
}

func (repo *Repository) CustomFunction(ctx context.Context, userName string) (string, error) {
  response, err := repo.client.Post(context.Background(), "http://test.com/users", []byte(userName), map[string]string{
    "Content-Type": "application/json",
  })
  if err != nil {
    return "", fmt.Errorf("error: %v", err)
  }
  return http.StatusText(response.StatusCode), nil
}
```

**repository_test.go**:
```go
func TestCustomFunction(t *testing.T) {
  var (
    url      = "http://test.com/users"
    body     = []byte("test")
    err      error
    response = &http.Response{
      StatusCode: http.StatusOK,
    }
    headers = map[string]string{
      "Content-Type": "application/json",
    }

    mockedClient = client_mock.New()
    mockedRepo = New(mockedClient)
  )

  mockedClient.ExpectPost(url, body, err, response, headers)

  result, err := mockedRepo.CustomFunction(context.Background(), "test")
  if err != nil {
    t.Error(err)
  }

  if result != http.StatusText(http.StatusOK) {
    t.Errorf("not equal responses, expect %s, go %s", http.StatusText(http.StatusOK), result)
  }

  if err := mockedClient.AllExpectationsDone(); err != nil {
    t.Error(err)
  }
}
```

