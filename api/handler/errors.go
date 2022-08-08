package handler

import (
	"net/http"

	"github.com/polldo/patweb/api/weberr"
)

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
type RequestError struct {
	Err    error
	Msg    string
	Status int
}

// ErrOpt defines the type for RequestError options.
type ErrOpt func(*RequestError)

// WithMsg returns an option that sets the error message.
func WithMsg(msg string) ErrOpt {
	return func(err *RequestError) {
		err.Msg = msg
	}
}

// WithMsg returns an option that decorates the error
// with the 'Fields' behavior.
func WithFields(fields map[string]interface{}) ErrOpt {
	return func(err *RequestError) {
		err.Err = weberr.Wrap(err.Err, weberr.WithFields(fields))
	}
}

// WithMsg returns an option that decorates the error
// with the 'Quiet' behavior.
func WithQuiet(quiet bool) ErrOpt {
	return func(err *RequestError) {
		err.Err = weberr.Wrap(err.Err, weberr.WithQuiet(quiet))
	}
}

// NewRequestError wraps a provided error with HTTP details that can be used later on
// to build and log an appropriate HTTP error response.
//
// This function should be used when handlers encounter expected errors.
func NewRequestError(err error, status int, opts ...ErrOpt) error {
	re := &RequestError{Err: err, Status: status}
	for _, opt := range opts {
		opt(re)
	}
	msg := re.Err.Error()
	if re.Msg != "" {
		msg = re.Msg
	}
	re.Err = weberr.Wrap(re.Err, weberr.WithResponse(
		&ErrorResponse{
			Error:  msg,
			Status: http.StatusText(status),
		},
		re.Status,
	))
	return re
}

// Unwrap allows to propagate inner error behaviors.
func (e *RequestError) Unwrap() error { return e.Err }

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (r *RequestError) Error() string { return r.Err.Error() }
