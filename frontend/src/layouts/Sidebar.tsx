import { NavLink } from "react-router-dom";
import { MessageSquare, Users, User, Settings } from "lucide-react";

export default function Sidebar() {
  
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
    <div className="w-full">
      {/* Optional: A small section label */}
      <p className="text-xs font-semibold uppercase tracking-widest text-slate-500 mb-4 px-4">
        Main Menu
      </p>
      
      <nav className="space-y-2">
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
    </div>
  );
}