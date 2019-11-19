package user

import (
	"github.com/albertwidi/go_project_example/debug/user"
	"github.com/albertwidi/go_project_example/internal/pkg/context"
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
	return nil
}
