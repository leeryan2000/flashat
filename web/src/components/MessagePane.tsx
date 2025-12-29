import { useMemo, useRef, useState, useLayoutEffect} from "react";
import type { MsgSlice} from "../context/ChatContext";
import { Composer } from "./Composer";
import { useAuth } from "../context/AuthContext";
import type { Message } from "../context/ChatContext";

type MessagePaneProps = {
  msg: MsgSlice;
  activeConvId: string;
  onLoadMore: () => Promise<void>; // Function to fetch older messages
}

// ---------- Messages pane ----------
export default function MessagesPane({ msg, activeConvId, onLoadMore}: MessagePaneProps) {
  const { user } =  useAuth();

  const msgList = useMemo(() => {
    const confirmed = msg.order.map(seq => msg.entities[seq]).filter(Boolean);
    // Map pending IDs to message objects
    const pending = (msg.pendingOrder || []).map(id => msg.pendingEntities[id]).filter(Boolean);
    
    return [...confirmed, ...pending];
  }, [msg]);

  const boxRef = useRef<HTMLDivElement|null>(null);
  const [isLoading, setIsLoading] =  useState(false);
  const prevScrollHeightRef = useRef<number|null>(null);

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
    // Otherwise, scroll to bottom (new message sent/received)
    else {
      box.scrollTop = box.scrollHeight;
    }
  }, [msgList]);
  
  const handleScroll = (e: React.UIEvent<HTMLDivElement>) => {
    const target = e.currentTarget;
    // Trigger when hitting top
    if (target.scrollTop === 0 && !isLoading && msgList.length > 0) {
      setIsLoading(true);
      prevScrollHeightRef.current = target.scrollHeight; // Snapshot height
      onLoadMore().finally(() => setIsLoading(false));
      console.log("Loading more messages...");
    }
  };

  const isSelf = (msg: Message) => {
    return user ? msg.fromUid === user.uid : false;
  }

  
  return (
    <div className="h-full grid grid-rows-[1fr_auto]">
      {/* messages list */}
      <div ref={boxRef} onScroll={handleScroll} className="overflow-y-auto space-y-2 pr-2">
        {msgList.map(msg => (
          <div key={msg.id || msg.clientMsgId || msg.ts} className={["flex", isSelf(msg) ?"justify-end":"justify-start"].join(" ") }>
            <div className={[
              "max-w-[75%] rounded-2xl px-3 py-2 text-sm shadow",
              msg.status === "failed" 
                ? "bg-red-100 text-red-800 border border-red-300" // Failed Style
                : isSelf(msg) 
                  ? "bg-indigo-600 text-white rounded-br-sm"     // Self Sent Style
                  : "bg-slate-100 text-slate-900 rounded-bl-sm", // Other Style
              msg.status === "sending" ? "opacity-70" : ""
           ].join(" ")}>
              <p className={`whitespace-pre-wrap break-words ${isSelf(msg) ? "text-right" : "text-left"}`}>
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
                  {new Date(msg.ts ?? 0).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})}
                </p>
                
              </div>
            </div>
          </div>
        ))}
        {msgList.length===0 && (
          <div className="h-full grid place-items-center text-slate-400">No messages yet.</div>
        )}
      </div>

      {/* composer */}
      <Composer convId={activeConvId} />
    </div>
  );
}

