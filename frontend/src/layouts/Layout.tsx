import { Outlet } from "react-router-dom";
import Footer from "./Footer";
import Header from "./Header";
import Sidebar from "./Sidebar";

export default function Layout() {
  return (
    <div className="h-screen w-screen flex flex-col bg-slate-100 overflow-hidden">
      {/* Header */}
      <header className="text-white p-4 shadow-lg" style={{ background: "var(--primary)" }}>
          <Header />
      </header>

      {/* Body: sidebar + main */}
      <div className="flex flex-1 overflow-hidden">
        {/* Sidebar */}
        <aside className="w-64 text-slate-100 p-4 shadow-xl flex flex-col h-full" style={{ background: "var(--sidebar-bg)" }}>
          <Sidebar/>
        </aside>

        {/* Main content → child pages render here */}
        <main className="h-full flex-1 p-4 bg-white overflow-hidden">
            <div className="h-full w-full rounded-xl bg-white ring-2 ring-slate-200 overflow-hidden">
              <div className="h-full w-full flex flex-col">
                <Outlet />
              </div>
            </div>
        </main>
      </div>
      <footer className="bg-slate-50 text-slate-600 text-center text-sm border-t border-slate-200">
        <Footer />
      </footer>
    </div>
  );
}
