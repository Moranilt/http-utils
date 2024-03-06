package tiny_errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"

	go_json "github.com/goccy/go-json"
)

var errorStorage atomic.Value

// initializes errorStorage to store a map that will contain
// error messages. This allows error messages to be reused for
// common error cases.
func init() {
	errorStorage.Store(make(map[int]string))
}

// Store errors globally
func Init(errors map[int]string) {
	errorStorage.Store(errors)
}

// Get global errors storage
func ErrorStorage() map[int]string {
	return errorStorage.Load().(map[int]string)
}

// Default type for error data
type Error struct {
	httpStatus  int
	httpMessage string
	Code        int            `json:"code"`
	Message     string         `json:"message"`
	Details     map[string]any `json:"details"`
}

// ErrorHandler defines an interface for handling errors that can be converted to JSON.
// It includes methods to get the HTTP status, message, error code, error message,
// and convert the error to JSON.
type ErrorHandler interface {
	JSON() string
	JSONOrigin() string
	Error() string

	GetHTTPStatus() int
	GetHTTPMessage() string
	GetCode() int
	GetMessage() string
	GetDetails() map[string]any
}

// PropertySetter defines methods to set error properties like message,
// code, details, and HTTP status. This interface allows abstracting away
// the specific error implementation.
type PropertySetter interface {
	SetMessage(string, ...any)
	FormatMessage(...any)
	SetCode(int)
	SetDetail(name, data string)
	SetHTTPStatus(int)
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.Message
}

// Converts the Error to a JSON string.
func (e *Error) JSON() string {
	b, _ := go_json.Marshal(e)
	return string(b)
}

// Converts the Error to a JSON string.
func (e *Error) JSONOrigin() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// Returns the message field of the Error.
func (e *Error) GetMessage() string {
	return e.Message
}

// Returns the code field of the Error.
func (e *Error) GetCode() int {
	return e.Code
}

// Returns the HTTP status code associated with the error.
func (e *Error) GetHTTPStatus() int {
	return e.httpStatus
}

// Returns the HTTP message associated with the error.
func (e *Error) GetHTTPMessage() string {
	return e.httpMessage
}

// Returns the details map of the Error.
func (e *Error) GetDetails() map[string]any {
	return e.Details
}

// Sets the code field of the Error.
func (e *Error) SetCode(code int) {
	e.Code = code
}

// Sets the HTTP status code and message for the error.
// The status text is looked up from the provided HTTP status code.
func (e *Error) SetHTTPStatus(code int) {
	e.httpStatus = code
	e.httpMessage = http.StatusText(code)
}

// Sets the message field of the Error, formatting it with fmt.Sprintf if
// format args are provided.
func (e *Error) SetMessage(msg string, format ...any) {
	if len(format) > 0 {
		msg = fmt.Sprintf(msg, format...)
	}
	e.Message = msg
}

// Adds a key-value pair to the Details map of the Error.
// This allows arbitrary additional context to be attached to the Error.
func (e *Error) SetDetail(name, data string) {
	if e.Details == nil {
		e.Details = make(map[string]any)
	}
	e.Details[name] = data
}

// Formats the message field of the Error using fmt.Sprintf.
// It replaces any %s and %v specifiers in the message with the provided args.
func (e *Error) FormatMessage(args ...any) {
	e.Message = fmt.Sprintf(e.Message, args...)
}

// ErrorOption is a function type that can be used to configure an Error.
// It accepts a PropertySetter func to set properties on the Error being constructed.
type ErrorOption func(PropertySetter)

// Returns an ErrorOption that sets the message field of the
// Error being constructed, formatting it with fmt.Sprintf if format
// args are provided.
func Message(message string, format ...any) ErrorOption {
	return func(err PropertySetter) {
		err.SetMessage(message, format...)
	}
}

// Returns an ErrorOption that formats the message field of the
// Error being constructed using fmt.Sprintf, replacing any %s and %v specifiers
// with the provided args.
func MessageArgs(args ...any) ErrorOption {
	return func(err PropertySetter) {
		err.FormatMessage(args...)
	}
}

// Returns an ErrorOption that adds a key-value pair to the Details map
// of the Error being constructed. This allows arbitrary additional
// context to be attached to the Error.
func Detail(name, data string) ErrorOption {
	return func(err PropertySetter) {
		err.SetDetail(name, data)
	}
}

// Returns an ErrorOption that sets the HTTPStatus field of the
// Error being constructed to the provided status code. This allows associating
// an HTTP status code with the error.
func HTTPStatus(status int) ErrorOption {
	return func(err PropertySetter) {
		err.SetHTTPStatus(status)
	}
}

// Creates a new ErrorHandler with the given error code and options.
// The code parameter specifies the error code. If a message is registered
// for the code in ErrorStorage, it will be used as the default message.
// The options allow configuring additional details like the message, HTTP
// status, and custom key-value pairs. Returns the configured ErrorHandler.
func New(code int, options ...ErrorOption) ErrorHandler {
	err := &Error{
		httpStatus:  http.StatusBadRequest,
		httpMessage: http.StatusText(http.StatusBadRequest),
		Code:        code,
	}

	if text, ok := ErrorStorage()[code]; ok {
		err.Message = text
	}

	for _, opt := range options {
		opt(err)
	}

	return err
}
