export interface AgentDecision {
  id: number;
  timestamp: string;
  agent_id: string;
  decision_type: string;
  context: Record<string, unknown>;
  action: string;
  policy_result: string;
  success: boolean;
  llm_provider?: string;
  tokens_used?: number;
  latency_ms?: number;
  model?: string;
}

export interface AuditTrailEntry {
  id: number;
  timestamp: string;
  agent_id: string;
  decision_type: string;
  action: string;
  description: string;
  policy_result: string;
  success: boolean;
  llm_provider?: string;
  tokens_used?: number;
}

export interface AgentRegistration {
  id: string;
  name: string;
  description?: string;
  capabilities?: Record<string, unknown>;
  registered_at: string;
  last_seen_at: string;
  status: string;
}

export interface PolicyRule {
  id: string;
  description: string;
  action_pattern: string;
  decision: string;
  severity: string;
}

export interface PolicyEvalResponse {
  decision: string;
  reason: string;
  policy_id?: string;
}
