import { useMemo, useRef, useState, useLayoutEffect, useEffect } from "react";
import { useChat, type MsgSlice } from "../context/ChatContext";
import { Composer } from "./Composer";
import { useAuth } from "../context/AuthContext";
import type { Message } from "../wire/message";
import type { Conversation } from "../wire/conversation";
import { MoreVertical, Ban, LogOut, Users, UserPlus, Clock, AlertCircle, MessageSquare } from "lucide-react";
import { ParticipantsPanel } from "./ParticipantsPanel";
import { AddMembersModal } from "./AddMembersModal";
import { avatarColor, nameInitials } from "../utils/avatar";
import { api } from "../api/api";

type MessagePaneProps = {
  conv?: Conversation;
  msg: MsgSlice;
  activeConvId: string;
  onLoadMore: () => Promise<void>;
  onBlock?: () => void;
  onLeaveGroup?: () => void;
};



function Avatar({ name }: { name: string }) {
  return (
    <div title={name} className={`w-8 h-8 rounded-full flex items-center justify-center text-white text-xs font-semibold shrink-0 ${avatarColor(name)}`}>
      {nameInitials(name)}
    </div>
  );
}

function getDateLabel(ts: number): string {
  const now = new Date();
  const msgDate = new Date(ts);
  const todayStart = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime();
  const msgStart = new Date(msgDate.getFullYear(), msgDate.getMonth(), msgDate.getDate()).getTime();
  const diffDays = Math.round((todayStart - msgStart) / 86400000);

  if (diffDays === 0) return "Today";
  if (diffDays === 1) return "Yesterday";
  if (diffDays < 7) return msgDate.toLocaleDateString("en-US", { weekday: "long" });
  if (msgDate.getFullYear() === now.getFullYear()) {
    return msgDate.toLocaleDateString("en-US", { month: "short", day: "numeric" });
  }
  return msgDate.toLocaleDateString("en-US", { month: "short", day: "numeric", year: "numeric" });
}

// ---------- Messages pane ----------
export default function MessagesPane({ conv, msg, activeConvId, onLoadMore, onBlock, onLeaveGroup }: MessagePaneProps) {
  const { user } = useAuth();
  const { convs, markAsRead, friendships, loadConvs } = useChat();

  const participantMap = useMemo(() => {
    const map = new Map<string, string>();
    conv?.participants.forEach((p) => map.set(p.uid, p.name));
    return map;
  }, [conv?.participants]);

  const eligibleFriends = useMemo(() => {
    if (!conv || conv.type !== "group") return [];
    const participantUids = new Set(conv.participants.map((p) => p.uid));
    return friendships.filter((f) => f.status === "accepted" && !participantUids.has(f.uid));
  }, [friendships, conv]);

  const msgList = useMemo(() => {
    const confirmed = msg.order.map((seq) => msg.entities[seq]).filter(Boolean);
    const pending = (msg.pendingOrder || [])
      .map((id) => msg.pendingEntities[id])
      .filter(Boolean);
    return [...confirmed, ...pending];
  }, [msg]);

  type DateSeparator = { type: "date"; label: string; key: string };
  type MsgItem = { type: "msg"; msg: Message; first: boolean };

  const renderedItems = useMemo(() => {
    const items: Array<DateSeparator | MsgItem> = [];
    let lastLabel = "";
    let prevFromUid: string | null = null;
    for (const m of msgList) {
      if (m.ts) {
        const label = getDateLabel(m.ts);
        if (label !== lastLabel) {
          items.push({ type: "date", label, key: `sep-${label}` });
          lastLabel = label;
          prevFromUid = null;
        }
      }
      items.push({ type: "msg", msg: m, first: m.fromUid !== prevFromUid });
      prevFromUid = m.fromUid;
    }
    return items;
  }, [msgList]);

  const boxRef = useRef<HTMLDivElement | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const prevScrollHeightRef = useRef<number | null>(null);

  const hasScrolledToUnreadRef = useRef(false);
  const isAtBottomRef = useRef(true);
  
  useEffect(() => {
    markAsRead(activeConvId);
  }, [activeConvId, markAsRead, msg]);

  // Uses useLayoutEffect to adjust scroll position after DOM updates
  useLayoutEffect(() => {
    const box = boxRef.current;
    if (!box) return;

    // restore scroll position
    if (prevScrollHeightRef.current !== null) {
      const newHeight = box.scrollHeight;
      const diff = newHeight - prevScrollHeightRef.current;
      box.scrollTop = diff;
      prevScrollHeightRef.current = null;
    }

    // only the first time would trigger before hasScrolledToUnreadRef is set the true
    if (!hasScrolledToUnreadRef.current && msgList.length > 0) {
      const conv = convs.entities[activeConvId];

      if (conv && conv.lastReadSeq) {
        // Locate the element by ID
        const element = document.getElementById(`msg-${conv.lastReadSeq}`);

        if (element) {
          element.scrollIntoView({ block: "start" });
          hasScrolledToUnreadRef.current = true;
          return; // Stop here
        }
      }
      // If no message found, scroll to bottom
      hasScrolledToUnreadRef.current = true;
    }

    // Check if the user is at the bottom for not flashing to bottom if scrolled up
    if (isAtBottomRef.current) {
      box.scrollTop = box.scrollHeight;
    }
  }, [msgList, activeConvId]);

  const handleScroll = (e: React.UIEvent<HTMLDivElement>) => {
    const target = e.currentTarget;
    
    // scrollHeight: total height of the content
    // scrollTop: how much hidden above the visible area
    // clientHeight: height of the visible area
    const isAtBottom = target.scrollHeight - target.scrollTop - target.clientHeight < 20;
    isAtBottomRef.current = isAtBottom;
    // Trigger when hitting top
    if (
      target.scrollTop === 0 &&
      !isLoading &&
      msgList.length > 0 &&
      target.scrollHeight > target.clientHeight
    ) {
      setIsLoading(true);
      prevScrollHeightRef.current = target.scrollHeight; // Snapshot height
      onLoadMore().finally(() => setIsLoading(false));
      console.log("Loading more messages...");
    }
  };

  const isSelf = (msg: Message) => {
    return user ? msg.fromUid === user.uid : false;
  };

  const [menuOpen, setMenuOpen] = useState(false);
  const [panelOpen, setPanelOpen] = useState(false);
  const [addUserOpen, setAddUserOpen] = useState(false);

  const handleAddParticipants = async (selected: { uid: string; name: string }[]) => {
    if (!conv || selected.length === 0) return;
    await Promise.all(
      selected.map((s) =>
        api("/conversation/participant", {
          method: "POST",
          body: JSON.stringify({ participant_uid: s.uid, conversation_id: conv.id, role: "member" }),
        })
      )
    );
    const newParticipants = selected.map((s) => ({ uid: s.uid, name: s.name, role: "member" as const }));
    loadConvs([{ ...conv, participants: [...conv.participants, ...newParticipants] }]);
  };

  return (
    <div className={`relative h-full grid ${conv ? "grid-rows-[auto_1fr_auto]" : "grid-rows-[1fr_auto]"}`}>
      {/* header */}
      {conv && (
        <div
          className="relative flex items-center justify-center px-4 py-3 border-b"
          style={{ background: "var(--panel)", borderColor: "var(--panel-border)" }}
        >
          <span className="font-semibold" style={{ color: "var(--text)" }}>
            {conv.title}{conv.type === "group" && <span className="font-normal" style={{ color: "var(--text-soft)" }}> ({conv.participants.length})</span>}
          </span>
          <div className="absolute right-4 flex items-center gap-1">
            {/* Participants panel button — group only */}
            {conv.type === "group" && (
              <button
                onClick={() => setPanelOpen((o) => !o)}
                className="p-2 rounded-lg transition"
                style={
                  panelOpen
                    ? { background: "color-mix(in srgb, var(--text) 10%, transparent)", color: "var(--text)" }
                    : { color: "var(--text-soft)" }
                }
                title="Members"
              >
                <Users size={18} />
              </button>
            )}
            {/* Add user button — group only */}
            {conv.type === "group" && (
              <button
                onClick={() => setAddUserOpen(true)}
                className="p-2 rounded-lg transition"
                style={
                  addUserOpen
                    ? { background: "color-mix(in srgb, var(--text) 10%, transparent)", color: "var(--text)" }
                    : { color: "var(--text-soft)" }
                }
                title="Add member"
              >
                <UserPlus size={18} />
              </button>
            )}
            {/* More options menu */}
            {(conv.type === "direct" && onBlock) || (conv.type === "group" && onLeaveGroup) ? (
              <div className="relative">
                <button
                  onClick={() => setMenuOpen((o) => !o)}
                  className="p-2 rounded-lg transition"
                  style={{ color: "var(--text-soft)" }}
                >
                  <MoreVertical size={18} />
                </button>
                {menuOpen && (
                  <>
                    <div className="fixed inset-0 z-10" onClick={() => setMenuOpen(false)} />
                    <div
                      className="absolute right-0 mt-1 w-44 border rounded-xl shadow-xl z-20 py-1"
                      style={{ background: "var(--sidebar-item)", borderColor: "var(--panel-border)" }}
                    >
                      {conv.type === "direct" && onBlock && (
                        <button
                          onClick={() => { setMenuOpen(false); onBlock(); }}
                          className="w-full text-left px-4 py-3 text-sm flex items-center gap-3 transition"
                          style={{ color: "var(--warning-text)" }}
                          onMouseEnter={e => (e.currentTarget.style.background = "var(--warning-bg)")}
                          onMouseLeave={e => (e.currentTarget.style.background = "")}
                        >
                          <Ban size={16} />
                          Block User
                        </button>
                      )}
                      {conv.type === "group" && onLeaveGroup && (
                        <button
                          onClick={() => { setMenuOpen(false); onLeaveGroup(); }}
                          className="w-full text-left px-4 py-3 text-sm flex items-center gap-3 transition"
                          style={{ color: "var(--danger-text)" }}
                          onMouseEnter={e => (e.currentTarget.style.background = "var(--danger-bg)")}
                          onMouseLeave={e => (e.currentTarget.style.background = "")}
                        >
                          <LogOut size={16} />
                          Leave Group
                        </button>
                      )}
                    </div>
                  </>
                )}
              </div>
            ) : null}
          </div>
        </div>
      )}

      {/* Participants panel */}
      {conv?.type === "group" && panelOpen && (
        <ParticipantsPanel
          participants={conv.participants}
          onClose={() => setPanelOpen(false)}
        />
      )}

      {/* Add members modal */}
      {conv?.type === "group" && addUserOpen && (
        <AddMembersModal
          friends={eligibleFriends}
          onClose={() => setAddUserOpen(false)}
          onAdd={handleAddParticipants}
        />
      )}

      {/* messages list */}
      <div
        ref={boxRef}
        onScroll={handleScroll}
        className="overflow-y-auto px-4 py-2"
      >
        {renderedItems.map((item) =>
          item.type === "date" ? (
            <div key={item.key} className="flex items-center gap-3 my-3">
              <div className="flex-1 h-px" style={{ background: "var(--panel-border)" }} />
              <span className="text-xs font-medium" style={{ color: "var(--text-faint)" }}>{item.label}</span>
              <div className="flex-1 h-px" style={{ background: "var(--panel-border)" }} />
            </div>
          ) : (
            <div
              id={item.msg.seq ? `msg-${item.msg.seq}` : undefined}
              key={item.msg.id || item.msg.clientMsgId}
              className={[
                "flex items-end gap-2",
                isSelf(item.msg) ? "justify-end" : "justify-start",
                item.first ? "mt-3" : "mt-0.5",
              ].join(" ")}
            >
              {/* Avatar for other users — only on the first message of a run */}
              {!isSelf(item.msg) && (
                item.first ? (
                  <Avatar name={participantMap.get(item.msg.fromUid) ?? item.msg.fromName ?? "?"} />
                ) : (
                  <div className="w-8 shrink-0" />
                )
              )}

              <div className="flex flex-col max-w-[70%]">
                {/* Sender name for group chats — only on the first message of a run */}
                {!isSelf(item.msg) && conv?.type === "group" && item.first && (
                  <span className="text-xs mb-1 ml-1" style={{ color: "var(--text-faint)" }}>
                    {participantMap.get(item.msg.fromUid) ?? item.msg.fromName ?? "Unknown"}
                  </span>
                )}
                <div
                  className={[
                    "rounded-2xl px-3 py-2 text-[15px]",
                    item.msg.status === "failed"
                      ? "border"
                      : isSelf(item.msg)
                      ? "text-white rounded-br-sm"
                      : "rounded-bl-sm",
                    item.msg.status === "sending" ? "opacity-70" : "",
                  ].join(" ")}
                  style={
                    item.msg.status === "failed"
                      ? { background: "var(--danger-bg)", color: "var(--danger-text)", borderColor: "var(--danger-bg-hover)" }
                      : { background: isSelf(item.msg) ? "var(--primary)" : "var(--bubble-other)", color: isSelf(item.msg) ? "white" : "var(--text)" }
                  }
                >
                  <div className="flex flex-wrap items-end justify-end gap-x-2">
                    <p className="whitespace-pre-wrap break-words text-left min-w-0">
                      {item.msg.text}
                    </p>
                    <span className="flex items-center gap-1 ml-auto pt-0.5 pb-0.5">
                      {isSelf(item.msg) && item.msg.status === "sending" && (
                        <Clock size={10} className="opacity-70" />
                      )}
                      {isSelf(item.msg) && item.msg.status === "failed" && (
                        <AlertCircle size={11} />
                      )}
                      <span className="text-[10px] opacity-60 leading-none">
                        {new Date(item.msg.ts ?? 0).toLocaleTimeString([], {
                          hour: "2-digit",
                          minute: "2-digit",
                        })}
                      </span>
                    </span>
                  </div>
                </div>
              </div>
            </div>
          )
        )}
        {msgList.length === 0 && (
          <div className="h-full flex flex-col items-center justify-center gap-2" style={{ color: "var(--text-faint)" }}>
            <MessageSquare size={28} className="opacity-50" />
            <p className="text-sm">No messages yet. Say hi!</p>
          </div>
        )}
      </div>

      {/* composer */}
      <Composer convId={activeConvId} />
    </div>
  );
}
