import { useEffect, useState } from "react";
import { fetchAgents } from "../api/client";
import type { AgentRegistration } from "../api/types";

export function Agents() {
  const [agents, setAgents] = useState<AgentRegistration[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchAgents()
      .then(setAgents)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

  if (error) {
    return <div className="text-red-600 p-4">Error: {error}</div>;
  }

  return (
    <div>
      <h1 className="text-xl font-semibold text-slate-900 mb-4">Registered Agents</h1>

      {loading ? (
        <div className="text-slate-400 p-4">Loading...</div>
      ) : agents.length === 0 ? (
        <div className="bg-white rounded-lg border border-slate-200 p-8 text-center text-slate-400">
          No agents registered yet.
        </div>
      ) : (
        <div className="grid gap-4">
          {agents.map((a) => (
            <div
              key={a.id}
              className="bg-white rounded-lg border border-slate-200 p-4"
            >
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-3">
                  <span
                    className={`w-2 h-2 rounded-full ${
                      a.status === "active" ? "bg-emerald-500" : "bg-slate-300"
                    }`}
                  />
                  <span className="font-medium text-slate-900">{a.name}</span>
                  <span className="text-xs text-slate-400 font-mono">{a.id}</span>
                </div>
                <span
                  className={`text-xs px-2 py-0.5 rounded ${
                    a.status === "active"
                      ? "bg-emerald-100 text-emerald-700"
                      : "bg-slate-100 text-slate-600"
                  }`}
                >
                  {a.status}
                </span>
              </div>
              {a.description && (
                <p className="text-sm text-slate-600 mb-2">{a.description}</p>
              )}
              <div className="flex gap-4 text-xs text-slate-400">
                <span>Registered: {new Date(a.registered_at).toLocaleDateString()}</span>
                <span>Last seen: {new Date(a.last_seen_at).toLocaleString()}</span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
