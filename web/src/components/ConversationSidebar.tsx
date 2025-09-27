import { useMemo, useState } from "react";
import type { Conversation, ConvState } from "../context/ChatContext";


function timeAgo(ts: number) {
  const mins = Math.max(1, Math.round((Date.now() - ts) / 60000));
  if (mins < 60) return `${mins}m`;
  const hrs = Math.round(mins / 60);
  if (hrs < 24) return `${hrs}h`;
  return `${Math.round(hrs / 24)}d`;
}

export default function ConversationsSidebar({ convs, activeConvId, onOpen }:{ convs: ConvState; activeConvId?: string; onOpen:(id:string)=>void }) {
  // Query for the title search
  const [q, setQ] = useState("");
  // filters the conversations by title
  const filtered = useMemo(() => {
    const needle = q.trim().toLowerCase();
    return convs.order
    .map(id => convs.entities[id])
    .filter((c): c is Conversation => !!c) 
    .filter(conv => conv.title?.toLowerCase().includes(needle));

  }, [convs.order, convs.entities, q]);

  return (
    <div className="h-full bg-slate-900 text-slate-100 w-80 flex-shrink-0 grid grid-rows-[auto_1fr]">
      {/* search */}
      <div className="p-3 border-b border-slate-800">
        <input
          value={q}
          onChange={e=>setQ(e.target.value)}
          placeholder="Search conversations"
          className="w-full rounded-xl bg-slate-800/60 px-3 py-2 text-sm placeholder:text-slate-400 outline-none focus:ring-2 focus:ring-indigo-400"
        />
      </div>

      {/* list */}
      <div className="overflow-y-auto">
        {filtered.map(conv => (
          <button
            key={conv.id}
            onClick={()=>onOpen(conv.id)}
            className={[
              "w-full text-left px-3 py-3 gap-3 flex items-start",
              activeConvId===conv.id ? "bg-slate-800" : "hover:bg-slate-800/60"
            ].join(" ")}
          >
            <div className="h-10 w-10 rounded-full bg-slate-700 flex items-center justify-center text-xs">{conv.title.slice(0,2)}</div>
            <div className="min-w-0 flex-1">
              <div className="flex items-center gap-2">
                <p className="truncate font-medium">{conv.title}</p>
                <span className="text-xs text-slate-400 ml-auto">
                  {typeof conv.lastMsgTs === "number" ? timeAgo(conv.lastMsgTs) : "--"}
                </span>
              </div>
              <p className="truncate text-sm text-slate-300">{conv.lastMsgText}</p>
            </div>
            {!!conv.unreadCount && (
              <span className="ml-2 inline-flex items-center justify-center rounded-full bg-emerald-500 text-white text-xs px-2 py-0.5">
                {conv.unreadCount}
              </span>
            )}
          </button>
        ))}
        {filtered.length===0 && (
          <div className="p-6 text-sm text-slate-400">No conversations</div>
        )}
      </div>
    </div>
  );
}