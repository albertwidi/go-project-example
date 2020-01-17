package user

//go:generate swagger generate spec

import (
	"net/http"

	"github.com/albertwidi/go-project-example/debug/user"
	"github.com/albertwidi/go-project-example/internal/pkg/context"
)

// Handler for user debug
type Handler struct {
}

// New handler for user debug
func New(userdebug *user.DebugUsecase) *Handler {
	h := Handler{}
	return &h
}

// BypassLogin handler for bypassing user login function
func (h *Handler) BypassLogin(rctx *context.RequestContext) error {
	rctx.ResponseWriter().WriteHeader(http.StatusOK)
	rctx.ResponseWriter().Write([]byte("OK"))
	return nil
}
