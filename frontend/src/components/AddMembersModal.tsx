import { useState } from "react";
import { X, Check, Loader2, UserPlus } from "lucide-react";
import type { Friendship } from "../wire/friendship";
import { avatarColor, nameInitials } from "../utils/avatar";

type AddMembersModalProps = {
  friends: Friendship[];
  onClose: () => void;
  onAdd: (selected: { uid: string; name: string }[]) => Promise<void>;
};

export function AddMembersModal({ friends, onClose, onAdd }: AddMembersModalProps) {
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [isAdding, setIsAdding] = useState(false);
  const [error, setError] = useState("");

  const toggle = (uid: string) => {
    setSelectedIds((prev) =>
      prev.includes(uid) ? prev.filter((id) => id !== uid) : [...prev, uid]
    );
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (selectedIds.length === 0) return;

    setIsAdding(true);
    setError("");

    const selected = friends
      .filter((f) => selectedIds.includes(f.uid))
      .map((f) => ({ uid: f.uid, name: f.name }));

    try {
      await onAdd(selected);
      onClose();
    } catch {
      setError("Failed to add members. Please try again.");
    } finally {
      setIsAdding(false);
    }
  };

  return (
    <div className="absolute inset-0 z-50 flex items-center justify-center bg-slate-900/80 p-4">
      <div
        className="border p-6 rounded-2xl w-full max-w-lg shadow-2xl animate-in zoom-in-95 duration-200 flex flex-col max-h-[90vh]"
        style={{
          background: "var(--sidebar-item)",
          borderColor: "color-mix(in srgb, var(--primary) 20%, transparent)",
        }}
      >
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-xl font-semibold text-white">Add Members</h3>
          <button onClick={onClose} className="text-slate-400 hover:text-white transition">
            <X size={20} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="flex flex-col flex-1 overflow-hidden">
          <div className="flex-1 overflow-hidden flex flex-col">
            <label className="block text-sm font-medium text-slate-400 mb-2">
              Select Friends ({selectedIds.length} selected)
            </label>

            <div
              className="flex-1 overflow-y-auto border rounded-xl p-2 space-y-1"
              style={{
                background: "var(--sidebar-bg)",
                borderColor: "color-mix(in srgb, var(--primary) 15%, transparent)",
              }}
            >
              {friends.length === 0 ? (
                <p className="text-slate-500 text-sm text-center py-4">
                  All your friends are already in this group.
                </p>
              ) : (
                friends.map((friend) => {
                  const isSelected = selectedIds.includes(friend.uid);
                  return (
                    <div
                      key={friend.uid}
                      onClick={() => toggle(friend.uid)}
                      className="flex items-center gap-3 p-2 rounded-lg cursor-pointer transition border"
                      style={
                        isSelected
                          ? {
                              background: "color-mix(in srgb, var(--primary) 15%, transparent)",
                              borderColor: "color-mix(in srgb, var(--primary) 50%, transparent)",
                            }
                          : { borderColor: "transparent" }
                      }
                    >
                      <div
                        className="w-5 h-5 rounded-md border flex items-center justify-center transition shrink-0"
                        style={
                          isSelected
                            ? { background: "var(--primary)", borderColor: "var(--primary)" }
                            : { borderColor: "#475569" }
                        }
                      >
                        {isSelected && <Check size={14} className="text-white" />}
                      </div>

                      <div
                        className={`h-10 w-10 rounded-full flex items-center justify-center text-white text-sm font-semibold shrink-0 ${avatarColor(friend.name)}`}
                      >
                        {nameInitials(friend.name)}
                      </div>

                      <p className="text-sm font-medium text-slate-200 truncate">
                        {friend.name}
                      </p>
                    </div>
                  );
                })
              )}
            </div>
          </div>

          {error && <p className="text-red-400 text-sm mt-4">{error}</p>}

          <div className="flex justify-end gap-3 mt-6 pt-4 border-t border-slate-700">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-slate-300 hover:text-white rounded-lg transition"
              onMouseEnter={(e) =>
                (e.currentTarget.style.background = "color-mix(in srgb, var(--primary) 10%, transparent)")
              }
              onMouseLeave={(e) => (e.currentTarget.style.background = "")}
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isAdding || selectedIds.length === 0}
              className="px-4 py-2 text-white rounded-lg font-medium transition flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
              style={{ background: "var(--primary)" }}
            >
              {isAdding ? (
                <Loader2 className="animate-spin" size={18} />
              ) : (
                <UserPlus size={18} />
              )}
              Add {selectedIds.length > 0 ? `(${selectedIds.length})` : ""}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
