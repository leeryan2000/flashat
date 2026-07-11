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
    <div className="w-full flex flex-col items-center flex-1">
      {/* Brand */}
      <div
        className="h-9 w-9 rounded-xl flex items-center justify-center rotate-3 border shrink-0 mb-6"
        style={{
          background: "color-mix(in srgb, var(--primary) 15%, transparent)",
          borderColor: "color-mix(in srgb, var(--primary) 30%, transparent)",
        }}
        title="Flashat"
      >
        <MessageSquare size={18} style={{ color: "var(--primary)" }} />
      </div>

      <nav className="flex flex-col items-center gap-2 flex-1">
        {navItems.map(({ to, end, icon: Icon, label }) => (
          <motion.div key={to} whileHover={{ y: -2 }} whileTap={{ scale: 0.95 }}>
            <NavLink to={to} end={end}>
              {({ isActive }) => (
                <div
                  className="flex flex-col items-center gap-1 px-3 py-2.5 rounded-xl font-medium transition-colors duration-150 cursor-pointer"
                  style={
                    isActive
                      ? { background: "var(--primary)", color: "white", boxShadow: "0 4px 12px var(--primary-shadow)" }
                      : { color: "var(--text-soft)" }
                  }
                  onMouseEnter={e => { if (!isActive) { (e.currentTarget as HTMLDivElement).style.background = "var(--sidebar-item)"; (e.currentTarget as HTMLDivElement).style.color = "var(--text)"; } }}
                  onMouseLeave={e => { if (!isActive) { (e.currentTarget as HTMLDivElement).style.background = ""; (e.currentTarget as HTMLDivElement).style.color = "var(--text-soft)"; } }}
                >
                  <Icon size={20} />
                  <span className="text-[10px] font-medium leading-none">{label}</span>
                </div>
              )}
            </NavLink>
          </motion.div>
        ))}
      </nav>

      {user && (
        <motion.div
          whileHover={{ scale: 1.08 }}
          title={`${user.name} · ${user.email}`}
          className="cursor-default mt-4"
        >
          <div className={`w-9 h-9 rounded-full flex items-center justify-center text-white text-xs font-semibold shrink-0 ${avatarColor(user.name)}`}>
            {nameInitials(user.name)}
          </div>
        </motion.div>
      )}
    </div>
  );
}
