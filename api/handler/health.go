package handler

import (
	"context"
	"net/http"

	"github.com/polldo/patweb/api/web"
)

func Health() web.Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return web.Respond(ctx, w, struct{}{}, http.StatusOK)
	}
	return h
}
