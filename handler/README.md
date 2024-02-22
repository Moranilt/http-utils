# HTTP handler
Package to make wrapped calls to HTTP endpoints. You can use it in service handlers.

## Examples

### Default
**service.go**:
```go
type Service interface {
	CreateUser(http.ResponseWriter, *http.Request)
	Files(w http.ResponseWriter, r *http.Request)
}

type service struct {
	log  *logger.SLogger
	repo *repository.Repository
}

func New(log *logger.SLogger, repo *repository.Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

// Parsing JSON body and binding it to CreateUserParams struct.
func (s *service) CreateUser(w http.ResponseWriter, r *http.Request) {
	handler.New(w, r, s.log, s.repo.CreateUser).
		WithJson().
		Run(http.StatusOK)
}

// Parsing multipart form and binding it to FilesParams struct.
func (s *service) Files(w http.ResponseWriter, r *http.Request) {
	handler.New(w, r, s.log, s.repo.Files).
		WithMultipart(32 << 20).
		Run(http.StatusOK)
}
```

### Extract data from multiple sources
Sometimes you should extract variables from query and router vars(gorilla/mux).

For example: you have router to get users from a specific users group `/users/{group_id}` with `limit` and `offset` query params. You want to parse all this stuff and pass to your service method:
```go
type UsersRequest struct {
  Limit   int `mapstructure:"limit"`
  Offset  int `mapstructure:"offset"`
  GroupID int `mapstructure:"group_id"`
}

func (r *Repository) GetGroupUsers(ctx context.Context, req *UsersRequest) ([]User, tiny_errors.ErrorHandler) {
  // ...

  return users, nil
}

// ...

// GetGroupUsers have request type *UsersRequest 
func (s *service) GetGroupUsers(w http.ResponseWriter, r *http.Request) {
	handler.New(w, r, s.log, s.repo.GetGroupUsers).WithQuery().WithVars().Run(http.StatusOK)
}
```

Or maybe you should extract router-vars and parse JSON-data from request body:
```go
type CreateGroupUserRequest struct {
  Name string `json:"name"`
  Description string `json:"description"`
  Avatar string `json:"avatar"`
  GroupID int `mapstructure:"group_id"`
}

func (r *Repository) CreateGroupUser(ctx context.Context, req *CreateGroupUserRequest) (User, tiny_errors.ErrorHandler) {
  // ...

  return newUser, nil
}

// ...

// CreateGroupUser have request type *CreateGroupUserRequest 
func (s *service) CreateGroupUser(w http.ResponseWriter, r *http.Request) {
	handler.New(w, r, s.log, s.repo.GetGroupUsers).WithVars().WithJSON().Run(http.StatusCreated)
}
```