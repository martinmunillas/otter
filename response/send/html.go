package send

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

type ErrorComponent = func(err error) templ.Component

type htmlSender struct {
	errorComponent ErrorComponent
	logInternals   bool
}

func (h *htmlSender) DisableLogInternals() {
	h.logInternals = false
}

func (h *htmlSender) SetErrorComponent(errorComponent ErrorComponent) {
	h.errorComponent = errorComponent
}

func (h htmlSender) Ok(w http.ResponseWriter, component templ.Component) error {
	err := component.Render(context.Background(), w)
	if err != nil {
		return err
	}
	return nil
}

func (h htmlSender) sendError(w http.ResponseWriter, errMessage errorMessage) error {
	w.WriteHeader(errMessage.Code)
	err := h.errorComponent(errors.New(errMessage.Message)).Render(context.Background(), w)
	if err != nil {
		return err
	}
	return nil
}

func (h htmlSender) InternalError(w http.ResponseWriter, err error) {
	if h.logInternals {
		slog.Error(err.Error())
	}
	h.sendError(w, errorMessage{
		Message: "Internal server error",
		Code:    http.StatusInternalServerError,
	})
}

func (h htmlSender) Unauthorized(w http.ResponseWriter, message string) {
	h.sendError(w, errorMessage{
		Message: message,
		Code:    http.StatusUnauthorized,
	})
}

func (h htmlSender) Forbidden(w http.ResponseWriter, message string) {
	h.sendError(w, errorMessage{
		Message: message,
		Code:    http.StatusForbidden,
	})
}

func (h htmlSender) NotFound(w http.ResponseWriter, message string) {
	h.sendError(w, errorMessage{
		Message: message,
		Code:    http.StatusNotFound,
	})
}

func (h htmlSender) BadRequest(w http.ResponseWriter, message string) {
	h.sendError(w, errorMessage{
		Message: message,
		Code:    http.StatusBadRequest,
	})
}
