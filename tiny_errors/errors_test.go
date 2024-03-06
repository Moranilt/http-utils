package tiny_errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name     string
		err      ErrorHandler
		expected string
	}{
		{
			name:     "code and message",
			err:      New(2, Message("error message")),
			expected: "{\"code\":2,\"message\":\"error message\",\"details\":null}",
		},
		{
			name:     "code, message and details",
			err:      New(2, Message("error message"), Detail("name", "John")),
			expected: "{\"code\":2,\"message\":\"error message\",\"details\":{\"name\":\"John\"}}",
		},
		{
			name:     "error message",
			err:      New(2, Message("error message"), Detail("name", "John")),
			expected: "error message",
		},
		{
			name:     "message with format args",
			err:      New(2, Message("error message %s, %d", "name", 20)),
			expected: "{\"code\":2,\"message\":\"error message name, 20\",\"details\":null}",
		},
		{
			name: "details map is nil",
			err: func() *Error {
				err := &Error{Code: 1, Message: "error message"}
				option := Detail("name", "John")
				option(err)
				return err
			}(),
			expected: "error message",
		},
		{
			name: "init array of errors without message option",
			err: func() ErrorHandler {
				var (
					ErrCodeBodyRequired = 1
					errors              = map[int]string{
						ErrCodeBodyRequired: "body required",
					}
				)
				Init(errors)
				return New(ErrCodeBodyRequired)
			}(),
			expected: "{\"code\":1,\"message\":\"body required\",\"details\":null}",
		},
		{
			name: "init array of errors with message option",
			err: func() ErrorHandler {
				var (
					ErrCodeBodyRequired = 1
					errors              = map[int]string{
						ErrCodeBodyRequired: "body required",
					}
				)
				Init(errors)
				return New(ErrCodeBodyRequired, Message("message option"))
			}(),
			expected: "{\"code\":1,\"message\":\"message option\",\"details\":null}",
		},
		{
			name: "init array of errors message args",
			err: func() ErrorHandler {
				var (
					ErrCodeValidation = 2
					errors            = map[int]string{
						ErrCodeValidation: "not valid field %s",
					}
				)
				Init(errors)
				return New(ErrCodeValidation, MessageArgs("fieldname"))
			}(),
			expected: "{\"code\":2,\"message\":\"not valid field fieldname\",\"details\":null}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "error message" || tt.name == "details map is nil" {
				assert.Equal(t, tt.expected, tt.err.Error())
			} else {
				assert.Equal(t, tt.expected, tt.err.JSON())
			}
		})
	}
}

func TestError_JSONOrigin(t *testing.T) {
	err := New(1, Message("error"))

	assert.Equal(t, `{"code":1,"message":"error","details":null}`, err.JSONOrigin())
}

func TestError_GetMessage(t *testing.T) {
	err := New(1, Message("error"))

	assert.Equal(t, "error", err.GetMessage())
}

func TestError_GetCode(t *testing.T) {
	err := New(1, Message("error"))

	assert.Equal(t, 1, err.GetCode())
}

func TestError_GetHTTPStatus(t *testing.T) {
	err := New(1, HTTPStatus(400))

	assert.Equal(t, 400, err.GetHTTPStatus())
}

func TestError_GetHTTPMessage(t *testing.T) {
	err := New(1, HTTPStatus(400))

	assert.Equal(t, "Bad Request", err.GetHTTPMessage())
}

func TestError_SetCode(t *testing.T) {
	err := &Error{}
	err.SetCode(500)

	assert.Equal(t, 500, err.Code)
}

func TestError_GetDetails(t *testing.T) {
	err := New(1, Detail("name", "John"))

	assert.Equal(t, map[string]any{"name": "John"}, err.GetDetails())
}

func BenchmarkHandlerError_GoJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := &Error{
			Code:    400,
			Message: "Bad Request",
			Details: map[string]any{"field": "value"},
		}
		_ = err.JSON()
	}
}

func BenchmarkHandlerError_OriginalJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := &Error{
			Code:    400,
			Message: "Bad Request",
			Details: map[string]any{"field": "value"},
		}
		_ = err.JSONOrigin()
	}
}
