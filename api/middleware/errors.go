package middleware

import (
	"context"
	"net/http"

	"github.com/polldo/patweb/api/web"
	"github.com/polldo/patweb/api/weberr"
	"github.com/sirupsen/logrus"
)

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
			if f, ok := weberr.Fields(err); ok {
				for k, v := range f {
					fields[k] = v
				}
			}

			// Log the error with the appropriate level.
			loglvl := log.WithFields(logrus.Fields(fields)).Error
			if weberr.IsQuiet(err) {
				loglvl = log.WithFields(logrus.Fields(fields)).Info
			}
			loglvl("ERROR")

			// Try to retrieve a response from the error.
			if body, code, ok := weberr.Response(err); ok {
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
