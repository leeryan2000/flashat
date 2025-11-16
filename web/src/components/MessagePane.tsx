import { useMemo, useRef, useEffect} from "react";
import type { MsgSlice} from "../context/ChatContext";
import { Composer } from "./Composer";
import { useAuth } from "../context/AuthContext";
import type { Message } from "../context/ChatContext";


// ***** messages in the conversation sometimes would display your message as others message
// ---------- Messages pane ----------
export default function MessagesPane({ msg, activeConvId }:{msg: MsgSlice, activeConvId: string}) {
  const { user, isAuthenticated } =  useAuth();
  const list = useMemo(() => msg.order.map(seq => msg.entities[seq]).filter(Boolean) ?? [], [msg]);
  const boxRef = useRef<HTMLDivElement|null>(null);

  useEffect(()=>{
    // auto scroll to bottom on convo change or new message
    boxRef.current?.scrollTo({ top: boxRef.current.scrollHeight })
  }, [list.length]);

  if (!isAuthenticated) 
    return (
      <div className="h-full grid grid-rows-[1fr_auto]">
        <div className="overflow-y-auto pr-2 grid place-items-center text-slate-400">
          Loading…
        </div>
        <Composer convId={activeConvId} />
      </div>
    );

  const isSelf = (msg: Message) => {
    console.log("Checking isSelf for msg:", user ? msg.fromUid === user.uid : false);
    return user ? msg.fromUid === user.uid : false;
  }
  
  return (
    <div className="h-full grid grid-rows-[1fr_auto]">
      {/* messages list */}
      <div ref={boxRef} className="overflow-y-auto space-y-2 pr-2">
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

