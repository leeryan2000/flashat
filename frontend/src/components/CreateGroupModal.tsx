import { useState } from "react";
import { X, Check, Loader2, MessageSquarePlus } from "lucide-react";
import type { Friendship } from "../wire/friendship";
import Avatar from "./Avatar";

type CreateGroupModalProps = {
  friends: Friendship[];
  onClose: () => void;
  onCreateGroup: (name: string, participantIds: string[]) => Promise<void>;
};

export function CreateGroupModal({ friends, onClose, onCreateGroup }: CreateGroupModalProps) {
  const [groupName, setGroupName] = useState("");
  const [selectedFriendIds, setSelectedFriendIds] = useState<string[]>([]);
  const [isCreatingGroup, setIsCreatingGroup] = useState(false);
  const [groupError, setGroupError] = useState("");

  const toggleFriendSelection = (uid: string) => {
    setSelectedFriendIds(prev =>
      prev.includes(uid)
        ? prev.filter(id => id !== uid)
        : [...prev, uid]
    );
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!groupName.trim() || selectedFriendIds.length === 0) {
      setGroupError("Please enter a name and select at least one friend.");
      return;
    }

    setIsCreatingGroup(true);
    setGroupError("");

    try {
      await onCreateGroup(groupName, selectedFriendIds);
      setGroupName("");
      setSelectedFriendIds([]);
      onClose();
    } catch {
      setGroupError("Failed to create group.");
    } finally {
      setIsCreatingGroup(false);
    }
  };

  const renderAvatar = (name: string, avatarUrl?: string | null) => (
    <Avatar name={name} avatarUrl={avatarUrl} size="sm" />
  );

  return (
    <div className="absolute inset-0 z-50 flex items-center justify-center p-4" style={{ background: "var(--overlay)" }}>
      <div className="border p-6 rounded-2xl w-full max-w-lg shadow-2xl animate-in zoom-in-95 duration-200 flex flex-col max-h-[90vh]" style={{ background: "var(--sidebar-item)", borderColor: "color-mix(in srgb, var(--primary) 20%, transparent)" }}>
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-xl font-semibold" style={{ color: "var(--text)" }}>Create Group</h3>
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

        <form onSubmit={handleSubmit} className="flex flex-col flex-1 overflow-hidden">
          <div className="mb-4">
            <label className="block text-sm font-medium mb-1" style={{ color: "var(--text-soft)" }}>
              Group Name
            </label>
            <input
              type="text"
              autoFocus
              required
              value={groupName}
              onChange={(e) => setGroupName(e.target.value)}
              placeholder="e.g. Project Team"
              className="w-full border rounded-xl px-4 py-3 focus:ring-2 outline-none placeholder:text-[color:var(--text-faint)]"
              style={{ background: "var(--chat-bg)", color: "var(--text)", borderColor: "color-mix(in srgb, var(--primary) 20%, transparent)", "--tw-ring-color": "var(--primary)" } as React.CSSProperties}
            />
          </div>

          <div className="flex-1 overflow-hidden flex flex-col">
            <label className="block text-sm font-medium mb-2" style={{ color: "var(--text-soft)" }}>
              Select Members ({selectedFriendIds.length})
            </label>

            <div className="flex-1 overflow-y-auto border rounded-xl p-2 space-y-1" style={{ background: "var(--chat-bg)", borderColor: "color-mix(in srgb, var(--primary) 15%, transparent)" }}>
              {friends.length === 0 ? (
                <p className="text-sm text-center py-4" style={{ color: "var(--text-faint)" }}>You have no friends to add yet.</p>
              ) : (
                friends.map(friend => {
                  const isSelected = selectedFriendIds.includes(friend.uid);
                  return (
                    <div
                      key={friend.uid}
                      onClick={() => toggleFriendSelection(friend.uid)}
                      className="flex items-center gap-3 p-2 rounded-lg cursor-pointer transition border"
                      style={isSelected ? { background: "color-mix(in srgb, var(--primary) 15%, transparent)", borderColor: "color-mix(in srgb, var(--primary) 50%, transparent)" } : { borderColor: "transparent" }}
                    >
                      <div
                        className="w-5 h-5 rounded-md border flex items-center justify-center transition"
                        style={isSelected ? { background: "var(--primary)", borderColor: "var(--primary)" } : { borderColor: "var(--panel-border)" }}
                      >
                        {isSelected && <Check size={14} className="text-white" />}
                      </div>

                      {renderAvatar(friend.name, friend.avatarUrl)}

                      <div className="min-w-0">
                        <p className="text-sm font-medium" style={{ color: "var(--text)" }}>
                          {friend.name}
                        </p>
                      </div>
                    </div>
                  );
                })
              )}
            </div>
          </div>

          {groupError && <p className="text-sm mt-4" style={{ color: "var(--danger-text)" }}>{groupError}</p>}

          <div className="flex justify-end gap-3 mt-6 pt-4 border-t" style={{ borderColor: "var(--panel-border)" }}>
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
              disabled={isCreatingGroup || selectedFriendIds.length === 0}
              className="px-4 py-2 text-white rounded-lg font-medium transition flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
              style={{ background: "var(--primary)" }}
            >
              {isCreatingGroup ? <Loader2 className="animate-spin" size={18} /> : <MessageSquarePlus size={18} />}
              Create Group
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
