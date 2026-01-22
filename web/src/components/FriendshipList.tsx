import React, { useState, useMemo } from 'react';
import { 
  MessageSquare, 
  UserPlus, 
  X, 
  Check, 
  Search, 
  UserX 
} from 'lucide-react'; // Assuming you use lucide-react for icons
import type { Friendship } from '../wire/friendship';
import { FriendshipOptions } from './FriendshipOptions';

// 1. Define the Data Shape

type FriendsListProps = {
  data: Friendship[]; // Pass your JSON data here
  onChatClick: (convId: string) => void;
  onAccept: (uid: string) => void;
  onDecline: (uid: string) => void;
  onCancel: (uid: string) => void;
  onUnfriend: (uid: string, convId: string | null) => void;
}

export default function FriendshipList({ 
  data, 
  onChatClick, 
  onAccept, 
  onDecline, 
  onCancel,
  onUnfriend
}: FriendsListProps) {
  const [query, setQuery] = useState("");

  // 2. Filter & Split Data
  const { friends, incoming, outgoing } = useMemo(() => {
    // First, filter by search query
    const filtered = data.filter(u => 
      u.name.toLowerCase().includes(query.toLowerCase()) || 
      u.email.toLowerCase().includes(query.toLowerCase())
    );

    return {
      friends: filtered.filter(u => u.status === 'accepted'),
      incoming: filtered.filter(u => u.status === 'pending' && u.direction === 'incoming'),
      outgoing: filtered.filter(u => u.status === 'pending' && u.direction === 'outgoing'),
    };
  }, [data, query]);

  // 3. Helper to render Avatar
  const renderAvatar = (name: string) => (
    <div className="h-10 w-10 rounded-full bg-slate-700 flex items-center justify-center text-slate-300 font-semibold shrink-0">
      {name.slice(0, 2).toUpperCase()}
    </div>
  );

  return (
    <div className="h-full w-full bg-slate-900 text-slate-100 flex flex-col">
      {/* --- Header --- */}
      <div className="p-6 border-b border-slate-800 flex items-center justify-between">
        <h1 className="text-2xl font-bold flex items-center gap-3">
          <UserPlus className="w-6 h-6 text-indigo-400" />
          Friends
        </h1>
        
        {/* Search Input */}
        <div className="relative w-64">
          <Search className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search friends..."
            className="w-full rounded-xl bg-slate-800/60 pl-9 pr-4 py-2 text-sm placeholder:text-slate-400 outline-none focus:ring-2 focus:ring-indigo-400 transition"
          />
        </div>
      </div>

      {/* --- Main Content (Scrollable) --- */}
      <div className="flex-1 overflow-y-auto p-6 space-y-8">
        
        {/* Section 1: Incoming Requests (Only show if exists) */}
        {incoming.length > 0 && (
          <section>
            <h2 className="text-sm font-semibold text-slate-400 uppercase tracking-wider mb-4">
              Friend Requests ({incoming.length})
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {incoming.map(user => (
                <div key={user.uid} className="bg-slate-800/40 p-4 rounded-xl border border-slate-700/50 flex items-center gap-4">
                  {renderAvatar(user.name)}
                  <div className="flex-1 min-w-0">
                    <p className="font-medium truncate">{user.name}</p>
                    <p className="text-xs text-slate-400 truncate">{user.email}</p>
                  </div>
                  <div className="flex gap-2">
                    <button 
                      onClick={() => onAccept(user.uid)}
                      className="p-2 rounded-lg bg-emerald-500/20 text-emerald-400 hover:bg-emerald-500/30 transition"
                      title="Accept"
                    >
                      <Check size={18} />
                    </button>
                    <button 
                      onClick={() => onDecline(user.uid)}
                      className="p-2 rounded-lg bg-red-500/20 text-red-400 hover:bg-red-500/30 transition"
                      title="Decline"
                    >
                      <X size={18} />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </section>
        )}

        {/* Section 2: Sent Requests (Only show if exists) */}
        {outgoing.length > 0 && (
          <section>
            <h2 className="text-sm font-semibold text-slate-400 uppercase tracking-wider mb-4">
              Sent Requests ({outgoing.length})
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {outgoing.map(user => (
                <div key={user.uid} className="bg-slate-800/40 p-4 rounded-xl border border-slate-700/50 flex items-center gap-4 opacity-75">
                  {renderAvatar(user.name)}
                  <div className="flex-1 min-w-0">
                    <p className="font-medium truncate">{user.name}</p>
                    <p className="text-xs text-slate-400 truncate">Pending...</p>
                  </div>
                  <button 
                    onClick={() => onCancel(user.uid)}
                    className="text-xs text-slate-400 hover:text-red-400 underline decoration-slate-600 hover:decoration-red-400 transition"
                  >
                    Cancel
                  </button>
                </div>
              ))}
            </div>
          </section>
        )}

        {/* Section 3: All Friends */}
        <section>
          <h2 className="text-sm font-semibold text-slate-400 uppercase tracking-wider mb-4">
            All Friends ({friends.length})
          </h2>
          
          {friends.length === 0 ? (
            <div className="text-slate-500 text-center py-10 italic">
              No friends found. Go add some people!
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
              {friends.map(user => (
                <div 
                  key={user.uid} 
                  className="group bg-slate-800 hover:bg-slate-700/80 transition p-4 rounded-xl border border-transparent hover:border-slate-600 flex items-center gap-4"
                >
                  {renderAvatar(user.name)}
                  
                  <div className="flex-1 min-w-0">
                    <p className="font-medium truncate text-slate-200 group-hover:text-white">
                      {user.name}
                    </p>
                    <p className="text-xs text-slate-400 truncate">{user.email}</p>
                  </div>

                  {/* Actions */}
                  <div className="flex items-center gap-2">
                    <button 
                      onClick={() => user.directConversationId && onChatClick(user.directConversationId)}
                      className="p-2 rounded-lg bg-indigo-500/10 text-indigo-400 hover:bg-indigo-500 hover:text-white transition"
                      title="Send Message"
                    >
                      <MessageSquare size={18} />
                    </button>
                    {/* Optional: Unfriend Button */}
                    <FriendshipOptions onUnfriend={() => onUnfriend(user.uid, user.directConversationId)} />
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>
      </div>
    </div>
  );
}