import { useState, useMemo } from 'react';
import {
  MessageSquare,
  UserPlus,
  Users,
  X,
  Check,
  Search,
  MessageSquarePlus,
  Ban,
  ChevronDown,
  ChevronRight,
  ShieldOff,
} from 'lucide-react';
import type { Friendship } from '../wire/friendship';
import { FriendshipOptions } from './FriendshipOptions';
import { AddFriendModal } from './AddFriendModal';
import { CreateGroupModal } from './CreateGroupModal';

type FriendsListProps = {
  data: Friendship[];
  onChatClick: (convId: string) => void;
  onAccept: (uid: string) => void;
  onDecline: (uid: string) => void;
  onCancel: (uid: string) => void;
  onUnfriend: (uid: string, convId: string | null) => void;
  onBlock: (uid: string, convId: string | null) => void;
  onUnblock: (uid: string) => void;
  onAddFriend: (email: string) => Promise<void>;
  onCreateGroup: (name: string, participantIds: string[]) => Promise<void>;
}

export default function FriendshipList({
  data,
  onChatClick,
  onAccept,
  onDecline,
  onCancel,
  onUnfriend,
  onBlock,
  onUnblock,
  onAddFriend,
  onCreateGroup,
}: FriendsListProps) {
  const [query, setQuery] = useState("");
  const [blockedOpen, setBlockedOpen] = useState(false);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [isGroupModalOpen, setIsGroupModalOpen] = useState(false);

  const { friends, incoming, outgoing, blocked } = useMemo(() => {
    const filtered = data.filter(u =>
      u.name.toLowerCase().includes(query.toLowerCase()) ||
      u.email.toLowerCase().includes(query.toLowerCase())
    );

    return {
      friends: filtered.filter(u => u.status === 'accepted'),
      incoming: filtered.filter(u => u.status === 'pending' && u.direction === 'incoming'),
      outgoing: filtered.filter(u => u.status === 'pending' && u.direction === 'outgoing'),
      blocked: data.filter(u => u.status === 'blocked' && u.direction === 'outgoing'),
    };
  }, [data, query]);

  const renderAvatar = (name: string) => (
    <div className="h-10 w-10 rounded-full flex items-center justify-center text-white font-semibold shrink-0" style={{ background: "var(--primary)" }}>
      {name.slice(0, 2).toUpperCase()}
    </div>
  );

  return (
    <div className="h-full w-full text-slate-100 flex flex-col" style={{ background: "var(--sidebar-bg)" }}>
      {/* --- Header --- */}
      <div className="p-6 border-b border-slate-700/50 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">

        <h1 className="text-2xl font-bold flex items-center gap-3">
          <Users className="w-6 h-6" style={{ color: "var(--primary)" }} />
          Friends
        </h1>

        <div className="flex gap-3">
          {/* Search Input */}
          <div className="relative w-64">
            <Search className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
            <input
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Search friends..."
              className="w-full rounded-xl pl-9 pr-4 py-2 text-sm text-slate-100 placeholder:text-slate-400 outline-none focus:ring-2 transition"
              style={{ background: "var(--sidebar-item)", "--tw-ring-color": "var(--primary)" } as React.CSSProperties}
              style={{ "--tw-ring-color": "var(--primary)" } as React.CSSProperties}
            />
          </div>

          {/* Add Friend Button */}
          <button
            onClick={() => setIsAddModalOpen(true)}
            className="flex items-center gap-2 px-4 py-2 text-white rounded-xl font-medium transition text-sm shadow-lg"
            style={{ background: "var(--primary)", boxShadow: "0 4px 12px var(--primary-shadow)" }}
          >
            <UserPlus size={18} />
          </button>

          {/* Create Group Button */}
          <button
            onClick={() => setIsGroupModalOpen(true)}
            className="flex items-center gap-2 px-4 py-2 text-slate-200 rounded-xl font-medium transition text-sm"
            style={{ background: "var(--sidebar-item)" }}
          >
            <MessageSquarePlus size={18} />
            <span className="hidden sm:inline">New Group</span>
          </button>
        </div>
      </div>

      {/* --- Main Content (Scrollable) --- */}
      <div className="flex-1 overflow-y-auto p-6 space-y-8">

        {/* Section 1: Incoming Requests */}
        {incoming.length > 0 && (
          <section>
            <h2 className="text-sm font-semibold text-slate-400 uppercase tracking-wider mb-4">
              Friend Requests ({incoming.length})
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {incoming.map(user => (
                <div key={user.uid} className="p-4 rounded-xl border border-slate-700/20 flex items-center gap-4" style={{ background: "var(--sidebar-item)" }}>
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
                    <button
                      onClick={() => onBlock(user.uid, null)}
                      className="p-2 rounded-lg bg-orange-500/20 text-orange-400 hover:bg-orange-500/30 transition"
                      title="Block"
                    >
                      <Ban size={18} />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </section>
        )}

        {/* Section 2: Sent Requests */}
        {outgoing.length > 0 && (
          <section>
            <h2 className="text-sm font-semibold text-slate-400 uppercase tracking-wider mb-4">
              Sent Requests ({outgoing.length})
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {outgoing.map(user => (
                <div key={user.uid} className="p-4 rounded-xl border border-slate-700/20 flex items-center gap-4 opacity-75" style={{ background: "var(--sidebar-item)" }}>
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
                  className="group transition p-4 rounded-xl border border-transparent flex items-center gap-4"
                  style={{ background: "var(--sidebar-item)" }}
                >
                  {renderAvatar(user.name)}

                  <div className="flex-1 min-w-0">
                    <p className="font-medium truncate text-slate-200 group-hover:text-white">
                      {user.name}
                    </p>
                    <p className="text-xs text-slate-400 truncate">{user.email}</p>
                  </div>

                  <div className="flex items-center gap-2">
                    <button
                      onClick={() => user.directConversationId && onChatClick(user.directConversationId)}
                      className="p-2 rounded-lg transition hover:text-white"
                      style={{ background: "color-mix(in srgb, var(--primary) 15%, transparent)", color: "var(--primary)" }}
                      title="Send Message"
                    >
                      <MessageSquare size={18} />
                    </button>
                    <FriendshipOptions
                      onUnfriend={() => onUnfriend(user.uid, user.directConversationId)}
                      onBlock={() => onBlock(user.uid, user.directConversationId)}
                    />
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>

        {/* Section 4: Blocked Users */}
        {blocked.length > 0 && (
          <section>
            <button
              onClick={() => setBlockedOpen(o => !o)}
              className="flex items-center gap-2 text-sm font-semibold text-slate-400 uppercase tracking-wider mb-4 hover:text-slate-300 transition"
            >
              {blockedOpen ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
              Blocked ({blocked.length})
            </button>

            {blockedOpen && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {blocked.map(user => (
                  <div key={user.uid} className="p-4 rounded-xl border border-slate-700/20 flex items-center gap-4 opacity-75" style={{ background: "var(--sidebar-item)" }}>
                    {renderAvatar(user.name)}
                    <div className="flex-1 min-w-0">
                      <p className="font-medium truncate">{user.name}</p>
                      <p className="text-xs text-slate-400 truncate">{user.email}</p>
                    </div>
                    <button
                      onClick={() => onUnblock(user.uid)}
                      className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium bg-slate-700 text-slate-300 hover:bg-slate-600 hover:text-white transition"
                      title="Unblock"
                    >
                      <ShieldOff size={14} />
                      Unblock
                    </button>
                  </div>
                ))}
              </div>
            )}
          </section>
        )}
      </div>

      {/* Modals — isolated components, their state changes don't re-render this list */}
      {isAddModalOpen && (
        <AddFriendModal
          onClose={() => setIsAddModalOpen(false)}
          onAddFriend={onAddFriend}
        />
      )}

      {isGroupModalOpen && (
        <CreateGroupModal
          friends={friends}
          onClose={() => setIsGroupModalOpen(false)}
          onCreateGroup={onCreateGroup}
        />
      )}
    </div>
  );
}
