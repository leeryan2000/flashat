import { useMemo, useRef, useEffect} from "react";
import type { MsgSlice} from "../context/ChatContext";
import { Composer } from "./Composer";

// const MOCK_MESSAGES: Message[] = [
//   { id: "m1", convId: "e745cf8a-4e24-4a3e-a85f-54464d2beadc", from: "Alice", text: "Hey!", ts: Date.now() - 1000 * 60 * 60, self: false },
//   { id: "m2", convId: "e745cf8a-4e24-4a3e-a85f-54464d2beadc", from: "Huai'en", text: "Yo!", ts: Date.now() - 1000 * 60 * 58, self: true },
//   { id: "m3", convId: "e745cf8a-4e24-4a3e-a85f-54464d2beadc", from: "Alice", text: "See you tomorrow!", ts: Date.now() - 1000 * 60 * 3, self: false },
//   { id: "m4", convId: "c2", from: "Bob", text: "Pushed a new branch.", ts: Date.now() - 1000 * 60 * 11, self: false },
//   { id: "m5", convId: "c3", from: "Jessie", text: "Prepare for trouble.", ts: Date.now() - 1000 * 60 * 60 * 5, self: false },
//   { id: "m6", convId: "c3", from: "James", text: "And make it double!", ts: Date.now() - 1000 * 60 * 60 * 5 + 10000, self: false },
//   { id: "m7", convId: "c4", from: "Huai'en", text: "Drafting the layout now", ts: Date.now() - 1000 * 60 * 60 * 10, self: true },
// ];



// ---------- Messages pane ----------
export default function MessagesPane({ msg, activeConvId }:{msg: MsgSlice, activeConvId: string}) {
  const list = useMemo(() => msg.order.map(seq => msg.entities[seq]).filter(Boolean) ?? [], [msg]);
  // const list = useMemo(() => MOCK_MESSAGES.filter(m => m.convId === activeConvId), [activeConvId]);
  const boxRef = useRef<HTMLDivElement|null>(null);

  useEffect(()=>{
    // auto scroll to bottom on convo change or new message
    boxRef.current?.scrollTo({ top: boxRef.current.scrollHeight })
  }, [list.length]);

  return (
    <div className="h-full grid grid-rows-[1fr_auto]">
      {/* messages list */}
      <div ref={boxRef} className="overflow-y-auto space-y-2 pr-2">
        {list.map(msg => (
          <div key={msg.id} className={["flex", msg.self ?"justify-end":"justify-start"].join(" ") }>
            <div className={[
              "max-w-[75%] rounded-2xl px-3 py-2 text-sm shadow",
              msg.self ? "bg-indigo-600 text-white rounded-br-sm" : "bg-slate-100 text-slate-900 rounded-bl-sm"
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

