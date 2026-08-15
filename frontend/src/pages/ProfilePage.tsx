import { useRef, useState } from 'react';
import { useAuth } from '../context/AuthContext';
import {
  User,
  Mail,
  Camera,
  Check,
  X,
  Edit2,
  Loader2,
  Upload,
  Trash2
} from 'lucide-react';
import Avatar from '../components/Avatar';

// Resizes/compresses an image client-side before upload, so avatars stay
// small in S3 regardless of how large the source photo is.
async function resizeImageToBlob(file: File, maxSize = 512, quality = 0.85): Promise<Blob> {
  const objectUrl = URL.createObjectURL(file);
  try {
    const img = document.createElement("img");
    await new Promise<void>((resolve, reject) => {
      img.onload = () => resolve();
      img.onerror = () => reject(new Error("Failed to read image"));
      img.src = objectUrl;
    });

    const scale = Math.min(1, maxSize / Math.max(img.width, img.height));
    const width = Math.round(img.width * scale);
    const height = Math.round(img.height * scale);

    const canvas = document.createElement("canvas");
    canvas.width = width;
    canvas.height = height;
    const ctx = canvas.getContext("2d");
    if (!ctx) throw new Error("Canvas is not supported");
    ctx.drawImage(img, 0, 0, width, height);

    return await new Promise<Blob>((resolve, reject) => {
      canvas.toBlob((blob) => {
        if (blob) resolve(blob);
        else reject(new Error("Failed to process image"));
      }, "image/jpeg", quality);
    });
  } finally {
    URL.revokeObjectURL(objectUrl);
  }
}

export default function ProfilePage() {
  // Pull current user data and logout function from your context
  const { user, updateProfile, getAvatarUploadUrl, updateAvatarUrl, removeAvatar } = useAuth();

  // State for toggling edit mode
  const [isEditing, setIsEditing] = useState(false);
  const [editName, setEditName] = useState(user?.name || "");
  const [isSaving, setIsSaving] = useState(false);

  // Avatar upload state
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [isAvatarMenuOpen, setIsAvatarMenuOpen] = useState(false);
  const [isAvatarBusy, setIsAvatarBusy] = useState(false);
  const [avatarError, setAvatarError] = useState("");

  const handleUploadOptionClick = () => {
    setIsAvatarMenuOpen(false);
    fileInputRef.current?.click();
  };

  const handleAvatarFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    e.target.value = ""; // allow re-selecting the same file next time
    if (!file) return;

    setIsAvatarBusy(true);
    setAvatarError("");
    try {
      const blob = await resizeImageToBlob(file);
      const { upload_url, object_url } = await getAvatarUploadUrl();

      const putRes = await fetch(upload_url, {
        method: "PUT",
        body: blob,
        headers: { "Content-Type": "image/jpeg" },
      });
      if (!putRes.ok) throw new Error("Failed to upload photo");

      await updateAvatarUrl(object_url);
    } catch (error) {
      console.error("Failed to upload avatar:", error);
      setAvatarError("Failed to upload photo. Please try again.");
    } finally {
      setIsAvatarBusy(false);
    }
  };

  const handleRemoveOptionClick = async () => {
    setIsAvatarMenuOpen(false);
    setIsAvatarBusy(true);
    setAvatarError("");
    try {
      await removeAvatar();
    } catch (error) {
      console.error("Failed to remove avatar:", error);
      setAvatarError("Failed to remove photo. Please try again.");
    } finally {
      setIsAvatarBusy(false);
    }
  };

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
          <div className="relative group cursor-pointer mb-4" onClick={() => setIsAvatarMenuOpen(o => !o)}>
            {/* Avatar Circle */}
            <Avatar name={user?.name || "??"} avatarUrl={user?.user_avatar_url} size="lg" className="shadow-lg" />

            <input
              ref={fileInputRef}
              type="file"
              accept="image/*"
              className="hidden"
              onChange={handleAvatarFileChange}
            />

            {/* Hover Overlay for changing avatar */}
            <div className="absolute inset-0 bg-black/50 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-200">
              {isAvatarBusy ? (
                <Loader2 className="w-8 h-8 text-white animate-spin" />
              ) : (
                <Camera className="w-8 h-8 text-white" />
              )}
            </div>

            {/* Upload / Remove menu */}
            {isAvatarMenuOpen && (
              <>
                <div
                  className="fixed inset-0 z-10 cursor-default"
                  onClick={(e) => { e.stopPropagation(); setIsAvatarMenuOpen(false); }}
                />
                <div
                  className="absolute left-1/2 -translate-x-1/2 top-full mt-2 w-48 border rounded-xl shadow-xl z-20 overflow-hidden py-1 animate-in fade-in zoom-in-95 duration-100 cursor-default"
                  style={{ background: "var(--sidebar-item)", borderColor: "var(--panel-border)" }}
                  onClick={(e) => e.stopPropagation()}
                >
                  <button
                    onClick={handleUploadOptionClick}
                    className="w-full text-left px-4 py-3 text-sm flex items-center gap-3 transition"
                    style={{ color: "var(--text)" }}
                    onMouseEnter={e => (e.currentTarget.style.background = "var(--surface-muted)")}
                    onMouseLeave={e => (e.currentTarget.style.background = "")}
                  >
                    <Upload size={16} />
                    Upload Photo
                  </button>

                  {user?.user_avatar_url && (
                    <button
                      onClick={handleRemoveOptionClick}
                      className="w-full text-left px-4 py-3 text-sm flex items-center gap-3 transition"
                      style={{ color: "var(--danger-text)" }}
                      onMouseEnter={e => (e.currentTarget.style.background = "var(--danger-bg)")}
                      onMouseLeave={e => (e.currentTarget.style.background = "")}
                    >
                      <Trash2 size={16} />
                      Remove Photo
                    </button>
                  )}
                </div>
              </>
            )}
          </div>

          {avatarError && <p className="text-sm mb-2" style={{ color: "var(--danger-text)" }}>{avatarError}</p>}

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