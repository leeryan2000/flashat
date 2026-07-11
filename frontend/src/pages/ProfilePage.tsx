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
import { avatarColor, nameInitials } from '../utils/avatar';

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


  return (
    <div className="h-full w-full flex flex-col items-center p-6 md:p-10 overflow-y-auto" style={{ background: "var(--sidebar-bg)", color: "var(--text)" }}>

      {/* Profile Card Container */}
      <div className="w-full max-w-2xl border rounded-3xl p-8 shadow-2xl animate-in fade-in zoom-in-95 duration-300" style={{ background: "var(--sidebar-item)", borderColor: "var(--panel-border)" }}>

        {/* Header & Avatar Section */}
        <div className="flex flex-col items-center mb-10">
          <div className="relative group cursor-pointer mb-4">
            {/* Avatar Circle */}
            <div className={`h-28 w-28 rounded-full flex items-center justify-center text-4xl font-bold text-white shadow-lg ${user?.name ? avatarColor(user.name) : ""}`} style={{ boxShadow: "0 8px 24px var(--primary-shadow)" }}>
              {user?.name ? nameInitials(user.name) : "??"}
            </div>

            {/* Hover Overlay for changing avatar (Optional UI) */}
            <div className="absolute inset-0 bg-black/50 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-200">
              <Camera className="w-8 h-8 text-white" />
            </div>
          </div>

          <h2 className="text-2xl font-bold" style={{ color: "var(--text)" }}>Profile</h2>
          <p className="text-sm" style={{ color: "var(--text-soft)" }}>Manage your Profile</p>
        </div>

        {/* User Details Section */}
        <div className="space-y-6">

          {/* Email Field (Read-only) */}
          <div className="rounded-2xl p-4 border flex items-center gap-4" style={{ background: "color-mix(in srgb, var(--primary) 6%, transparent)", borderColor: "var(--panel-border)" }}>
            <div className="p-3 rounded-xl" style={{ background: "var(--chat-bg)", color: "var(--primary)" }}>
              <Mail size={20} />
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-xs uppercase tracking-wider font-semibold mb-1" style={{ color: "var(--text-faint)" }}>Email Address</p>
              <p className="font-medium truncate" style={{ color: "var(--text)" }}>{user?.email || "user@example.com"}</p>
            </div>
          </div>

          {/* Name Field (Editable) */}
          <div className="rounded-2xl p-4 border flex items-center gap-4 transition-all" style={{ background: "color-mix(in srgb, var(--primary) 6%, transparent)", borderColor: "var(--panel-border)" }}>
            <div className="p-3 rounded-xl" style={{ background: "var(--chat-bg)", color: "var(--primary)" }}>
              <User size={20} />
            </div>

            <div className="flex-1 min-w-0">
              <p className="text-xs uppercase tracking-wider font-semibold mb-1" style={{ color: "var(--text-faint)" }}>Display Name</p>
              {isEditing ? (
                <input
                  type="text"
                  autoFocus
                  value={editName}
                  onChange={(e) => setEditName(e.target.value)}
                  className="w-full rounded-lg px-3 py-1 outline-none focus:ring-2 border"
                  style={{ background: "var(--chat-bg)", color: "var(--text)", borderColor: "color-mix(in srgb, var(--primary) 50%, transparent)", "--tw-ring-color": "var(--primary)" } as React.CSSProperties}
                />
              ) : (
                <p className="font-medium truncate" style={{ color: "var(--text)" }}>{user?.name || "Username"}</p>
              )}
            </div>

            <div className="flex gap-2">
              {isEditing ? (
                <>
                  <button
                    onClick={handleSave}
                    disabled={isSaving}
                    className="p-2 rounded-lg transition disabled:opacity-50"
                    style={{ background: "var(--success-bg)", color: "var(--success-text)" }}
                    onMouseEnter={e => (e.currentTarget.style.background = "var(--success-bg-hover)")}
                    onMouseLeave={e => (e.currentTarget.style.background = "var(--success-bg)")}
                  >
                    <Check size={18} />
                  </button>
                  <button
                    onClick={() => { setIsEditing(false); setEditName(user?.name || ""); }}
                    className="p-2 rounded-lg transition"
                    style={{ background: "var(--chat-bg)", color: "var(--text-soft)" }}
                    onMouseEnter={e => { e.currentTarget.style.background = "var(--surface-muted)"; e.currentTarget.style.color = "var(--text)"; }}
                    onMouseLeave={e => { e.currentTarget.style.background = "var(--chat-bg)"; e.currentTarget.style.color = "var(--text-soft)"; }}
                  >
                    <X size={18} />
                  </button>
                </>
              ) : (
                <button
                  onClick={() => setIsEditing(true)}
                  className="p-2 rounded-lg border transition"
                  style={{ background: "var(--chat-bg)", borderColor: "var(--panel-border)", color: "var(--text-soft)" }}
                  onMouseEnter={e => { e.currentTarget.style.background = "var(--surface-muted)"; e.currentTarget.style.color = "var(--text)"; }}
                  onMouseLeave={e => { e.currentTarget.style.background = "var(--chat-bg)"; e.currentTarget.style.color = "var(--text-soft)"; }}
                  title="Edit Name"
                >
                  <Edit2 size={18} />
                </button>
              )}
            </div>
          </div>
        </div>

        {/* Divider */}
        <hr className="my-8" style={{ borderColor: "var(--panel-border)" }} />
      </div>
    </div>
  );
}