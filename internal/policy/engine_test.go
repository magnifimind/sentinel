package policy

import (
	"testing"

	"github.com/magnifimind/sentinel/internal/model"
)

func TestEvaluate_BlockDestructive(t *testing.T) {
	e := NewEngine()
	tests := []struct {
		action   string
		wantDecision string
	}{
		{"rm -rf /var/log", "block"},
		{"shutdown -h now", "block"},
		{"reboot", "block"},
	}
	for _, tt := range tests {
		resp := e.Evaluate(model.PolicyEvalRequest{AgentID: "test", Action: tt.action})
		if resp.Decision != tt.wantDecision {
			t.Errorf("action %q: got %q, want %q", tt.action, resp.Decision, tt.wantDecision)
		}
	}
}

func TestEvaluate_Escalate(t *testing.T) {
	e := NewEngine()
	tests := []struct {
		action string
	}{
		{"docker rm abc123"},
		{"systemctl stop nginx"},
		{"kubectl delete pod my-pod"},
	}
	for _, tt := range tests {
		resp := e.Evaluate(model.PolicyEvalRequest{AgentID: "test", Action: tt.action})
		if resp.Decision != "escalate" {
			t.Errorf("action %q: got %q, want escalate", tt.action, resp.Decision)
		}
	}
}

func TestEvaluate_ApproveDefault(t *testing.T) {
	e := NewEngine()
	tests := []struct {
		action string
	}{
		{"df -h"},
		{"docker ps"},
		{"cat /var/log/syslog"},
		{"systemctl status nginx"},
	}
	for _, tt := range tests {
		resp := e.Evaluate(model.PolicyEvalRequest{AgentID: "test", Action: tt.action})
		if resp.Decision != "approve" {
			t.Errorf("action %q: got %q, want approve", tt.action, resp.Decision)
		}
	}
}

func TestEvaluate_CaseInsensitive(t *testing.T) {
	e := NewEngine()
	resp := e.Evaluate(model.PolicyEvalRequest{AgentID: "test", Action: "RM -RF /tmp"})
	if resp.Decision != "block" {
		t.Errorf("uppercase rm -rf: got %q, want block", resp.Decision)
	}
}
