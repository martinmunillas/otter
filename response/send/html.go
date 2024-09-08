package send

import (
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

func (h htmlSender) send(w http.ResponseWriter, r *http.Request, component templ.Component, status int) {
	w.WriteHeader(status)
	err := component.Render(r.Context(), w)
	if err != nil {
		h.logger.Error(err.Error())
	}
}

func (h htmlSender) Ok(w http.ResponseWriter, r *http.Request, component templ.Component) {
	h.send(w, r, component, http.StatusOK)
}

func (h htmlSender) InternalError(w http.ResponseWriter, r *http.Request, err error, component templ.Component) {
	h.logger.Error(err.Error())
	h.send(w, r, component, http.StatusInternalServerError)
}

func (h htmlSender) Unauthorized(w http.ResponseWriter, r *http.Request, component templ.Component) {
	h.send(w, r, component, http.StatusUnauthorized)
}

func (h htmlSender) Forbidden(w http.ResponseWriter, r *http.Request, component templ.Component) {
	h.send(w, r, component, http.StatusForbidden)
}

func (h htmlSender) NotFound(w http.ResponseWriter, r *http.Request, component templ.Component) {
	h.send(w, r, component, http.StatusNotFound)
}

func (h htmlSender) BadRequest(w http.ResponseWriter, r *http.Request, component templ.Component) {
	h.send(w, r, component, http.StatusBadRequest)
}
