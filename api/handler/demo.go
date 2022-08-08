package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/polldo/patweb/api/web"
	"github.com/polldo/patweb/api/weberr"
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

		// Just a request error.
		case "vanilla":
			err := errors.New("internal reasons here: wrapping other internal errors")
			return NewRequestError(err, http.StatusBadRequest)

		// Mask the error in the response with a custom message but keep it in the logs.
		case "mask with msg":
			err := fmt.Errorf("internal reasons here: wrapping other internal errors")
			return NewRequestError(err, http.StatusBadRequest, WithMsg("This is a bad request m8!"))

		// Add fields to error log.
		case "log fields":
			err := errors.New("add fields to log the err")
			f := map[string]interface{}{"description": "this is an additional info"}
			return NewRequestError(err, http.StatusBadRequest, WithFields(f))

		// Mark the error as 'quiet' to log it as INFO rather than ERR.
		case "be quiet":
			err := errors.New("some physiological error, logged as info")
			return NewRequestError(err, http.StatusBadRequest, WithQuiet(true))

		// Wrap a normal error with quiet behavior.
		case "non responder but quiet error":
			err := errors.New("some normal error with quiet behavior")
			return weberr.Wrap(err, weberr.WithQuiet(true))

		default:
			return web.Respond(ctx, w, struct{}{}, http.StatusOK)
		}
	}
	return h
}
