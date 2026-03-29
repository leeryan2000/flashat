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
        className={`p-2 rounded-lg transition ${isOpen ? 'bg-slate-700 text-white' : 'text-slate-500 hover:bg-slate-700/50 hover:text-slate-300'}`}
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
          <div className="absolute right-0 mt-2 w-48 bg-slate-800 border border-slate-700 rounded-xl shadow-xl z-20 overflow-hidden py-1 animate-in fade-in zoom-in-95 duration-100">
            
            {/* Option: Unfriend */}
            <button 
              onClick={() => {
                onUnfriend();
                setIsOpen(false);
              }}
              className="w-full text-left px-4 py-3 text-sm text-red-400 hover:bg-red-500/10 hover:text-red-300 flex items-center gap-3 transition"
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
              className="w-full text-left px-4 py-3 text-sm text-orange-400 hover:bg-orange-500/10 hover:text-orange-300 flex items-center gap-3 transition"
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