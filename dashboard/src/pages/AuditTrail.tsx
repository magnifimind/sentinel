import { useEffect, useState } from "react";
import { fetchAuditTrail } from "../api/client";
import type { AuditTrailEntry } from "../api/types";
import { PolicyBadge } from "../components/PolicyBadge";
import { StatusDot } from "../components/StatusDot";

export function AuditTrail() {
  const [entries, setEntries] = useState<AuditTrailEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(0);
  const pageSize = 25;

  useEffect(() => {
    setLoading(true);
    fetchAuditTrail(pageSize, page * pageSize)
      .then(setEntries)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [page]);

  if (error) {
    return <div className="text-red-600 p-4">Error: {error}</div>;
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-xl font-semibold text-slate-900">Audit Trail</h1>
        <div className="flex gap-2">
          <button
            onClick={() => setPage((p) => Math.max(0, p - 1))}
            disabled={page === 0}
            className="px-3 py-1 text-sm rounded border border-slate-300 disabled:opacity-40"
          >
            Previous
          </button>
          <span className="px-3 py-1 text-sm text-slate-600">Page {page + 1}</span>
          <button
            onClick={() => setPage((p) => p + 1)}
            disabled={entries.length < pageSize}
            className="px-3 py-1 text-sm rounded border border-slate-300 disabled:opacity-40"
          >
            Next
          </button>
        </div>
      </div>

      <div className="bg-white rounded-lg border border-slate-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="bg-slate-50 border-b border-slate-200">
              <th className="text-left px-4 py-2 font-medium text-slate-600">Time</th>
              <th className="text-left px-4 py-2 font-medium text-slate-600">Agent</th>
              <th className="text-left px-4 py-2 font-medium text-slate-600">Type</th>
              <th className="text-left px-4 py-2 font-medium text-slate-600">Action</th>
              <th className="text-left px-4 py-2 font-medium text-slate-600">Policy</th>
              <th className="text-center px-4 py-2 font-medium text-slate-600">Status</th>
              <th className="text-right px-4 py-2 font-medium text-slate-600">Provider</th>
              <th className="text-right px-4 py-2 font-medium text-slate-600">Tokens</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr>
                <td colSpan={8} className="px-4 py-8 text-center text-slate-400">
                  Loading...
                </td>
              </tr>
            ) : entries.length === 0 ? (
              <tr>
                <td colSpan={8} className="px-4 py-8 text-center text-slate-400">
                  No audit entries yet.
                </td>
              </tr>
            ) : (
              entries.map((e) => (
                <tr key={e.id} className="border-b border-slate-100 hover:bg-slate-50">
                  <td className="px-4 py-2 text-slate-500 font-mono text-xs">
                    {new Date(e.timestamp).toLocaleString()}
                  </td>
                  <td className="px-4 py-2 font-medium">{e.agent_id}</td>
                  <td className="px-4 py-2 text-slate-600">{e.decision_type}</td>
                  <td className="px-4 py-2 font-mono text-xs max-w-xs truncate">{e.action}</td>
                  <td className="px-4 py-2">
                    <PolicyBadge result={e.policy_result} />
                  </td>
                  <td className="px-4 py-2 text-center">
                    <StatusDot success={e.success} />
                  </td>
                  <td className="px-4 py-2 text-right text-slate-500">{e.llm_provider || "—"}</td>
                  <td className="px-4 py-2 text-right font-mono text-slate-500">
                    {e.tokens_used || "—"}
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
