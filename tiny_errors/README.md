# Tiny Errors
This package is needed to make your errors standardized.

The main idea is to have an array of errors with a specific codes and messages.

# Examples
## Initialization
You can init your errors globally for a project.

**main.go**:
```go
package main

import (
  "github.com/Moranilt/http-utils/tiny_errors"
)

const (
  ERR_NotValidToken = 1000
  ERR_TokenExpired = 1001

  ERR_UserNotFound = 2000
  ERR_UserAlreadyExists = 2001
)

var globalErrors = map[int]string{
  ERR_NotValidToken: "Not valid token",
  ERR_TokenExpired: "Token expired",
  ERR_UserNotFound: "User not found",
  ERR_UserAlreadyExists: "User already exists",
}

func main() {
  // ...
  tiny_errors.Init(globalErrors)

}
```

**repository.go**:
```go
func (repo *Repository) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, tiny_errors.ErrorHandler) {
  // ...

  if err != nil {
    return nil, tiny_errors.New(
      ERR_UserNotFound,
      tiny_errors.SetDetail("user_id", req.ID),
    )
  }

  // ...
}
// Output: {"error": {"code": 2000, "message": "User not found", "details": {"user_id": "1"}}}
```

You can use details to provide more information into response to help you to debug or make human readable messages with extra data from details field.

## Without initialization
You can use it without initialization as simple as you think:
```go
func (repo *Repository) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, tiny_errors.ErrorHandler) {
  // ...

  if err != nil {
    return nil, tiny_errors.New(
      ERR_UserNotFound,
      tiny_errors.Message("User not found"),
      tiny_errors.SetDetail("user_id", req.ID),
    )
  }

  // ...
}