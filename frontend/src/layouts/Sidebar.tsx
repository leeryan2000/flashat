import { NavLink } from "react-router-dom";
import { MessageSquare, Users, User, Settings } from "lucide-react";
import { useAuth } from "../context/AuthContext";
import { avatarColor, nameInitials } from "../utils/avatar";

export default function Sidebar() {
  const { user } = useAuth();

  // 🛠️ Helper function to style the links dynamically based on active state
  const getNavClass = ({ isActive }: { isActive: boolean }) => {
    const baseClasses = "flex items-center gap-3 w-full rounded-xl px-4 py-3 transition-all active:scale-[0.98] font-medium";
    
    if (isActive) {
      // Active State: Highlighted with your theme's indigo color
      return `${baseClasses} bg-indigo-600 text-white shadow-lg shadow-indigo-500/20`;
    } else {
      // Inactive State: Dimmed, turns lighter on hover
      return `${baseClasses} text-slate-400 hover:bg-slate-800 hover:text-slate-100`;
    }
  };

  return (
    <div className="w-full flex flex-col flex-1">
      {/* Optional: A small section label */}
      <p className="text-xs font-semibold uppercase tracking-widest text-slate-500 mb-4 px-4">
        Main Menu
      </p>

      <nav className="space-y-2 flex-1">
        {/* Note: Added 'end' so the Chat button doesn't stay highlighted when you visit sub-pages */}
        <NavLink to="" end className={getNavClass}>
          <MessageSquare size={20} />
          <span>Chat</span>
        </NavLink>

        <NavLink to="friends" className={getNavClass}>
          <Users size={20} />
          <span>Friends</span>
        </NavLink>

        <NavLink to="profile" className={getNavClass}>
          <User size={20} />
          <span>Profile</span>
        </NavLink>

        <NavLink to="settings" className={getNavClass}>
          <Settings size={20} />
          <span>Settings</span>
        </NavLink>
      </nav>

      {/* User card at the bottom */}
      {user && (
        <div className="flex items-center gap-3 px-2 py-3 rounded-xl bg-slate-800 mt-4">
          <div className={`w-9 h-9 rounded-full flex items-center justify-center text-white text-xs font-semibold shrink-0 ${avatarColor(user.uid)}`}>
            {nameInitials(user.name)}
          </div>
          <div className="flex flex-col min-w-0">
            <span className="text-sm font-medium text-slate-100 truncate">{user.name}</span>
            <span className="text-xs text-slate-400 truncate">{user.email}</span>
          </div>
        </div>
      )}
    </div>
  );
}