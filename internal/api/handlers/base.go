package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	apierrors "github.com/DARREN-2000/ai-hypervisor-platform/pkg/errors"
)

// BaseHandler provides shared response helpers for API handlers.
type BaseHandler struct {
	logger *logrus.Logger
}

// NewBaseHandler creates a base handler with a logger.
func NewBaseHandler(logger *logrus.Logger) *BaseHandler {
	return &BaseHandler{logger: logger}
}

// WriteJSON writes a JSON response with a status code.
func (h *BaseHandler) WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil && h.logger != nil {
		h.logger.WithError(err).Error("Failed to encode JSON response")
	}
}

// WriteError renders an API error payload when possible.
func (h *BaseHandler) WriteError(w http.ResponseWriter, err error) {
	if apiErr, ok := err.(*apierrors.APIError); ok {
		h.WriteJSON(w, apiErr.StatusCode, apiErr)
		return
	}

	if err != nil && h.logger != nil {
		h.logger.WithError(err).Error("Internal server error")
	}

	payload := map[string]string{
		"error":   "internal_error",
		"message": "An unexpected error occurred",
	}

	h.WriteJSON(w, http.StatusInternalServerError, payload)
}
