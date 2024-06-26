package handler

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/Moranilt/http-utils/logger"
	"github.com/Moranilt/http-utils/response"
	"github.com/Moranilt/http-utils/tiny_errors"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

const (
	ErrNotValidBodyFormat = "unable to unmarshal request body "
	ErrEmptyMultipartData = "empty multipart form data "
)

const (
	ERR_CODE_UnexpectedBody = 999
)

type HandlerMaker[ReqT any, RespT any] struct {
	request     *http.Request
	response    http.ResponseWriter
	requestBody ReqT
	logger      logger.Logger
	caller      CallerFunc[ReqT, RespT]
	err         tiny_errors.ErrorHandler
}

// A function that is called to process request.
//
// ReqT - type of request body
// RespT - type of response body
type CallerFunc[ReqT any, RespT any] func(ctx context.Context, req ReqT) (RespT, tiny_errors.ErrorHandler)

// Create new handler instance
//
// **caller** should be a function that implements type CallerFunc[ReqT, RespT]
func New[ReqT any, RespT any](w http.ResponseWriter, r *http.Request, logger logger.Logger, caller CallerFunc[ReqT, RespT]) *HandlerMaker[ReqT, RespT] {
	log := logger.WithRequestInfo(r)
	return &HandlerMaker[ReqT, RespT]{
		logger:   log,
		request:  r,
		caller:   caller,
		response: w,
	}
}

func (h *HandlerMaker[ReqT, RespT]) setError(errs ...string) {
	h.err = tiny_errors.New(ERR_CODE_UnexpectedBody, tiny_errors.Message(strings.Join(errs, ",")))
}

// Parsing JSON-body of request.
//
// # Request type should include fields with tags of json
//
// Example:
//
//	type YourRequest struct {
//			FieldName string `json:"field_name"`
//	}
func (h *HandlerMaker[ReqT, RespT]) WithJSON() *HandlerMaker[ReqT, RespT] {
	if h.err != nil {
		return h
	}
	if h.request.Method == http.MethodGet {
		return h
	}
	err := json.NewDecoder(h.request.Body).Decode(&h.requestBody)
	if err != nil {
		h.setError(ErrNotValidBodyFormat, err.Error())
		return h
	}
	return h
}

// Parsing URI vars using gorilla/mux
//
// Request type should include fields with tags of mapstructure.
//
// Example:
//
//	type YourRequest struct {
//			FieldName string `mapstructure:"field_name"`
//	}
func (h *HandlerMaker[ReqT, RespT]) WithVars() *HandlerMaker[ReqT, RespT] {
	if h.err != nil {
		return h
	}
	vars := mux.Vars(h.request)
	cfg := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &h.requestBody,
	}

	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		h.setError(ErrNotValidBodyFormat, err.Error())
		return h
	}

	err = decoder.Decode(vars)
	if err != nil {
		h.setError(ErrNotValidBodyFormat, err.Error())
		return h
	}

	return h
}

// Parsing URL-query params from request.
//
// Request type should include fields with tags of mapstructure.
//
// Example:
//
//	type YourRequest struct {
//			FieldName string `mapstructure:"field_name"`
//	}
func (h *HandlerMaker[ReqT, RespT]) WithQuery() *HandlerMaker[ReqT, RespT] {
	if h.err != nil {
		return h
	}
	query := h.request.URL.Query()
	if len(query) == 0 {
		return h
	}

	queryVars := make(map[string]any)
	for name, q := range query {
		queryVars[name] = q[0]
	}

	cfg := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &h.requestBody,
	}

	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		h.setError(ErrNotValidBodyFormat, err.Error())
		return h
	}

	err = decoder.Decode(queryVars)
	if err != nil {
		h.setError(ErrNotValidBodyFormat, err.Error())
		return h
	}

	return h
}

// Parsing multipart-data from request body.
//
// Request type should include fields with tags of mapstructure.
//
// If field is an array of files you should set tag name as files[] and type []*multipart.FileHeader([mime/multipart.FileHeader])
//
// If field is file and not array of files you should set tag with field name without brackets and type *multipart.FileHeader([mime/multipart.FileHeader])
//
// Other fields should have string type([mime/multipart.Form])
//
// # File types
//
//   - []*multipart.FileHeader -	field with array of files. Should contain square brackets in name
//   - *multipart.FileHeader -	field with single file. Should not contain square brackets in field name
//
// Example
//
//	type YourRequest struct {
//		MultipleFiles []*multipart.FileHeader `mapstructure:"multiple_files[]"`
//		SingleFile *multipart.FileHeader 	`mapstructure:"single_file"`
//		Name string `mapstructure:"name"`
//	}
//
// # Supported nested structures.
//
// Example:
//
//	type Recipient struct {
//		Name string `json:"name,omitempty" mapstructure:"name"`
//		Age  string `json:"age,omitempty" mapstructure:"age"`
//	}
//
//	type CreateOrder struct {
//		Recipient  Recipient               `json:"recipient" mapstructure:"recipient"`
//		Content    map[string]string       `json:"content" mapstructure:"content"`
//	}
//
// Request body(multipart-form):
//
//	{
//		"recipient[name]": "John",
//		"recipient[age]": "30",
//		"content[title]": "content title",
//		"content[body]": "content body"
//	}
//
// Result:
//
//	func main() {
//		// ...
//		var order CreateOrder
//		fmt.Println(order.Recipient.Name) // John
//		fmt.Println(order.Recipient.Age) // 30
//		fmt.Println(order.Content["title"]) // content title
//		fmt.Println(order.Content["body"]) // content body
//	}
func (h *HandlerMaker[ReqT, RespT]) WithMultipart(maxMemory int64) *HandlerMaker[ReqT, RespT] {
	if h.err != nil {
		return h
	}
	if h.request.Method == http.MethodGet {
		return h
	}
	err := h.request.ParseMultipartForm(maxMemory)
	if err != nil {
		h.setError(err.Error())
		return h
	}

	if len(h.request.MultipartForm.Value) == 0 && len(h.request.MultipartForm.File) == 0 {
		h.setError(ErrEmptyMultipartData)
		return h
	}

	result := make(map[string]any, len(h.request.MultipartForm.Value)+len(h.request.MultipartForm.File))
	for name, value := range h.request.MultipartForm.Value {
		if len(value) > 0 {
			fieldName, subName, validName := extractSubName(name)
			if validName {
				if _, ok := result[fieldName]; !ok {
					result[fieldName] = make(map[string]any)
				}
				result[fieldName].(map[string]any)[subName] = value[0]
			} else {
				result[name] = value[0]
			}
		}
	}

	for name, value := range h.request.MultipartForm.File {
		if len(value) > 0 {
			fieldName, validName := extractArrayName(name)
			if validName {
				safeValue := make([]*multipart.FileHeader, 0)
				safeValue = append(safeValue, value...)
				result[fieldName] = safeValue
			} else {
				result[name] = value[0]
			}
		}
	}

	cfg := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &h.requestBody,
	}

	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		h.setError(ErrNotValidBodyFormat, err.Error())
		return h
	}

	err = decoder.Decode(result)
	if err != nil {
		h.setError(ErrNotValidBodyFormat, err.Error())
		return h
	}

	return h
}

// Run handler and send response with status code
func (h *HandlerMaker[ReqT, RespT]) Run(successStatus int) {
	h.logger.With("body", h.requestBody).Info("request")
	if h.err != nil {
		h.logger.Error(h.err.Error(), "code", h.err.GetCode(), "details", h.err.GetDetails())
		response.ErrorResponse(h.response, h.err, http.StatusBadRequest)
		return
	}

	resp, err := h.caller(h.request.Context(), h.requestBody)
	if err != nil {
		h.logger.Error(err.Error(), "code", err.GetCode(), "details", err.GetDetails())
		response.ErrorResponse(h.response, err, err.GetHTTPStatus())
		return
	}
	response.SuccessResponse(h.response, resp, successStatus)
}
