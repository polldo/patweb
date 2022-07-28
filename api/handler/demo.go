package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/polldo/patweb/api/web"
)

// Demo is a simple handler that shows off the various properties of errors.
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
			err := fmt.Errorf("internal reasons here: wrapping other internal errors")
			return NewRequestError(err, http.StatusBadRequest)
		default:
			return web.Respond(ctx, w, struct{}{}, http.StatusOK)
		}
	}
	return h
}
