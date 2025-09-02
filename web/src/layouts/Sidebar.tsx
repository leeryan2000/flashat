import { Link } from "react-router-dom";

export default function Sidebar() {
  return (
    <div>
      <p className="text-xs uppercase tracking-widest text-slate-400">Nav</p>
      <nav className="space-y-2">
        <Link
          to=""
          className="block rounded-lg px-3 py-2 bg-slate-800 hover:bg-slate-700 active:scale-[0.98] transition"
        >
          Chat
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
    </div>
  );
}
