package send

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

type htmlSender struct {
	logger *slog.Logger
}

func (h *htmlSender) SetLogger(logger *slog.Logger) {
	h.logger = logger
}

func (h htmlSender) send(w http.ResponseWriter, ctx context.Context, component templ.Component, status int) {
	w.WriteHeader(status)
	if component == nil {
		return
	}
	err := component.Render(ctx, w)
	if err != nil {
		h.logger.Error(err.Error())
	}
}

func (h htmlSender) Ok(w http.ResponseWriter, ctx context.Context, component templ.Component) {
	h.send(w, ctx, component, http.StatusOK)
}

func (h htmlSender) InternalError(w http.ResponseWriter, ctx context.Context, err error, component templ.Component) {
	if err != nil {
		h.logger.Error(err.Error())
	}
	h.send(w, ctx, component, http.StatusInternalServerError)
}

func (h htmlSender) Unauthorized(w http.ResponseWriter, ctx context.Context, component templ.Component) {
	h.send(w, ctx, component, http.StatusUnauthorized)
}

func (h htmlSender) Forbidden(w http.ResponseWriter, ctx context.Context, component templ.Component) {
	h.send(w, ctx, component, http.StatusForbidden)
}

func (h htmlSender) NotFound(w http.ResponseWriter, ctx context.Context, component templ.Component) {
	h.send(w, ctx, component, http.StatusNotFound)
}

func (h htmlSender) BadRequest(w http.ResponseWriter, ctx context.Context, component templ.Component) {
	h.send(w, ctx, component, http.StatusBadRequest)
}
