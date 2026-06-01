package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/magnifimind/sentinel/internal/model"
	"github.com/magnifimind/sentinel/internal/policy"
)

type PolicyHandler struct {
	engine *policy.Engine
	logger *slog.Logger
}

func NewPolicyHandler(engine *policy.Engine, logger *slog.Logger) *PolicyHandler {
	return &PolicyHandler{engine: engine, logger: logger}
}

// Evaluate handles POST /api/v1/policy/evaluate.
func (h *PolicyHandler) Evaluate(w http.ResponseWriter, r *http.Request) {
	var req model.PolicyEvalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.AgentID == "" || req.Action == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "agent_id and action are required"})
		return
	}

	resp := h.engine.Evaluate(req)
	h.logger.Info("policy evaluated",
		"agent_id", req.AgentID,
		"action", req.Action,
		"decision", resp.Decision,
		"policy_id", resp.PolicyID,
	)
	writeJSON(w, http.StatusOK, resp)
}

// ListRules handles GET /api/v1/policy/rules.
func (h *PolicyHandler) ListRules(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.engine.Rules())
}
