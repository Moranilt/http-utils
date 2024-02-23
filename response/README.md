# Response builder
Default response builder with `body` and `error`.

# Usage

```go
func CreateUser(w http.ResponseWriter, r *http.Request) {
  type User struct {
    Name string
    Age int
  }

  user := User{
    Name: "John",
    Age: 21,
  }

  SuccessResponse(w, user, http.StatusCreated)
}
```