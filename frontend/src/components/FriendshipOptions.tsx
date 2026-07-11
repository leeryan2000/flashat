import { MoreVertical, UserX, Ban } from 'lucide-react';
import { useState } from 'react';

// 🛠️ Helper Component for the Dropdown
export function FriendshipOptions({ onUnfriend, onBlock }: { onUnfriend: () => void; onBlock: () => void }) {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="relative">
      {/* 1. The Trigger Button (Three Dots) */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="p-2 rounded-lg transition"
        style={isOpen ? { background: "color-mix(in srgb, var(--primary) 20%, transparent)", color: "var(--primary)" } : { color: "var(--text-faint)" }}
        onMouseEnter={e => { if (!isOpen) { e.currentTarget.style.background = "color-mix(in srgb, var(--primary) 10%, transparent)"; e.currentTarget.style.color = "var(--text-soft)"; } }}
        onMouseLeave={e => { if (!isOpen) { e.currentTarget.style.background = ""; e.currentTarget.style.color = "var(--text-faint)"; } }}
        title="Options"
      >
        <MoreVertical size={18} />
      </button>

      {/* 2. The Dropdown Menu */}
      {isOpen && (
        <>
          {/* Invisible Backdrop to close menu when clicking outside */}
          <div
            className="fixed inset-0 z-10 cursor-default"
            onClick={() => setIsOpen(false)}
          />

          {/* The Actual Menu List */}
          <div
            className="absolute right-0 mt-2 w-48 border rounded-xl shadow-xl z-20 overflow-hidden py-1 animate-in fade-in zoom-in-95 duration-100"
            style={{ background: "var(--sidebar-item)", borderColor: "var(--panel-border)" }}
          >

            {/* Option: Unfriend */}
            <button
              onClick={() => {
                onUnfriend();
                setIsOpen(false);
              }}
              className="w-full text-left px-4 py-3 text-sm flex items-center gap-3 transition"
              style={{ color: "var(--danger-text)" }}
              onMouseEnter={e => (e.currentTarget.style.background = "var(--danger-bg)")}
              onMouseLeave={e => (e.currentTarget.style.background = "")}
            >
              <UserX size={16} />
              Unfriend User
            </button>

            {/* Option: Block */}
            <button
              onClick={() => {
                onBlock();
                setIsOpen(false);
              }}
              className="w-full text-left px-4 py-3 text-sm flex items-center gap-3 transition"
              style={{ color: "var(--warning-text)" }}
              onMouseEnter={e => (e.currentTarget.style.background = "var(--warning-bg)")}
              onMouseLeave={e => (e.currentTarget.style.background = "")}
            >
              <Ban size={16} />
              Block User
            </button>
          </div>
        </>
      )}
    </div>
  );
}