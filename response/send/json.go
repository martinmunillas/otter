package send

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type jsonSender struct {
	logInternals bool
}

type errorResponse struct {
	Error errorMessage `json:"error"`
}

func (j *jsonSender) DisableLogInternals() {
	j.logInternals = false
}

func (j jsonSender) sendError(w http.ResponseWriter, errResponse errorResponse) error {
	w.WriteHeader(errResponse.Error.Code)
	err := json.NewEncoder(w).Encode(errResponse)
	if err != nil {
		return err
	}
	return nil
}

func (j jsonSender) Ok(w http.ResponseWriter, response any) error {
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return j.InternalError(w, err)
	}
	return nil
}

func (j jsonSender) InternalError(w http.ResponseWriter, err error) error {
	if j.logInternals {
		slog.Error(err.Error())
	}
	return j.sendError(w, errorResponse{
		Error: errorMessage{
			Message: "Internal server error",
			Code:    http.StatusInternalServerError,
		},
	})
}

func (j jsonSender) Unauthorized(w http.ResponseWriter, message string) error {
	return j.sendError(w, errorResponse{
		Error: errorMessage{
			Message: message,
			Code:    http.StatusUnauthorized,
		},
	})
}

func (j jsonSender) Forbidden(w http.ResponseWriter, message string) error {
	return j.sendError(w, errorResponse{
		Error: errorMessage{
			Message: message,
			Code:    http.StatusForbidden,
		},
	})
}

func (j jsonSender) NotFound(w http.ResponseWriter, message string) error {
	return j.sendError(w, errorResponse{
		Error: errorMessage{
			Message: message,
			Code:    http.StatusNotFound,
		},
	})
}

func (j jsonSender) BadRequest(w http.ResponseWriter, message string) error {
	return j.sendError(w, errorResponse{
		Error: errorMessage{
			Message: message,
			Code:    http.StatusBadRequest,
		},
	})
}
