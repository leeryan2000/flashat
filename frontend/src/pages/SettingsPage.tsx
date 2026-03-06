import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { PATHS } from '../routes/paths';
import { 
  LogOut, 
  Settings, 
  ChevronRight,
  Bell, 
  Lock 
} from 'lucide-react';

export default function SettingsPage() {
  const { logout } = useAuth();
  const navigate = useNavigate();

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
      ]
    },
    {
      title: "Account Actions",
      items: [
        {
          id: "signout",
          label: "Sign Out",
          icon: LogOut,
          onClick: handleLogout,
          isDanger: true, // 👈 This flag tells the UI to make it red
        }
      ]
    }
  ];

  return (
    <div className="h-full w-full bg-slate-900 text-slate-100 flex flex-col items-center p-6 md:p-10 overflow-y-auto">
      
      {/* Container */}
      <div className="w-full max-w-2xl animate-in fade-in zoom-in-95 duration-300">
        
        {/* Header */}
        <div className="flex items-center gap-3 mb-8 pb-6 border-b border-slate-800">
          <div className="p-3 bg-indigo-500/10 rounded-xl text-indigo-400">
            <Settings size={24} />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-white">Settings</h1>
            <p className="text-sm text-slate-400">Manage your application preferences</p>
          </div>
        </div>

        {/* Settings Groups Loop */}
        <div className="space-y-8">
          {settingGroups.map((group, groupIndex) => (
            <div key={groupIndex}>
              {/* Group Title */}
              <h2 className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3 pl-2">
                {group.title}
              </h2>
              
              {/* Group Card */}
              <div className="bg-slate-800/50 border border-slate-700/50 rounded-2xl overflow-hidden">
                {group.items.map((item, itemIndex) => {
                  const Icon = item.icon;
                  const isLast = itemIndex === group.items.length - 1;

                  return (
                    <button
                      key={item.id}
                      onClick={item.onClick}
                      className={`w-full flex items-center justify-between p-4 hover:bg-slate-700/50 transition-colors ${
                        !isLast ? 'border-b border-slate-700/50' : ''
                      }`}
                    >
                      <div className="flex items-center gap-4">
                        <div className={`p-2 rounded-lg ${
                          item.isDanger 
                            ? 'bg-red-500/10 text-red-400' 
                            : 'bg-slate-700 text-slate-300'
                        }`}>
                          <Icon size={18} />
                        </div>
                        <span className={`font-medium ${
                          item.isDanger ? 'text-red-400' : 'text-slate-200'
                        }`}>
                          {item.label}
                        </span>
                      </div>
                      
                      {/* Only show the arrow if it's NOT a danger action (like sign out) */}
                      {!item.isDanger && (
                        <ChevronRight size={18} className="text-slate-500" />
                      )}
                    </button>
                  );
                })}
              </div>
            </div>
          ))}
        </div>

      </div>
    </div>
  );
}