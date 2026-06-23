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
    <div className="absolute inset-0 z-50 flex items-center justify-center bg-slate-900/80 p-4">
      <div className="border p-6 rounded-2xl w-full max-w-md shadow-2xl animate-in zoom-in-95 duration-200" style={{ background: "var(--sidebar-item)", borderColor: "color-mix(in srgb, var(--primary) 20%, transparent)" }}>
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-xl font-semibold text-white">Send Friend Request</h3>
          <button onClick={onClose} className="text-slate-400 hover:text-white">
            <X size={20} />
          </button>
        </div>

        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-sm font-medium text-slate-400 mb-1">
              Friend's Email Address
            </label>
            <input
              type="email"
              autoFocus
              required
              value={addEmail}
              onChange={(e) => setAddEmail(e.target.value)}
              placeholder="name@example.com"
              className="w-full border rounded-xl px-4 py-3 text-slate-100 placeholder:text-slate-500 focus:ring-2 outline-none"
              style={{ background: "var(--sidebar-bg)", borderColor: "color-mix(in srgb, var(--primary) 20%, transparent)", "--tw-ring-color": "var(--primary)" } as React.CSSProperties}
            />
            {addError && <p className="text-red-400 text-sm mt-2">{addError}</p>}
          </div>

          <div className="flex justify-end gap-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-slate-300 hover:text-white rounded-lg transition"
              onMouseEnter={e => (e.currentTarget.style.background = "color-mix(in srgb, var(--primary) 10%, transparent)")}
              onMouseLeave={e => (e.currentTarget.style.background = "")}
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
