import type {
  AgentDecision,
  AgentRegistration,
  AuditTrailEntry,
  PolicyRule,
} from "./types";

const BASE = "/api/v1";

async function get<T>(path: string, params?: Record<string, string>): Promise<T> {
  const url = new URL(path, window.location.origin);
  if (params) {
    Object.entries(params).forEach(([k, v]) => url.searchParams.set(k, v));
  }
  const res = await fetch(url.toString());
  if (!res.ok) {
    throw new Error(`API error: ${res.status} ${res.statusText}`);
  }
  return res.json();
}

export async function fetchDecisions(
  limit = 50,
  offset = 0
): Promise<AgentDecision[]> {
  return get<AgentDecision[]>(`${BASE}/audit/decisions`, {
    limit: String(limit),
    offset: String(offset),
  });
}

export async function fetchAgentDecisions(
  agentId: string,
  limit = 50
): Promise<AgentDecision[]> {
  return get<AgentDecision[]>(`${BASE}/audit/decisions/${agentId}`, {
    limit: String(limit),
  });
}

export async function fetchAuditTrail(
  limit = 50,
  offset = 0
): Promise<AuditTrailEntry[]> {
  return get<AuditTrailEntry[]>(`${BASE}/audit/trail`, {
    limit: String(limit),
    offset: String(offset),
  });
}

export async function fetchAgents(): Promise<AgentRegistration[]> {
  return get<AgentRegistration[]>(`${BASE}/agents`);
}

export async function fetchPolicyRules(): Promise<PolicyRule[]> {
  return get<PolicyRule[]>(`${BASE}/policy/rules`);
}
