package policy

import (
	"strings"
	"sync"

	"github.com/magnifimind/sentinel/internal/model"
)

// Rule defines a single policy rule evaluated against agent actions.
type Rule struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	// ActionPattern is matched against the requested action (prefix match).
	ActionPattern string `json:"action_pattern"`
	// Decision is the outcome when this rule matches: "block" or "escalate".
	// If no rule matches, the default decision is "approve".
	Decision string `json:"decision"`
	// Severity used when creating a governance violation for blocked actions.
	Severity string `json:"severity"`
}

// Engine evaluates policy rules against agent action requests.
type Engine struct {
	mu    sync.RWMutex
	rules []Rule
}

func NewEngine() *Engine {
	return &Engine{
		rules: defaultRules(),
	}
}

// Evaluate checks the request against all rules. First matching rule wins.
// If no rule matches, the action is approved.
func (e *Engine) Evaluate(req model.PolicyEvalRequest) model.PolicyEvalResponse {
	e.mu.RLock()
	defer e.mu.RUnlock()

	action := strings.ToLower(req.Action)
	for _, r := range e.rules {
		if strings.HasPrefix(action, strings.ToLower(r.ActionPattern)) {
			return model.PolicyEvalResponse{
				Decision: r.Decision,
				Reason:   r.Description,
				PolicyID: r.ID,
			}
		}
	}

	return model.PolicyEvalResponse{
		Decision: "approve",
		Reason:   "no matching policy rule — action permitted",
	}
}

// Rules returns a copy of the current rule set.
func (e *Engine) Rules() []Rule {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]Rule, len(e.rules))
	copy(out, e.rules)
	return out
}

// SetRules replaces the rule set (for future YAML/API-driven config).
func (e *Engine) SetRules(rules []Rule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules = rules
}

// defaultRules returns the initial built-in policy rules.
// These demonstrate the governance model — block destructive ops, escalate risky ones.
func defaultRules() []Rule {
	return []Rule{
		{
			ID:            "block-rm-rf",
			Description:   "block recursive force-delete operations",
			ActionPattern: "rm -rf",
			Decision:      "block",
			Severity:      "critical",
		},
		{
			ID:            "block-shutdown",
			Description:   "block system shutdown/reboot commands",
			ActionPattern: "shutdown",
			Decision:      "block",
			Severity:      "critical",
		},
		{
			ID:            "block-reboot",
			Description:   "block system reboot commands",
			ActionPattern: "reboot",
			Decision:      "block",
			Severity:      "critical",
		},
		{
			ID:            "escalate-docker-rm",
			Description:   "escalate container removal — requires human approval",
			ActionPattern: "docker rm",
			Decision:      "escalate",
			Severity:      "high",
		},
		{
			ID:            "escalate-systemctl-stop",
			Description:   "escalate service stop — requires human approval",
			ActionPattern: "systemctl stop",
			Decision:      "escalate",
			Severity:      "high",
		},
		{
			ID:            "escalate-kubectl-delete",
			Description:   "escalate kubernetes resource deletion — requires human approval",
			ActionPattern: "kubectl delete",
			Decision:      "escalate",
			Severity:      "high",
		},
	}
}
