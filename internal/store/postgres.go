package store

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magnifimind/sentinel/internal/model"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*Store, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return &Store{pool: pool}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

// RunMigrations applies the SQL migration file. Simple and sufficient for now.
func (s *Store) RunMigrations(ctx context.Context) error {
	sql, err := os.ReadFile("migrations/001_initial.sql")
	if err != nil {
		return fmt.Errorf("read migration: %w", err)
	}
	_, err = s.pool.Exec(ctx, string(sql))
	if err != nil {
		return fmt.Errorf("apply migration: %w", err)
	}
	return nil
}

// InsertDecision writes an agent decision to PostgreSQL.
func (s *Store) InsertDecision(ctx context.Context, d *model.AgentDecision) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO agent_decisions
			(timestamp, agent_id, decision_type, context, action, policy_result, success, llm_provider, tokens_used, latency_ms, model)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		d.Timestamp, d.AgentID, d.DecisionType, d.Context, d.Action,
		d.PolicyResult, d.Success, d.LLMProvider, d.TokensUsed, d.LatencyMS, d.Model,
	)
	return err
}

// InsertViolation records a governance violation.
func (s *Store) InsertViolation(ctx context.Context, v *model.GovernanceViolation) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO governance_violations (decision_id, timestamp, agent_id, action, severity, resolution)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		v.DecisionID, v.Timestamp, v.AgentID, v.Action, v.Severity, v.Resolution,
	)
	return err
}

// ListDecisions returns recent decisions, most recent first.
func (s *Store) ListDecisions(ctx context.Context, limit, offset int) ([]model.AgentDecision, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, timestamp, agent_id, decision_type, context, action, policy_result, success, llm_provider, tokens_used, latency_ms, model
		FROM agent_decisions ORDER BY timestamp DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.AgentDecision, error) {
		var d model.AgentDecision
		err := row.Scan(&d.ID, &d.Timestamp, &d.AgentID, &d.DecisionType, &d.Context,
			&d.Action, &d.PolicyResult, &d.Success, &d.LLMProvider, &d.TokensUsed, &d.LatencyMS, &d.Model)
		return d, err
	})
}

// AuditTrail returns the audit trail view entries.
func (s *Store) AuditTrail(ctx context.Context, limit, offset int) ([]model.AuditTrailEntry, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, timestamp, agent_id, decision_type, action, description, policy_result, success, llm_provider, tokens_used
		FROM audit_trail LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.AuditTrailEntry, error) {
		var e model.AuditTrailEntry
		err := row.Scan(&e.ID, &e.Timestamp, &e.AgentID, &e.DecisionType, &e.Action,
			&e.Description, &e.PolicyResult, &e.Success, &e.LLMProvider, &e.TokensUsed)
		return e, err
	})
}

// GetDecisionsByAgent returns decisions filtered by agent ID.
func (s *Store) GetDecisionsByAgent(ctx context.Context, agentID string, limit int) ([]model.AgentDecision, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, timestamp, agent_id, decision_type, context, action, policy_result, success, llm_provider, tokens_used, latency_ms, model
		FROM agent_decisions WHERE agent_id = $1 ORDER BY timestamp DESC LIMIT $2`,
		agentID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.AgentDecision, error) {
		var d model.AgentDecision
		err := row.Scan(&d.ID, &d.Timestamp, &d.AgentID, &d.DecisionType, &d.Context,
			&d.Action, &d.PolicyResult, &d.Success, &d.LLMProvider, &d.TokensUsed, &d.LatencyMS, &d.Model)
		return d, err
	})
}

// RegisterAgent upserts an agent registration.
func (s *Store) RegisterAgent(ctx context.Context, a *model.AgentRegistration) error {
	caps, _ := json.Marshal(a.Capabilities)
	_, err := s.pool.Exec(ctx,
		`INSERT INTO agent_registry (id, name, description, capabilities, registered_at, last_seen_at, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			capabilities = EXCLUDED.capabilities,
			last_seen_at = EXCLUDED.last_seen_at,
			status = EXCLUDED.status`,
		a.ID, a.Name, a.Description, caps, a.RegisteredAt, a.LastSeenAt, a.Status,
	)
	return err
}

// ListAgents returns all registered agents.
func (s *Store) ListAgents(ctx context.Context) ([]model.AgentRegistration, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, name, description, capabilities, registered_at, last_seen_at, status
		FROM agent_registry ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.AgentRegistration, error) {
		var a model.AgentRegistration
		err := row.Scan(&a.ID, &a.Name, &a.Description, &a.Capabilities, &a.RegisteredAt, &a.LastSeenAt, &a.Status)
		return a, err
	})
}
