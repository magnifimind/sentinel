import { NavLink, Outlet } from "react-router-dom";

const links = [
  { to: "/", label: "Activity" },
  { to: "/audit", label: "Audit Trail" },
  { to: "/agents", label: "Agents" },
  { to: "/policies", label: "Policies" },
];

export function Layout() {
  return (
    <div className="min-h-screen bg-slate-50">
      <nav className="bg-slate-900 text-white">
        <div className="max-w-7xl mx-auto px-4 flex items-center h-14 gap-8">
          <span className="text-lg font-semibold tracking-tight">Sentinel</span>
          <div className="flex gap-1">
            {links.map((l) => (
              <NavLink
                key={l.to}
                to={l.to}
                className={({ isActive }) =>
                  `px-3 py-1.5 rounded text-sm transition-colors ${
                    isActive
                      ? "bg-slate-700 text-white"
                      : "text-slate-300 hover:text-white hover:bg-slate-800"
                  }`
                }
              >
                {l.label}
              </NavLink>
            ))}
          </div>
        </div>
      </nav>
      <main className="max-w-7xl mx-auto px-4 py-6">
        <Outlet />
      </main>
    </div>
  );
}
