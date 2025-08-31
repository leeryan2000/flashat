import { Outlet, Link } from "react-router-dom";

export default function Layout() {
  return (
    <div className="min-h-screen flex flex-col bg-slate-100">
      {/* Header */}
      <header className="bg-indigo-600 text-white p-4 shadow-lg">
        <h1 className="text-2xl font-extrabold tracking-wide">
          🚀 Tailwind Smoke Test
        </h1>
        <p className="text-indigo-100 text-sm">
          If you can see this bright header, Tailwind is working.
        </p>
      </header>

      {/* Body: sidebar + main */}
      <div className="flex flex-1">
        {/* Sidebar */}
        <aside className="w-64 bg-slate-900 text-slate-100 p-4 space-y-3 shadow-xl">
          <p className="text-xs uppercase tracking-widest text-slate-400">
            Nav
          </p>
          <nav className="space-y-2">
            <Link
              to=""
              className="block rounded-lg px-3 py-2 bg-slate-800 hover:bg-slate-700 active:scale-[0.98] transition"
            >
              Dashboard
            </Link>
            <Link
              to="profile"
              className="block rounded-lg px-3 py-2 bg-slate-800 hover:bg-slate-700 active:scale-[0.98] transition"
            >
              Profile
            </Link>
            <Link
              to="settings"
              className="block rounded-lg px-3 py-2 bg-slate-800 hover:bg-slate-700 active:scale-[0.98] transition"
            >
              Settings
            </Link>
          </nav>

          {/* obvious interactive test */}
          <button className="mt-4 w-full rounded-xl bg-emerald-500 text-white py-2 font-semibold shadow-lg hover:shadow-emerald-300/50 hover:brightness-110 active:scale-95 transition">
            Click me
          </button>
        </aside>

        {/* Main content → child pages render here */}
        <main className="flex-1 p-8 bg-white">
          <div className="rounded-2xl border-4 border-dashed border-indigo-300 p-6 shadow-2xl">
            <h2 className="text-xl font-bold mb-2">Main content area</h2>
            <p className="text-slate-600 mb-4">
              You should see dashed borders, heavy shadows, and spacing.
            </p>

            {/* Where child routes render */}
            <div className="rounded-xl p-4 bg-amber-50 ring-2 ring-amber-300">
              <Outlet />
            </div>
          </div>
        </main>
      </div>

      {/* Footer */}
      <footer className="bg-indigo-50 text-indigo-900 text-center text-sm p-4 border-t border-indigo-200">
        Tailwind footer • <span className="font-semibold">very visible</span>
      </footer>
    </div>
  );
}
