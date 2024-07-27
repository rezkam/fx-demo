package echo

import (
	"io"
	"log/slog"
	"net/http"
)

// Handler is an http.Handler that copies its request body to its response body.
type Handler struct {
	logger *slog.Logger
}

// NewHandler constructs a new Handler.
func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{logger: logger}
}

// ServeHTTP handles incoming HTTP requests.
func (e *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if the request body is empty
	if r.ContentLength == 0 {
		// Respond with a 400 Bad Request status code and error message
		http.Error(w, "request body is empty", http.StatusBadRequest)
		return
	}

	// Try to copy the request body to the response body
	if _, err := io.Copy(w, r.Body); err != nil {
		// If copying fails, respond with a 500 Internal Server Error status code and error message
		http.Error(w, "failed to copy request body", http.StatusInternalServerError)
		// Log the error
		e.logger.Error("failed to handle request", "error", err)
		return
	}

	// If everything is successful, respond with a 200 OK status code
	w.WriteHeader(http.StatusOK)

	// Log the successful handling of the request
	e.logger.Info("request handled", "method", r.Method, "url", r.URL)
}

func (*Handler) Pattern() string {
	return "/echo"
}
