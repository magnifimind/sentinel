export function StatusDot({ success }: { success: boolean }) {
  return (
    <span
      className={`inline-block w-2 h-2 rounded-full ${success ? "bg-emerald-500" : "bg-red-500"}`}
      title={success ? "Success" : "Failed"}
    />
  );
}
