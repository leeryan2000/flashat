import React, { useState, useMemo } from 'react';
import { 
  MessageSquare, 
  UserPlus,
  Users, 
  X, 
  Check, 
  Search, 
  Loader2,
  MessageSquarePlus,
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
  onAddFriend: (email: string) => void;
  onCreateGroup: (name: string, participantIds: string[]) => Promise<void>;
}

export default function FriendshipList({ 
  data, 
  onChatClick, 
  onAccept, 
  onDecline, 
  onCancel,
  onUnfriend,
  onAddFriend,
  onCreateGroup
}: FriendsListProps) {
  const [query, setQuery] = useState("");

  // -- State for "Add Friend" ---
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [addEmail, setAddEmail] = useState("");
  const [isAdding, setIsAdding] = useState(false);
  const [addError, setAddError] = useState("");

  // --- State for "Create Group Conversation" ---
  const [isGroupModalOpen, setIsGroupModalOpen] = useState(false);
  const [groupName, setGroupName] = useState("");
  const [selectedFriendIds, setSelectedFriendIds] = useState<string[]>([]);
  const [isCreatingGroup, setIsCreatingGroup] = useState(false);
  const [groupError, setGroupError] = useState("");

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

  const handleAddSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!addEmail.trim()) return;
    
    setIsAdding(true);
    setAddError("");
    
    try {
      await onAddFriend(addEmail);
      setAddEmail("");
      setIsAddModalOpen(false); // Close on success
    } catch (err) {
      setAddError("Failed to send request. Check the email.");
    } finally {
      setIsAdding(false);
    }
  };

  const handleGroupSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!groupName.trim() || selectedFriendIds.length === 0) {
      setGroupError("Please enter a name and select at least one friend.");
      return;
    }

    setIsCreatingGroup(true);
    setGroupError("");

    try {
      await onCreateGroup(groupName, selectedFriendIds);
      // Reset State
      setGroupName("");
      setSelectedFriendIds([]);
      setIsGroupModalOpen(false);
    } catch (err) {
      setGroupError("Failed to create group.");
    } finally {
      setIsCreatingGroup(false);
    }
  };

  const toggleFriendSelection = (uid: string) => {
    setSelectedFriendIds(prev => 
      prev.includes(uid) 
        ? prev.filter(id => id !== uid) 
        : [...prev, uid]
    );
  };

  return (
    <div className="h-full w-full bg-slate-900 text-slate-100 flex flex-col">
      {/* --- Header --- */}
      <div className="p-6 border-b border-slate-800 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        
        {/* Title (Icon Changed to generic Users) */}
        <h1 className="text-2xl font-bold flex items-center gap-3">
          <Users className="w-6 h-6 text-indigo-400" /> 
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
              className="w-full rounded-xl bg-slate-800/60 pl-9 pr-4 py-2 text-sm placeholder:text-slate-400 outline-none focus:ring-2 focus:ring-indigo-400 transition"
            />
          </div>

          {/* Add Friend Button */}
          <button 
            onClick={() => setIsAddModalOpen(true)}
            className="flex items-center gap-2 px-4 py-2 bg-indigo-600 hover:bg-indigo-500 text-white rounded-xl font-medium transition text-sm shadow-lg shadow-indigo-500/20"
          >
            <UserPlus size={18} />
          </button>

          {/* Create Group Button */}
          <button 
            onClick={() => setIsGroupModalOpen(true)}
            className="flex items-center gap-2 px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-xl font-medium transition text-sm shadow-lg"
          >
            <MessageSquarePlus size={18} />
            <span className="hidden sm:inline">New Group</span>
          </button>
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

      {isAddModalOpen && (
        <div className="absolute inset-0 z-50 flex items-center justify-center bg-slate-900/80 backdrop-blur-sm p-4">
          <div className="bg-slate-800 border border-slate-700 p-6 rounded-2xl w-full max-w-md shadow-2xl animate-in zoom-in-95 duration-200">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-xl font-semibold text-white">Send Friend Request</h3>
              <button onClick={() => setIsAddModalOpen(false)} className="text-slate-400 hover:text-white">
                <X size={20} />
              </button>
            </div>
            
            <form onSubmit={(e) => handleAddSubmit(e)}>
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
                  className="w-full bg-slate-900 border border-slate-700 rounded-xl px-4 py-3 text-slate-100 placeholder:text-slate-500 focus:ring-2 focus:ring-indigo-500 outline-none"
                />
                {addError && <p className="text-red-400 text-sm mt-2">{addError}</p>}
              </div>

              <div className="flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setIsAddModalOpen(false)}
                  className="px-4 py-2 text-slate-300 hover:text-white hover:bg-slate-700 rounded-lg transition"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isAdding}
                  className="px-4 py-2 bg-indigo-600 hover:bg-indigo-500 text-white rounded-lg font-medium transition flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {isAdding ? <Loader2 className="animate-spin" size={18} /> : <UserPlus size={18} />}
                  Send Request
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {isGroupModalOpen && (
        <div className="absolute inset-0 z-50 flex items-center justify-center bg-slate-900/80 backdrop-blur-sm p-4">
          <div className="bg-slate-800 border border-slate-700 p-6 rounded-2xl w-full max-w-lg shadow-2xl animate-in zoom-in-95 duration-200 flex flex-col max-h-[90vh]">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-xl font-semibold text-white">Create Group</h3>
              <button onClick={() => setIsGroupModalOpen(false)} className="text-slate-400 hover:text-white">
                <X size={20} />
              </button>
            </div>
            
            <form onSubmit={handleGroupSubmit} className="flex flex-col flex-1 overflow-hidden">
              <div className="mb-4">
                <label className="block text-sm font-medium text-slate-400 mb-1">
                  Group Name
                </label>
                <input
                  type="text"
                  autoFocus
                  required
                  value={groupName}
                  onChange={(e) => setGroupName(e.target.value)}
                  placeholder="e.g. Project Team"
                  className="w-full bg-slate-900 border border-slate-700 rounded-xl px-4 py-3 text-slate-100 placeholder:text-slate-500 focus:ring-2 focus:ring-indigo-500 outline-none"
                />
              </div>

              <div className="flex-1 overflow-hidden flex flex-col">
                <label className="block text-sm font-medium text-slate-400 mb-2">
                  Select Members ({selectedFriendIds.length})
                </label>
                
                <div className="flex-1 overflow-y-auto border border-slate-700 rounded-xl bg-slate-900/50 p-2 space-y-1">
                  {friends.length === 0 ? (
                    <p className="text-slate-500 text-sm text-center py-4">You have no friends to add yet.</p>
                  ) : (
                    friends.map(friend => {
                      const isSelected = selectedFriendIds.includes(friend.uid);
                      return (
                        <div 
                          key={friend.uid}
                          onClick={() => toggleFriendSelection(friend.uid)}
                          className={`flex items-center gap-3 p-2 rounded-lg cursor-pointer transition ${
                            isSelected ? 'bg-indigo-600/20 border border-indigo-500/50' : 'hover:bg-slate-800 border border-transparent'
                          }`}
                        >
                          {/* Custom Checkbox Appearance */}
                          <div className={`w-5 h-5 rounded-md border flex items-center justify-center transition ${
                            isSelected ? 'bg-indigo-500 border-indigo-500' : 'border-slate-600'
                          }`}>
                            {isSelected && <Check size={14} className="text-white" />}
                          </div>
                          
                          {renderAvatar(friend.name)}
                          
                          <div className="min-w-0">
                            <p className={`text-sm font-medium ${isSelected ? 'text-indigo-200' : 'text-slate-200'}`}>
                              {friend.name}
                            </p>
                          </div>
                        </div>
                      )
                    })
                  )}
                </div>
              </div>

              {groupError && <p className="text-red-400 text-sm mt-4">{groupError}</p>}

              <div className="flex justify-end gap-3 mt-6 pt-4 border-t border-slate-700">
                <button
                  type="button"
                  onClick={() => setIsGroupModalOpen(false)}
                  className="px-4 py-2 text-slate-300 hover:text-white hover:bg-slate-700 rounded-lg transition"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isCreatingGroup || selectedFriendIds.length === 0}
                  className="px-4 py-2 bg-indigo-600 hover:bg-indigo-500 text-white rounded-lg font-medium transition flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {isCreatingGroup ? <Loader2 className="animate-spin" size={18} /> : <MessageSquarePlus size={18} />}
                  Create Group
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

    </div>
  );
}