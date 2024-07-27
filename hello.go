package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type HelloHandler struct {
	logger *slog.Logger
}

func NewHelloHandler(logger *slog.Logger) *HelloHandler {
	return &HelloHandler{logger: logger}
}

func (*HelloHandler) Pattern() string {
	return "/hello"
}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "failed to read request body", http.StatusInternalServerError)
		h.logger.Error("failed to read request body", "error", err)
		return
	}
	if _, err := fmt.Fprintf(w, "Hello, %s!", body); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		h.logger.Error("failed to write response", "error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
