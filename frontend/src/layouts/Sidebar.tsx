import { Link } from "react-router-dom";

export default function Sidebar() {
  return (
    <div>
      {/* <p className="text-xs uppercase tracking-widest text-slate-400">Nav</p> */}
      <nav className="space-y-2">
        <Link
          to=""
          className="block rounded-lg px-3 py-2 bg-slate-800 hover:bg-slate-700 active:scale-[0.98] transition"
        >
          Chat
        </Link>
        <Link
          to="friends"
          className="block rounded-lg px-3 py-2 bg-slate-800 hover:bg-slate-700 active:scale-[0.98] transition"
        >
          Friends
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
    </div>
  );
}
