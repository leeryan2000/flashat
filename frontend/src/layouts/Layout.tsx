import { Outlet } from "react-router-dom";
import Sidebar from "./Sidebar";

export default function Layout() {
  return (
    <div className="h-screen w-screen flex overflow-hidden" style={{ background: "var(--sidebar-bg)" }}>
      {/* Sidebar */}
      <aside
        className="w-[84px] py-4 px-2 flex flex-col h-full shrink-0 border-r"
        style={{ background: "var(--nav-bg)", borderColor: "var(--panel-border)", color: "var(--text)" }}
      >
        <Sidebar />
      </aside>

      {/* Main content → child pages render here */}
      <main className="h-full flex-1 overflow-hidden">
        <div className="h-full w-full flex flex-col">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
