import { NavLink } from "react-router-dom";
import { MessageSquare, Users, User, Settings } from "lucide-react";
import { useAuth } from "../context/AuthContext";
import { avatarColor, nameInitials } from "../utils/avatar";
import { motion } from "framer-motion";

const navItems = [
  { to: "", end: true, icon: MessageSquare, label: "Chat" },
  { to: "friends", icon: Users, label: "Friends" },
  { to: "profile", icon: User, label: "Profile" },
  { to: "settings", icon: Settings, label: "Settings" },
];

export default function Sidebar() {
  const { user } = useAuth();

  return (
    <div className="w-full flex flex-col flex-1">
      <nav className="space-y-2 flex-1">
        {navItems.map(({ to, end, icon: Icon, label }) => (
          <motion.div key={to} whileHover={{ x: 4 }} whileTap={{ scale: 0.97 }}>
            <NavLink to={to} end={end}>
              {({ isActive }) => (
                <div
                  className="flex items-center gap-3 w-full rounded-xl px-4 py-3 font-medium transition-colors duration-150 cursor-pointer"
                  style={
                    isActive
                      ? { background: "var(--primary)", color: "white", boxShadow: "0 4px 12px var(--primary-shadow)" }
                      : { color: "#cbd5e1" }
                  }
                  onMouseEnter={e => { if (!isActive) (e.currentTarget as HTMLDivElement).style.background = "var(--sidebar-item)"; (e.currentTarget as HTMLDivElement).style.color = "#f8fafc"; }}
                  onMouseLeave={e => { if (!isActive) { (e.currentTarget as HTMLDivElement).style.background = ""; (e.currentTarget as HTMLDivElement).style.color = "#cbd5e1"; } }}
                >
                  <Icon size={20} />
                  <span>{label}</span>
                </div>
              )}
            </NavLink>
          </motion.div>
        ))}
      </nav>

      {user && (
        <motion.div
          whileHover={{ scale: 1.02 }}
          className="flex items-center gap-3 px-2 py-3 rounded-xl mt-4 cursor-default"
          style={{ background: "var(--sidebar-item)" }}
        >
          <div className={`w-9 h-9 rounded-full flex items-center justify-center text-white text-xs font-semibold shrink-0 ${avatarColor(user.uid)}`}>
            {nameInitials(user.name)}
          </div>
          <div className="flex flex-col min-w-0">
            <span className="text-sm font-medium text-slate-100 truncate">{user.name}</span>
            <span className="text-xs text-slate-400 truncate">{user.email}</span>
          </div>
        </motion.div>
      )}
    </div>
  );
}
