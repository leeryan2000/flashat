import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useTheme } from '../context/ThemeContext';
import { PATHS } from '../routes/paths';
import { LogOut, Settings, ChevronRight, Bell, Lock, Palette, X, Sun, Moon } from 'lucide-react';

export default function SettingsPage() {
  const { logout } = useAuth();
  const { theme, setTheme, themes, mode, setMode } = useTheme();
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
        { id: "theme", label: "Appearance", icon: Palette, onClick: () => setThemeOpen(true), isDanger: false },
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
    <div className="h-full w-full flex flex-col items-center p-6 md:p-10 overflow-y-auto" style={{ background: "var(--sidebar-bg)", color: "var(--text)" }}>
      <div className="w-full max-w-2xl animate-in fade-in zoom-in-95 duration-300">

        {/* Header */}
        <div className="flex items-center gap-3 mb-8 pb-6 border-b" style={{ borderColor: "var(--panel-border)" }}>
          <div className="p-3 rounded-xl" style={{ color: "var(--primary)", background: "color-mix(in srgb, var(--primary) 15%, transparent)" }}>
            <Settings size={24} />
          </div>
          <div>
            <h1 className="text-2xl font-bold" style={{ color: "var(--text)" }}>Settings</h1>
            <p className="text-sm" style={{ color: "var(--text-soft)" }}>Manage your application preferences</p>
          </div>
        </div>

        {/* Settings Groups */}
        <div className="space-y-8">
          {settingGroups.map((group, groupIndex) => (
            <div key={groupIndex}>
              <h2 className="text-xs font-semibold uppercase tracking-wider mb-3 pl-2" style={{ color: "var(--text-faint)" }}>
                {group.title}
              </h2>
              <div className="border rounded-2xl overflow-hidden" style={{ background: "var(--sidebar-item)", borderColor: "var(--panel-border)" }}>
                {group.items.map((item, itemIndex) => {
                  const Icon = item.icon;
                  const isLast = itemIndex === group.items.length - 1;
                  return (
                    <button
                      key={item.id}
                      onClick={item.onClick}
                      className="w-full flex items-center justify-between p-4 transition-colors"
                      style={!isLast ? { borderBottom: "1px solid var(--panel-border)" } : undefined}
                      onMouseEnter={e => (e.currentTarget.style.background = "color-mix(in srgb, var(--primary) 8%, transparent)")}
                      onMouseLeave={e => (e.currentTarget.style.background = "")}
                    >
                      <div className="flex items-center gap-4">
                        <div
                          className="p-2 rounded-lg"
                          style={
                            item.isDanger
                              ? { background: "var(--danger-bg)", color: "var(--danger-text)" }
                              : { background: "color-mix(in srgb, var(--primary) 15%, transparent)", color: "var(--primary)" }
                          }
                        >
                          <Icon size={18} />
                        </div>
                        <span className="font-medium" style={{ color: item.isDanger ? "var(--danger-text)" : "var(--text)" }}>
                          {item.label}
                        </span>
                      </div>
                      <div className="flex items-center gap-2">
                        {item.id === "theme" && (
                          <span className="text-xs capitalize" style={{ color: "var(--text-soft)" }}>{theme} · {mode}</span>
                        )}
                        {!item.isDanger && <ChevronRight size={18} style={{ color: "var(--text-faint)" }} />}
                      </div>
                    </button>
                  );
                })}
              </div>
            </div>
          ))}
        </div>

      </div>

      {/* Appearance modal */}
      {themeOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center" style={{ background: "var(--overlay)" }} onClick={() => setThemeOpen(false)}>
          <div
            className="border rounded-2xl p-6 w-80 shadow-xl"
            style={{ background: "var(--sidebar-item)", borderColor: "var(--panel-border)" }}
            onClick={e => e.stopPropagation()}
          >
            <div className="flex items-center justify-between mb-5">
              <h3 className="font-semibold" style={{ color: "var(--text)" }}>Appearance</h3>
              <button onClick={() => setThemeOpen(false)} className="transition" style={{ color: "var(--text-soft)" }}>
                <X size={18} />
              </button>
            </div>

            {/* Light / Dark mode toggle */}
            <div
              className="grid grid-cols-2 gap-1 p-1 rounded-xl mb-5"
              style={{ background: "var(--chat-bg)", border: "1px solid var(--panel-border)" }}
            >
              {(["dark", "light"] as const).map(m => (
                <button
                  key={m}
                  onClick={() => setMode(m)}
                  className="flex items-center justify-center gap-2 px-3 py-2 rounded-lg text-sm font-medium capitalize transition-all duration-150"
                  style={
                    mode === m
                      ? { background: "var(--primary)", color: "white" }
                      : { color: "var(--text-soft)" }
                  }
                >
                  {m === "dark" ? <Moon size={15} /> : <Sun size={15} />}
                  {m}
                </button>
              ))}
            </div>

            <p className="text-xs font-semibold uppercase tracking-wider mb-2" style={{ color: "var(--text-faint)" }}>
              Accent
            </p>
            <div className="flex flex-col gap-2">
              {themes.map(t => (
                <button
                  key={t.id}
                  onClick={() => { setTheme(t.id); }}
                  className="flex items-center gap-3 px-4 py-3 rounded-xl border transition-all duration-150 text-sm font-medium"
                  style={{
                    borderColor: theme === t.id ? t.color : "transparent",
                    background: theme === t.id ? `color-mix(in srgb, ${t.color} 15%, transparent)` : "var(--chat-bg)",
                    color: theme === t.id ? t.color : "var(--text-soft)",
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
