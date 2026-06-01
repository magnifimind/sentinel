import { useEffect, useState } from "react";
import { fetchDecisions } from "../api/client";
import type { AgentDecision } from "../api/types";
import { PolicyBadge } from "../components/PolicyBadge";
import { StatusDot } from "../components/StatusDot";

export function Activity() {
  const [decisions, setDecisions] = useState<AgentDecision[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = () => {
    fetchDecisions(20)
      .then(setDecisions)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    load();
    const interval = setInterval(load, 10_000);
    return () => clearInterval(interval);
  }, []);

  const stats = {
    total: decisions.length,
    approved: decisions.filter((d) => d.policy_result === "approve").length,
    blocked: decisions.filter((d) => d.policy_result === "block").length,
    escalated: decisions.filter((d) => d.policy_result === "escalate").length,
  };

  if (error) {
    return <div className="text-red-600 p-4">Error: {error}</div>;
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-xl font-semibold text-slate-900">Live Activity</h1>
        <span className="text-xs text-slate-400">Auto-refreshes every 10s</span>
      </div>

      <div className="grid grid-cols-4 gap-4 mb-6">
        <StatCard label="Recent Decisions" value={stats.total} />
        <StatCard label="Approved" value={stats.approved} color="text-emerald-600" />
        <StatCard label="Blocked" value={stats.blocked} color="text-red-600" />
        <StatCard label="Escalated" value={stats.escalated} color="text-amber-600" />
      </div>

      <div className="bg-white rounded-lg border border-slate-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="bg-slate-50 border-b border-slate-200">
              <th className="text-left px-4 py-2 font-medium text-slate-600">Time</th>
              <th className="text-left px-4 py-2 font-medium text-slate-600">Agent</th>
              <th className="text-left px-4 py-2 font-medium text-slate-600">Action</th>
              <th className="text-left px-4 py-2 font-medium text-slate-600">Policy</th>
              <th className="text-center px-4 py-2 font-medium text-slate-600">Status</th>
              <th className="text-right px-4 py-2 font-medium text-slate-600">Latency</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr>
                <td colSpan={6} className="px-4 py-8 text-center text-slate-400">
                  Loading...
                </td>
              </tr>
            ) : decisions.length === 0 ? (
              <tr>
                <td colSpan={6} className="px-4 py-8 text-center text-slate-400">
                  No activity yet. Waiting for agent decisions...
                </td>
              </tr>
            ) : (
              decisions.map((d) => (
                <tr key={d.id} className="border-b border-slate-100 hover:bg-slate-50">
                  <td className="px-4 py-2 text-slate-500 font-mono text-xs">
                    {new Date(d.timestamp).toLocaleTimeString()}
                  </td>
                  <td className="px-4 py-2 font-medium">{d.agent_id}</td>
                  <td className="px-4 py-2 font-mono text-xs max-w-xs truncate">{d.action}</td>
                  <td className="px-4 py-2">
                    <PolicyBadge result={d.policy_result} />
                  </td>
                  <td className="px-4 py-2 text-center">
                    <StatusDot success={d.success} />
                  </td>
                  <td className="px-4 py-2 text-right font-mono text-slate-500">
                    {d.latency_ms ? `${d.latency_ms}ms` : "—"}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}

function StatCard({
  label,
  value,
  color = "text-slate-900",
}: {
  label: string;
  value: number;
  color?: string;
}) {
  return (
    <div className="bg-white rounded-lg border border-slate-200 p-4">
      <div className="text-xs text-slate-500 mb-1">{label}</div>
      <div className={`text-2xl font-semibold ${color}`}>{value}</div>
    </div>
  );
}
