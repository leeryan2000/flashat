import { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { 
  User, 
  Mail, 
  Camera, 
  Check,
  X,
  Edit2
} from 'lucide-react';

export default function ProfilePage() {
  // Pull current user data and logout function from your context
  const { user, updateProfile } = useAuth();

  // State for toggling edit mode
  const [isEditing, setIsEditing] = useState(false);
  const [editName, setEditName] = useState(user?.name || "");
  const [isSaving, setIsSaving] = useState(false);

  const handleSave = async () => {
    if (!editName.trim() || editName === user?.name) {
      setIsEditing(false);
      return;
    }

    setIsSaving(true);
    try {
      // Assuming your context has an update function wired to your Go backend
      if (updateProfile) {
        await updateProfile({ name: editName });
      }
      setIsEditing(false);
    } catch (error) {
      console.error("Failed to update profile", error);
      // You could add an error state here similar to your login form
    } finally {
      setIsSaving(false);
    }
  };

  // Helper to get avatar initials (e.g., "Ryan" -> "RY")
  const getInitials = (name?: string) => {
    return name ? name.slice(0, 2).toUpperCase() : "??";
  };

  return (
    <div className="h-full w-full bg-slate-900 text-slate-100 flex flex-col items-center p-6 md:p-10 overflow-y-auto">
      
      {/* Profile Card Container */}
      <div className="w-full max-w-2xl bg-slate-800 border border-slate-700 rounded-3xl p-8 shadow-2xl animate-in fade-in zoom-in-95 duration-300">
        
        {/* Header & Avatar Section */}
        <div className="flex flex-col items-center mb-10">
          <div className="relative group cursor-pointer mb-4">
            {/* Avatar Circle */}
            <div className="h-28 w-28 rounded-full bg-indigo-600 flex items-center justify-center text-4xl font-bold shadow-lg shadow-indigo-500/30">
              {getInitials(user?.name)}
            </div>
            
            {/* Hover Overlay for changing avatar (Optional UI) */}
            <div className="absolute inset-0 bg-slate-900/60 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-200">
              <Camera className="w-8 h-8 text-white" />
            </div>
          </div>
          
          <h2 className="text-2xl font-bold text-white">Profile</h2>
          <p className="text-sm text-slate-400">Manage your Profile</p>
        </div>

        {/* User Details Section */}
        <div className="space-y-6">
          
          {/* Email Field (Read-only) */}
          <div className="bg-slate-900/50 rounded-2xl p-4 border border-slate-700/50 flex items-center gap-4">
            <div className="p-3 bg-slate-800 rounded-xl text-slate-400">
              <Mail size={20} />
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-xs text-slate-500 uppercase tracking-wider font-semibold mb-1">Email Address</p>
              <p className="text-slate-200 font-medium truncate">{user?.email || "user@example.com"}</p>
            </div>
          </div>

          {/* Name Field (Editable) */}
          <div className="bg-slate-900/50 rounded-2xl p-4 border border-slate-700/50 flex items-center gap-4 transition-all">
            <div className="p-3 bg-slate-800 rounded-xl text-indigo-400">
              <User size={20} />
            </div>
            
            <div className="flex-1 min-w-0">
              <p className="text-xs text-slate-500 uppercase tracking-wider font-semibold mb-1">Display Name</p>
              
              {isEditing ? (
                <input
                  type="text"
                  autoFocus
                  value={editName}
                  onChange={(e) => setEditName(e.target.value)}
                  className="w-full bg-slate-800 border border-indigo-500/50 rounded-lg px-3 py-1 text-white outline-none focus:ring-2 focus:ring-indigo-500"
                />
              ) : (
                <p className="text-slate-200 font-medium truncate">{user?.name || "Username"}</p>
              )}
            </div>
            {/* Edit / Save Actions */}
            <div className="flex gap-2">
              {isEditing ? (
                <>
                  <button 
                    onClick={handleSave}
                    disabled={isSaving}
                    className="p-2 rounded-lg bg-emerald-500/20 text-emerald-400 hover:bg-emerald-500/30 transition disabled:opacity-50"
                  >
                    <Check size={18} />
                  </button>
                  <button 
                    onClick={() => {
                      setIsEditing(false);
                      setEditName(user?.name || ""); // Reset
                    }}
                    className="p-2 rounded-lg bg-slate-700 text-slate-400 hover:text-white transition"
                  >
                    <X size={18} />
                  </button>
                </>
              ) : (
                <button 
                  onClick={() => setIsEditing(true)}
                  className="p-2 rounded-lg bg-slate-800 border border-slate-700 text-slate-400 hover:text-indigo-400 hover:border-indigo-500/50 transition"
                  title="Edit Name"
                >
                  <Edit2 size={18} />
                </button>
              )}
            </div>
          </div>
        </div>

        {/* Divider */}
        <hr className="my-8 border-slate-700" />
      </div>
    </div>
  );
}