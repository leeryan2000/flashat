import { Outlet, Link } from "react-router-dom";
import Footer from "./Footer";
import Header from "./Header";
import Sidebar from "./Sidebar";

export default function Layout() {
  return (
    <div className="min-h-screen flex flex-col bg-slate-100">
      {/* Header */}
      <header className="bg-indigo-600 text-white p-4 shadow-lg">
          <Header />
      </header>

      {/* Body: sidebar + main */}
      <div className="flex flex-1">
        {/* Sidebar */}
        <aside className="w-64 bg-slate-900 text-slate-100 p-4 space-y-3 shadow-xl">
          <Sidebar/>
        </aside>

        {/* Main content → child pages render here */}
        <main className="flex-1 p-8 bg-white">
            <div className="rounded-xl p-4 bg-amber-50 ring-2 ring-amber-300">
              <Outlet />
            </div>
        </main>
      </div>
      <footer className="bg-indigo-50 text-indigo-900 text-center text-sm p-4 border-t border-indigo-200">
        <Footer />
      </footer>
    </div>
  );
}
