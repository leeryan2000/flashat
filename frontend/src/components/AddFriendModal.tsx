import { useState } from "react";
import { X, UserPlus, Loader2 } from "lucide-react";

type AddFriendModalProps = {
  onClose: () => void;
  onAddFriend: (email: string) => Promise<void>;
};

export function AddFriendModal({ onClose, onAddFriend }: AddFriendModalProps) {
  const [addEmail, setAddEmail] = useState("");
  const [isAdding, setIsAdding] = useState(false);
  const [addError, setAddError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!addEmail.trim()) return;

    setIsAdding(true);
    setAddError("");

    try {
      await onAddFriend(addEmail);
      setAddEmail("");
      onClose();
    } catch {
      setAddError("Failed to send request. Check the email.");
    } finally {
      setIsAdding(false);
    }
  };

  return (
    <div className="absolute inset-0 z-50 flex items-center justify-center p-4" style={{ background: "var(--overlay)" }}>
      <div className="border p-6 rounded-2xl w-full max-w-md shadow-2xl animate-in zoom-in-95 duration-200" style={{ background: "var(--sidebar-item)", borderColor: "color-mix(in srgb, var(--primary) 20%, transparent)" }}>
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-xl font-semibold" style={{ color: "var(--text)" }}>Send Friend Request</h3>
          <button
            onClick={onClose}
            className="transition"
            style={{ color: "var(--text-soft)" }}
            onMouseEnter={e => (e.currentTarget.style.color = "var(--text)")}
            onMouseLeave={e => (e.currentTarget.style.color = "var(--text-soft)")}
          >
            <X size={20} />
          </button>
        </div>

        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-sm font-medium mb-1" style={{ color: "var(--text-soft)" }}>
              Friend's Email Address
            </label>
            <input
              type="email"
              autoFocus
              required
              value={addEmail}
              onChange={(e) => setAddEmail(e.target.value)}
              placeholder="name@example.com"
              className="w-full border rounded-xl px-4 py-3 focus:ring-2 outline-none placeholder:text-[color:var(--text-faint)]"
              style={{ background: "var(--chat-bg)", color: "var(--text)", borderColor: "color-mix(in srgb, var(--primary) 20%, transparent)", "--tw-ring-color": "var(--primary)" } as React.CSSProperties}
            />
            {addError && <p className="text-sm mt-2" style={{ color: "var(--danger-text)" }}>{addError}</p>}
          </div>

          <div className="flex justify-end gap-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 rounded-lg transition"
              style={{ color: "var(--text-soft)" }}
              onMouseEnter={e => { e.currentTarget.style.background = "color-mix(in srgb, var(--primary) 10%, transparent)"; e.currentTarget.style.color = "var(--text)"; }}
              onMouseLeave={e => { e.currentTarget.style.background = ""; e.currentTarget.style.color = "var(--text-soft)"; }}
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isAdding}
              className="px-4 py-2 text-white rounded-lg font-medium transition flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
              style={{ background: "var(--primary)" }}
            >
              {isAdding ? <Loader2 className="animate-spin" size={18} /> : <UserPlus size={18} />}
              Send Request
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
