package model

import (
	"encoding/json"
	"time"
)

// AgentDecision represents a single agent decision record in the audit trail.
type AgentDecision struct {
	ID           int64           `json:"id"`
	Timestamp    time.Time       `json:"timestamp"`
	AgentID      string          `json:"agent_id"`
	DecisionType string          `json:"decision_type"`
	Context      json.RawMessage `json:"context"`
	Action       string          `json:"action"`
	PolicyResult string          `json:"policy_result"`
	Success      bool            `json:"success"`
	LLMProvider  string          `json:"llm_provider,omitempty"`
	TokensUsed   int             `json:"tokens_used,omitempty"`
	LatencyMS    int             `json:"latency_ms,omitempty"`
	Model        string          `json:"model,omitempty"`
}

// GovernanceViolation records a policy violation by an agent.
type GovernanceViolation struct {
	ID         int64     `json:"id"`
	DecisionID int64     `json:"decision_id"`
	Timestamp  time.Time `json:"timestamp"`
	AgentID    string    `json:"agent_id"`
	Action     string    `json:"action"`
	Severity   string    `json:"severity"`
	Resolution string    `json:"resolution,omitempty"`
}

// AgentRegistration holds metadata about a registered agent.
type AgentRegistration struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description,omitempty"`
	Capabilities json.RawMessage `json:"capabilities,omitempty"`
	RegisteredAt time.Time       `json:"registered_at"`
	LastSeenAt   time.Time       `json:"last_seen_at"`
	Status       string          `json:"status"`
}

// PolicyEvalRequest is the payload for POST /api/v1/policy/evaluate.
type PolicyEvalRequest struct {
	AgentID  string          `json:"agent_id"`
	Action   string          `json:"action"`
	Context  json.RawMessage `json:"context,omitempty"`
}

// PolicyEvalResponse is the response from policy evaluation.
type PolicyEvalResponse struct {
	Decision string `json:"decision"` // "approve", "block", "escalate"
	Reason   string `json:"reason"`
	PolicyID string `json:"policy_id,omitempty"`
}

// KafkaEvent is the envelope for events consumed from Kafka topics.
type KafkaEvent struct {
	Topic     string          `json:"topic"`
	AgentID   string          `json:"agent_id"`
	Timestamp time.Time       `json:"timestamp"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
}

// AuditTrailEntry is the view model returned by audit trail queries.
type AuditTrailEntry struct {
	ID           int64           `json:"id"`
	Timestamp    time.Time       `json:"timestamp"`
	AgentID      string          `json:"agent_id"`
	DecisionType string          `json:"decision_type"`
	Action       string          `json:"action"`
	Description  string          `json:"description"`
	PolicyResult string          `json:"policy_result"`
	Success      bool            `json:"success"`
	LLMProvider  string          `json:"llm_provider,omitempty"`
	TokensUsed   int             `json:"tokens_used,omitempty"`
}
