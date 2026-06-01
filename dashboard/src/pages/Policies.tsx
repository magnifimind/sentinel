import { useEffect, useState } from "react";
import { fetchPolicyRules } from "../api/client";
import type { PolicyRule } from "../api/types";
import { PolicyBadge } from "../components/PolicyBadge";

const severityColors: Record<string, string> = {
  critical: "bg-red-100 text-red-800",
  high: "bg-orange-100 text-orange-800",
  medium: "bg-amber-100 text-amber-800",
  low: "bg-blue-100 text-blue-800",
};

export function Policies() {
  const [rules, setRules] = useState<PolicyRule[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchPolicyRules()
      .then(setRules)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

  if (error) {
    return <div className="text-red-600 p-4">Error: {error}</div>;
  }

  return (
    <div>
      <h1 className="text-xl font-semibold text-slate-900 mb-4">Policy Rules</h1>

      {loading ? (
        <div className="text-slate-400 p-4">Loading...</div>
      ) : (
        <div className="bg-white rounded-lg border border-slate-200 overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="bg-slate-50 border-b border-slate-200">
                <th className="text-left px-4 py-2 font-medium text-slate-600">Rule ID</th>
                <th className="text-left px-4 py-2 font-medium text-slate-600">Description</th>
                <th className="text-left px-4 py-2 font-medium text-slate-600">Pattern</th>
                <th className="text-left px-4 py-2 font-medium text-slate-600">Decision</th>
                <th className="text-left px-4 py-2 font-medium text-slate-600">Severity</th>
              </tr>
            </thead>
            <tbody>
              {rules.map((r) => (
                <tr key={r.id} className="border-b border-slate-100">
                  <td className="px-4 py-2 font-mono text-xs">{r.id}</td>
                  <td className="px-4 py-2 text-slate-700">{r.description}</td>
                  <td className="px-4 py-2 font-mono text-xs bg-slate-50">{r.action_pattern}</td>
                  <td className="px-4 py-2">
                    <PolicyBadge result={r.decision} />
                  </td>
                  <td className="px-4 py-2">
                    <span
                      className={`inline-block px-2 py-0.5 rounded text-xs font-medium ${
                        severityColors[r.severity] ?? "bg-gray-100 text-gray-700"
                      }`}
                    >
                      {r.severity}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
