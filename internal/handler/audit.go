package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/magnifimind/sentinel/internal/store"
)

type AuditHandler struct {
	store  *store.Store
	logger *slog.Logger
}

func NewAuditHandler(store *store.Store, logger *slog.Logger) *AuditHandler {
	return &AuditHandler{store: store, logger: logger}
}

// ListDecisions handles GET /api/v1/audit/decisions.
func (h *AuditHandler) ListDecisions(w http.ResponseWriter, r *http.Request) {
	limit, offset := pagination(r)
	decisions, err := h.store.ListDecisions(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("list decisions failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, decisions)
}

// GetAgentDecisions handles GET /api/v1/audit/decisions/{agentID}.
func (h *AuditHandler) GetAgentDecisions(w http.ResponseWriter, r *http.Request) {
	agentID := chi.URLParam(r, "agentID")
	if agentID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "agent_id is required"})
		return
	}
	limit, _ := pagination(r)
	decisions, err := h.store.GetDecisionsByAgent(r.Context(), agentID, limit)
	if err != nil {
		h.logger.Error("get agent decisions failed", "error", err, "agent_id", agentID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, decisions)
}

// AuditTrail handles GET /api/v1/audit/trail.
func (h *AuditHandler) AuditTrail(w http.ResponseWriter, r *http.Request) {
	limit, offset := pagination(r)
	entries, err := h.store.AuditTrail(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("audit trail query failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

// ListAgents handles GET /api/v1/agents.
func (h *AuditHandler) ListAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := h.store.ListAgents(r.Context())
	if err != nil {
		h.logger.Error("list agents failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, agents)
}

func pagination(r *http.Request) (limit, offset int) {
	limit = queryInt(r, "limit", 50)
	offset = queryInt(r, "offset", 0)
	if limit > 1000 {
		limit = 1000
	}
	return
}

func queryInt(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return def
	}
	return n
}
