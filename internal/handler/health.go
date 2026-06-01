package handler

import (
	"encoding/json"
	"net/http"

	"github.com/magnifimind/sentinel/internal/store"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ReadyzHandler returns a readiness check that verifies DB connectivity.
type ReadyzHandler struct {
	store *store.Store
}

func NewReadyzHandler(s *store.Store) *ReadyzHandler {
	return &ReadyzHandler{store: s}
}

func (h *ReadyzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.store.Ping(r.Context()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status": "not ready",
			"error":  "database unreachable",
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
