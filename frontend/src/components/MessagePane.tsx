import { useMemo, useRef, useState, useLayoutEffect, useEffect } from "react";
import { useChat, type MsgSlice } from "../context/ChatContext";
import { Composer } from "./Composer";
import { useAuth } from "../context/AuthContext";
import type { Message } from "../wire/message";
import type { Conversation } from "../wire/conversation";
import { MoreVertical, Ban } from "lucide-react";

type MessagePaneProps = {
  conv?: Conversation;
  msg: MsgSlice;
  activeConvId: string;
  onLoadMore: () => Promise<void>;
  onBlock?: () => void;
};

// ---------- Messages pane ----------
export default function MessagesPane({ conv, msg, activeConvId, onLoadMore, onBlock }: MessagePaneProps) {
  const { user } = useAuth();
  const { convs, markAsRead } = useChat();

  const msgList = useMemo(() => {
    const confirmed = msg.order.map((seq) => msg.entities[seq]).filter(Boolean);
    // Map pending IDs to message objects
    const pending = (msg.pendingOrder || [])
      .map((id) => msg.pendingEntities[id])
      .filter(Boolean);

    return [...confirmed, ...pending];
  }, [msg]);

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

  return (
    <div className={`h-full grid ${conv ? "grid-rows-[auto_1fr_auto]" : "grid-rows-[1fr_auto]"}`}>
      {/* header */}
      {conv && (
        <div className="relative flex items-center justify-center px-4 py-3 border-b border-slate-200 bg-white">
          <span className="font-semibold text-slate-800">{conv.title}</span>
          {conv.type === "direct" && onBlock && (
            <div className="absolute right-4">
              <button
                onClick={() => setMenuOpen((o) => !o)}
                className="p-2 rounded-lg text-slate-400 hover:bg-slate-100 hover:text-slate-700 transition"
              >
                <MoreVertical size={18} />
              </button>
              {menuOpen && (
                <>
                  <div className="fixed inset-0 z-10" onClick={() => setMenuOpen(false)} />
                  <div className="absolute right-0 mt-1 w-44 bg-white border border-slate-200 rounded-xl shadow-lg z-20 py-1">
                    <button
                      onClick={() => { setMenuOpen(false); onBlock(); }}
                      className="w-full text-left px-4 py-3 text-sm text-orange-500 hover:bg-orange-50 flex items-center gap-3 transition"
                    >
                      <Ban size={16} />
                      Block User
                    </button>
                  </div>
                </>
              )}
            </div>
          )}
        </div>
      )}

      {/* messages list */}
      <div
        ref={boxRef}
        onScroll={handleScroll}
        className="overflow-y-auto space-y-2 px-4 py-2"
      >
        {msgList.map((msg) => (
          <div
            id={msg.seq ? `msg-${msg.seq}` : undefined}
            key={msg.id || msg.clientMsgId}
            className={[
              "flex",
              isSelf(msg) ? "justify-end" : "justify-start",
            ].join(" ")}
          >
            <div
              className={[
                "max-w-[75%] rounded-2xl px-3 py-2 text-sm shadow",
                msg.status === "failed"
                  ? "bg-red-100 text-red-800 border border-red-300" // Failed Style
                  : isSelf(msg)
                  ? "bg-indigo-600 text-white rounded-br-sm" // Self Sent Style
                  : "bg-slate-100 text-slate-900 rounded-bl-sm", // Other Style
                msg.status === "sending" ? "opacity-70" : "",
              ].join(" ")}
            >
              <p
                className={`whitespace-pre-wrap break-words ${
                  isSelf(msg) ? "text-right" : "text-left"
                }`}
              >
                {msg.text}
              </p>

              <div className="flex items-center justify-end gap-1 mt-1">
                {isSelf(msg) && msg.status === "sending" && (
                  <span className="text-[10px] italic">...</span>
                )}

                {isSelf(msg) && msg.status === "failed" && (
                  <span className="text-[10px] font-bold">X</span>
                )}

                {/* Timestamp is now last (rightmost) */}
                <p className="text-[11px] opacity-70">
                  {new Date(msg.ts ?? 0).toLocaleTimeString([], {
                    hour: "2-digit",
                    minute: "2-digit",
                  })}
                </p>
              </div>
            </div>
          </div>
        ))}
        {msgList.length === 0 && (
          <div className="h-full grid place-items-center text-slate-400">
            No messages yet.
          </div>
        )}
      </div>

      {/* composer */}
      <Composer convId={activeConvId} />
    </div>
  );
}
