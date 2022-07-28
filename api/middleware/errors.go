package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/polldo/patweb/api/web"
	"github.com/sirupsen/logrus"
)

// Info indicates if the error is informative (non critical)
// and it should be logged with Info level rather than Err.
// If the error does not implement the Info behavior, it returns
// false.
func Info(err error) bool {
	var infoErr interface{ Info() bool }
	if errors.As(err, &infoErr) {
		return infoErr.Info()
	}
	return false
}

// Fields extracts fields to be logged together with the error.
// If the error does not implement the Fields behavior, it returns
// 'ok' to false and other parameters should be ignored.
func Fields(err error) (fields map[string]interface{}, ok bool) {
	var fieldsErr interface{ Fields() map[string]interface{} }
	if errors.As(err, &fieldsErr) {
		return fieldsErr.Fields(), true
	}
	return nil, false
}

// Response returns a body and status code to use as a web response.
// If the error does not implement the Response behavior, it returns
// 'ok' to false and other parameters should be ignored.
func Response(err error) (body interface{}, code int, ok bool) {
	var respErr interface{ Response() (interface{}, int) }
	if errors.As(err, &respErr) {
		body, code := respErr.Response()
		return body, code, true
	}
	return nil, 0, false
}

// Errors handles errors coming out of the call chain.
// This middleware leverages a technique of opaque errors that
// allows to customize errors with behaviors without coupling them to
// a specific type.
// In this way, it's easier to create new errors compatible with
// the behavior used here.
func Errors(log logrus.FieldLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			err := handler(ctx, w, r)
			if err == nil {
				return nil
			}

			// Prepare fields to be logged.
			fields := map[string]interface{}{
				"req_id":  ContextRequestID(ctx),
				"message": err,
			}
			if f, ok := Fields(err); ok {
				for k, v := range f {
					fields[k] = v
				}
			}

			// Log the error with the appropriate level.
			loglvl := log.WithFields(logrus.Fields(fields)).Error
			if Info(err) {
				loglvl = log.WithFields(logrus.Fields(fields)).Info
			}
			loglvl("ERROR")

			// Try to retrieve a response from the error.
			if body, code, ok := Response(err); ok {
				return web.Respond(ctx, w, body, code)
			}

			// Unknown error, respond with Internal Server Error.
			er := struct {
				Error string `json:"error"`
			}{
				http.StatusText(http.StatusInternalServerError),
			}
			return web.Respond(ctx, w, er, http.StatusInternalServerError)
		}
		return h
	}
	return m
}