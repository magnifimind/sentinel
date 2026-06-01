-- Sentinel initial schema
-- Source: Confluence page 3309569 "Set up R740 for ML"

CREATE TABLE IF NOT EXISTS agent_registry (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    capabilities JSONB DEFAULT '{}',
    registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status      TEXT NOT NULL DEFAULT 'active'
);

CREATE TABLE IF NOT EXISTS agent_decisions (
    id            BIGSERIAL PRIMARY KEY,
    timestamp     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    agent_id      TEXT NOT NULL,
    decision_type TEXT NOT NULL,
    context       JSONB DEFAULT '{}',
    action        TEXT NOT NULL,
    policy_result TEXT NOT NULL,
    success       BOOLEAN NOT NULL DEFAULT TRUE,
    llm_provider  TEXT,
    tokens_used   INTEGER DEFAULT 0,
    latency_ms    INTEGER DEFAULT 0,
    model         TEXT
);

CREATE INDEX idx_decisions_agent_id ON agent_decisions(agent_id);
CREATE INDEX idx_decisions_timestamp ON agent_decisions(timestamp DESC);
CREATE INDEX idx_decisions_policy_result ON agent_decisions(policy_result);

CREATE TABLE IF NOT EXISTS governance_violations (
    id          BIGSERIAL PRIMARY KEY,
    decision_id BIGINT NOT NULL REFERENCES agent_decisions(id),
    timestamp   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    agent_id    TEXT NOT NULL,
    action      TEXT NOT NULL,
    severity    TEXT NOT NULL DEFAULT 'medium',
    resolution  TEXT
);

CREATE INDEX idx_violations_agent_id ON governance_violations(agent_id);
CREATE INDEX idx_violations_severity ON governance_violations(severity);

-- Audit trail view (from Confluence)
CREATE OR REPLACE VIEW audit_trail AS
SELECT
    id, timestamp, agent_id, decision_type,
    action,
    COALESCE(context->>'description', '') AS description,
    policy_result, success, llm_provider, tokens_used
FROM agent_decisions
ORDER BY timestamp DESC;
