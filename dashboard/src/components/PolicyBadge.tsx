const colors: Record<string, string> = {
  approve: "bg-emerald-100 text-emerald-800",
  block: "bg-red-100 text-red-800",
  escalate: "bg-amber-100 text-amber-800",
  recorded: "bg-slate-100 text-slate-700",
};

export function PolicyBadge({ result }: { result: string }) {
  const cls = colors[result] ?? "bg-gray-100 text-gray-700";
  return (
    <span className={`inline-block px-2 py-0.5 rounded text-xs font-medium ${cls}`}>
      {result}
    </span>
  );
}
