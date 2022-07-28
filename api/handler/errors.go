package handler

import "net/http"

// ErrorResponse is the structure of error responses returned to clients:
// { "status": "Internal Server Error", "error": "some error message" } .
type ErrorResponse struct {
	Error  string `json:"error"`
	Status string `json:"status"`
}

// RequestError is used to pass an error during the request through the
// application with web specific context.
// RequestError wraps a provided error with HTTP details that can be used later on
// to build an appropriate HTTP error response.
// 'Err' is the complete error description that will be logged, it will be returned in HTTP response if 'Msg' is empty.
// 'Msg' is the text error that, if not empty, will be returned in the HTTP response.
// 'Status' indicates the status code of the response to be built.
// 'LogInfo' determines whether the error should be logged with the Info level.
// 'LogFields' contains the fields to be logged together with the error.
type RequestError struct {
	Err       error
	Msg       string
	Status    int
	LogInfo   bool
	LogFields map[string]interface{}
}

// ErrOpt defines the type for RequestError options.
type ErrOpt func(*RequestError)

func WithMsg(msg string) ErrOpt {
	return func(err *RequestError) {
		err.Msg = msg
	}
}

// Info sets the error as informative (to be logged with Info level).
func Info() ErrOpt {
	return func(err *RequestError) {
		err.LogInfo = true
	}
}

// Fields sets the error as empty.
func Fields(f map[string]interface{}) ErrOpt {
	return func(err *RequestError) {
		err.LogFields = f
	}
}

// NewRequestError wraps a provided error with HTTP details that can be used later on
// to build and log an appropriate HTTP error response.
//
// This function should be used when handlers encounter expected errors.
func NewRequestError(err error, status int, opts ...ErrOpt) error {
	e := &RequestError{Err: err, Status: status}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (r *RequestError) Error() string {
	return r.Err.Error()
}

// Response converts and returns the error in a body and status code
// to be written as response to vernemq.
func (r *RequestError) Response() (body interface{}, code int) {
	err := r.Err.Error()
	if r.Msg != "" {
		err = r.Msg
	}
	return &ErrorResponse{
		Error:  err,
		Status: http.StatusText(r.Status),
	}, r.Status
}

// Info indicates whether the error is only informative (not critical)
// and should be logged with Info level.
func (r *RequestError) Info() bool {
	return r.LogInfo
}

// Fields returns the fields to be logged together with the error.
func (r *RequestError) Fields() map[string]interface{} {
	return r.LogFields
}
