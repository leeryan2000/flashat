import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useTheme } from '../context/ThemeContext';
import { PATHS } from '../routes/paths';
import { LogOut, Settings, ChevronRight, Bell, Lock, Palette, X } from 'lucide-react';

export default function SettingsPage() {
  const { logout } = useAuth();
  const { theme, setTheme, themes } = useTheme();
  const navigate = useNavigate();
  const [themeOpen, setThemeOpen] = useState(false);

  const handleLogout = () => {
    if (logout) logout();
    navigate(PATHS.login);
  };

  const settingGroups = [
    {
      title: "General",
      items: [
        { id: "notifications", label: "Notifications", icon: Bell, onClick: () => console.log("Notify"), isDanger: false },
        { id: "privacy", label: "Privacy & Security", icon: Lock, onClick: () => console.log("Privacy"), isDanger: false },
        { id: "theme", label: "Theme", icon: Palette, onClick: () => setThemeOpen(true), isDanger: false },
      ]
    },
    {
      title: "Account Actions",
      items: [
        { id: "logout", label: "Log Out", icon: LogOut, onClick: handleLogout, isDanger: true }
      ]
    }
  ];

  return (
    <div className="h-full w-full text-slate-100 flex flex-col items-center p-6 md:p-10 overflow-y-auto" style={{ background: "var(--sidebar-bg)" }}>
      <div className="w-full max-w-2xl animate-in fade-in zoom-in-95 duration-300">

        {/* Header */}
        <div className="flex items-center gap-3 mb-8 pb-6 border-b border-slate-800">
          <div className="p-3 rounded-xl" style={{ color: "var(--primary)", background: "color-mix(in srgb, var(--primary) 15%, transparent)" }}>
            <Settings size={24} />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-white">Settings</h1>
            <p className="text-sm text-slate-400">Manage your application preferences</p>
          </div>
        </div>

        {/* Settings Groups */}
        <div className="space-y-8">
          {settingGroups.map((group, groupIndex) => (
            <div key={groupIndex}>
              <h2 className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3 pl-2">
                {group.title}
              </h2>
              <div className="border border-slate-700/20 rounded-2xl overflow-hidden" style={{ background: "var(--sidebar-item)" }}>
                {group.items.map((item, itemIndex) => {
                  const Icon = item.icon;
                  const isLast = itemIndex === group.items.length - 1;
                  return (
                    <button
                      key={item.id}
                      onClick={item.onClick}
                      className={`w-full flex items-center justify-between p-4 transition-colors ${!isLast ? 'border-b border-slate-700/20' : ''}`}
                      onMouseEnter={e => (e.currentTarget.style.background = "color-mix(in srgb, var(--primary) 8%, transparent)")}
                      onMouseLeave={e => (e.currentTarget.style.background = "")}
                    >
                      <div className="flex items-center gap-4">
                        <div
                          className={`p-2 rounded-lg ${item.isDanger ? 'bg-red-500/10 text-red-400' : ''}`}
                          style={!item.isDanger ? { background: "color-mix(in srgb, var(--primary) 15%, transparent)", color: "var(--primary)" } : undefined}
                        >
                          <Icon size={18} />
                        </div>
                        <span className={`font-medium ${item.isDanger ? 'text-red-400' : 'text-slate-200'}`}>
                          {item.label}
                        </span>
                      </div>
                      <div className="flex items-center gap-2">
                        {item.id === "theme" && (
                          <span className="text-xs text-slate-400 capitalize">{theme}</span>
                        )}
                        {!item.isDanger && <ChevronRight size={18} className="text-slate-500" />}
                      </div>
                    </button>
                  );
                })}
              </div>
            </div>
          ))}
        </div>

      </div>

      {/* Theme modal */}
      {themeOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60" onClick={() => setThemeOpen(false)}>
          <div className="border border-slate-700/20 rounded-2xl p-6 w-80 shadow-xl" style={{ background: "var(--sidebar-item)" }} onClick={e => e.stopPropagation()}>
            <div className="flex items-center justify-between mb-5">
              <h3 className="font-semibold text-white">Choose Theme</h3>
              <button onClick={() => setThemeOpen(false)} className="text-slate-400 hover:text-slate-200 transition">
                <X size={18} />
              </button>
            </div>
            <div className="flex flex-col gap-2">
              {themes.map(t => (
                <button
                  key={t.id}
                  onClick={() => { setTheme(t.id); setThemeOpen(false); }}
                  className="flex items-center gap-3 px-4 py-3 rounded-xl border transition-all duration-150 text-sm font-medium"
                  style={{
                    borderColor: theme === t.id ? t.color : "transparent",
                    background: theme === t.id ? `color-mix(in srgb, ${t.color} 15%, transparent)` : "rgba(255,255,255,0.05)",
                    color: theme === t.id ? t.color : "#94a3b8",
                  }}
                >
                  <span className="w-4 h-4 rounded-full shrink-0" style={{ background: t.color }} />
                  {t.label}
                  {theme === t.id && <span className="ml-auto text-xs opacity-70">Active</span>}
                </button>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
