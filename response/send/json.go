package send

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type jsonSender struct {
	logger *slog.Logger
}
type errorResponse struct {
	Error errorMessage `json:"error"`
}

func (j *jsonSender) SetLogger(logger *slog.Logger) {
	j.logger = logger
}

func (j jsonSender) sendError(w http.ResponseWriter, errResponse errorResponse) {
	w.WriteHeader(errResponse.Error.Code)
	err := json.NewEncoder(w).Encode(errResponse)
	if err != nil {
		j.logger.Error(err.Error())
	}
}

func (j jsonSender) Ok(w http.ResponseWriter, response any) {
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		j.logger.Error(err.Error())
	}
}

func (j jsonSender) InternalError(w http.ResponseWriter, err error) {
	j.logger.Error(err.Error())
	j.sendError(w, errorResponse{
		Error: errorMessage{
			Message: "Internal server error",
			Code:    http.StatusInternalServerError,
		},
	})
}

func (j jsonSender) Unauthorized(w http.ResponseWriter, message string) {
	j.sendError(w, errorResponse{
		Error: errorMessage{
			Message: message,
			Code:    http.StatusUnauthorized,
		},
	})
}

func (j jsonSender) Forbidden(w http.ResponseWriter, message string) {
	j.sendError(w, errorResponse{
		Error: errorMessage{
			Message: message,
			Code:    http.StatusForbidden,
		},
	})
}

func (j jsonSender) NotFound(w http.ResponseWriter, message string) {
	j.sendError(w, errorResponse{
		Error: errorMessage{
			Message: message,
			Code:    http.StatusNotFound,
		},
	})
}

func (j jsonSender) BadRequest(w http.ResponseWriter, message string) {
	j.sendError(w, errorResponse{
		Error: errorMessage{
			Message: message,
			Code:    http.StatusBadRequest,
		},
	})
}
