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
  const list = useMemo(() => msg.order.map(seq => msg.entities[seq]).filter(Boolean) ?? [], [msg]);
  const boxRef = useRef<HTMLDivElement|null>(null);
  const [isLoading, setIsLoading] =  useState(false);
  const prevScrollHeightRef = useRef<number|null>(null);

  useLayoutEffect(() => {
    const box = boxRef.current;
    if (!box) return;

    // 1. Restore scroll position if we just loaded history
    if (prevScrollHeightRef.current !== null) {
      const newHeight = box.scrollHeight;
      const diff = newHeight - prevScrollHeightRef.current;
      box.scrollTop = diff;
      prevScrollHeightRef.current = null;
    } 
    // 2. Otherwise, scroll to bottom (new message sent/received)
    else {
      box.scrollTop = box.scrollHeight;
    }
  }, [list]);
  
  const handleScroll = (e: React.UIEvent<HTMLDivElement>) => {
    const target = e.currentTarget;
    // Trigger when hitting top
    if (target.scrollTop === 0 && !isLoading && list.length > 0) {
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
        {list.map(msg => (
          <div key={msg.id} className={["flex", isSelf(msg) ?"justify-end":"justify-start"].join(" ") }>
            <div className={[
              "max-w-[75%] rounded-2xl px-3 py-2 text-sm shadow",
              isSelf(msg) ? "bg-indigo-600 text-white rounded-br-sm" : "bg-slate-100 text-slate-900 rounded-bl-sm"
            ].join(" ")}>
              <p className="whitespace-pre-wrap break-words">{msg.text}</p>
              <p className="mt-1 text-[11px] opacity-70">{new Date(msg.ts ?? 0).toLocaleTimeString()}</p>
            </div>
          </div>
        ))}
        {list.length===0 && (
          <div className="h-full grid place-items-center text-slate-400">No messages yet.</div>
        )}
      </div>

      {/* composer */}
      <Composer convId={activeConvId} />
    </div>
  );
}

