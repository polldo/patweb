package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/polldo/patweb/api/web"
)

// Demo is a simple handler that shows off the various errors behaviors.
func Demo() web.Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var input struct{ Value string }
		if err := web.Decode(r, &input); err != nil {
			err = fmt.Errorf("unable to decode payload: %w", err)
			return NewRequestError(err, http.StatusInternalServerError, WithMsg("we couldn't decode your payload!"))
		}

		switch input.Value {
		case "mask":
			// Mask the error in the response but keep it in the logs.
			err := fmt.Errorf("internal reasons here: wrapping other internal errors")
			return NewRequestError(err, http.StatusBadRequest, WithMsg("This is a bad request m8!"))

		case "dont mask":
			// Keep the whole error in the response.
			err := errors.New("internal reasons here: wrapping other internal errors")
			return NewRequestError(err, http.StatusBadRequest)

		case "log fields":
			// Add fields to error log.
			err := errors.New("add fields to log the err")
			f := map[string]interface{}{"description": "this is an additional info"}
			return NewRequestError(err, http.StatusBadRequest, WithFields(f))

		case "be quiet":
			// Mark the error as 'quiet' to log it as INFO rather than ERR.
			err := errors.New("some physiological error, logged as info")
			return NewRequestError(err, http.StatusBadRequest, Quiet())

		default:
			return web.Respond(ctx, w, struct{}{}, http.StatusOK)
		}
	}
	return h
}
